package excludes

import (
	"AirExcludes/db"
	"AirExcludes/model/device"
	"AirExcludes/model/post"
	"database/sql"
	"errors"
	"github.com/tealeg/xlsx/v3"
	"log"
	"strconv"
	"strings"
	"time"
)

// Структура для получения данных с бд и обменом с фротом
type Exception struct {
	Id         int    `json:"id"`
	PostId     int    `json:"postId"`
	SensorType int    `json:"sensorType"`
	DateStart  string `json:"dateStart"`
	DateEnd    string `json:"dateEnd"`
	Comment    string `json:"comment"`
}

// Структура для длбавления в БД
type exceptionDb struct {
	id        int
	sensorId  int
	dateStart time.Time
	dateEnd   time.Time
	comment   string
}

// Набор изменений таблицы
type ModifiedExceptions struct {
	UpdatingRows []Exception `json:"updatingRows"`
	AddedRows    []Exception `json:"addedRows"`
	DeletedRows  []Exception `json:"deletedRows"`
}

// Преобразует в структуру для базы данных
func convertForDb(exceptions []Exception) (exceptionsDb []exceptionDb) {
	con := db.ConnectDB()
	defer con.Close()
	for i := 0; i < len(exceptions); i++ {
		var exceptionDb exceptionDb
		var query = `select top 1 id from krasecology.dbo.Svc_AirSensors where PostId=@p1 and SensorType=@p2`
		rows, err := con.Query(query, exceptions[i].PostId, exceptions[i].SensorType)
		if err != nil {
			rows.Close()
			log.Println(err)
			return
		}
		for rows.Next() {
			if err := rows.Scan(&exceptionDb.sensorId); err != nil {
				rows.Close()
				log.Println(err)
				return
			}
		}
		exceptionDb.id = exceptions[i].Id
		exceptionDb.dateStart, _ = time.Parse("2006-01-02T15:04", exceptions[i].DateStart)
		exceptionDb.dateEnd, _ = time.Parse("2006-01-02T15:04", exceptions[i].DateEnd)
		exceptionDb.comment = exceptions[i].Comment
		exceptionsDb = append(exceptionsDb, exceptionDb)
		err = rows.Close()
		if err != nil {
			return nil
		}
	}
	return exceptionsDb
}

// Собирает id из массива в строчку для удаления записей
func sliceInString(exceptions []exceptionDb) string {
	var str = ""
	for i := 0; i < len(exceptions); i++ {
		str += strconv.Itoa(exceptions[i].id)
		if i < len(exceptions)-1 {
			str = str + ","
		}
	}
	return str
}

// GetExceptions Получает "исключения" из БД
func GetExceptions() []Exception {
	var exceptions []Exception
	con := db.ConnectDB()
	defer con.Close()
	var comment sql.NullString
	var dateStart time.Time
	var dateEnd time.Time
	//Запрос на получение записей за последние 5 лет
	var query = `SELECT Id, PostId, SensorType, DateStart, DateEnd, Comment FROM krasecology.dbo.v_AIR_Excludes where DateEnd>='` + time.Now().AddDate(-10, 0, 0).Format("2006-02-01 15:04:05") + `' order by DateStart DESC`
	rows, err := con.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println("Не удалось выполнить запрос к БД: ", err)
		return exceptions
	}
	for rows.Next() {
		var exception Exception
		if err = rows.Scan(&exception.Id, &exception.PostId, &exception.SensorType, &dateStart, &dateEnd, &comment); err != nil {
			log.Println(err)
			return exceptions
		}
		exception.DateStart = dateStart.Format("2006-01-02T15:04")
		exception.DateEnd = dateEnd.Format("2006-01-02T15:04")
		if comment.Valid {
			exception.Comment = comment.String
		} else {
			exception.Comment = ""
		}

		exceptions = append(exceptions, exception)
	}
	return exceptions
}

// SaveExceptions Добавление изменений в БД
func SaveExceptions(updates ModifiedExceptions) error {
	//Замняем postId и sensorType  на sensorId
	//TODO перенести преобразование в кконечные функции
	var convertedDeletedRows = convertForDb(updates.DeletedRows)
	var convertedAddedRows = convertForDb(updates.AddedRows)
	var convertedUpdatingRows = convertForDb(updates.UpdatingRows)
	var allExceptions = GetExceptions()
	if len(convertedDeletedRows) > 0 {
		deleteRows(convertedDeletedRows)
		for _, row := range updates.DeletedRows {
			var erroneouslyReturnedSamples []Exception
			timeStart, _ := time.Parse("2006-01-02T15:04", row.DateStart)
			timeEnd, _ := time.Parse("2006-01-02T15:04", row.DateEnd)
			for _, exception := range allExceptions {
				//Пропускаем если это другой сенсор
				if exception.PostId != row.PostId || exception.SensorType != row.SensorType {
					continue
				}
				continueOuter := false
				for _, deletedRow := range updates.DeletedRows {
					if exception.Id == deletedRow.Id {
						continueOuter = true
						break
					}
				}
				//Пропускаем если exception есть в списке удаляемых
				if continueOuter {
					continue
				}
				time2Start, _ := time.Parse("2006-01-02T15:04", exception.DateStart)
				time2End, _ := time.Parse("2006-01-02T15:04", exception.DateEnd)
				if time2End.Before(timeStart) || time2End.Equal(timeStart) {
					continue
				}
				timeIntersection1 := maxDate(timeStart, time2Start)
				timeIntersection2 := minDate(timeEnd, time2End)
				if timeIntersection1.Before(timeIntersection2) || timeIntersection1.Equal(timeIntersection2) {
					erroneouslyReturnedSamples = append(erroneouslyReturnedSamples, Exception{
						SensorType: row.SensorType,
						PostId:     row.PostId,
						DateStart:  timeIntersection1.Format("2006-01-02T15:04"),
						DateEnd:    timeIntersection2.Format("2006-01-02T15:04"),
					})
				}
			}
			returnSamples(erroneouslyReturnedSamples)
		}
	}
	allExceptions = GetExceptions()
	if len(convertedAddedRows) > 0 {
		err := checkIntersection(updates.AddedRows, allExceptions)
		if err != nil {
			return err
		}
		addRows(convertedAddedRows)
	}
	if len(convertedUpdatingRows) > 0 {
		err := checkIntersection(updates.UpdatingRows, allExceptions)
		if err != nil {
			return err
		}
		updateRows(convertedUpdatingRows)
	}
	recalculationOfAverageDaily(convertedDeletedRows, convertedAddedRows, convertedUpdatingRows)
	return nil
}

// Возвращаем "исключения" после удаления одной из нескольких пересекающихся записей
func returnSamples(rows []Exception) {
	con := db.ConnectDB()
	defer con.Close()
	convertedRows := convertForDb(rows)
	query := `Update krasecology.dbo.Svc_AirData Set isConfirm = 0 WHERE [Type] = 1
		AND SensorId = @p1 and [Date]>dbo.ToUTCTime(@p2)
		AND [Date] <= dbo.ToUTCTime(@p3)`
	for _, row := range convertedRows {
		_, err := con.Exec(query, row.sensorId, row.dateStart, row.dateEnd)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func minDate(date1 time.Time, date2 time.Time) time.Time {
	if date1.Before(date2) {
		return date1
	}
	if date1.Equal(date2) {
		return date1
	}
	return date2
}
func maxDate(date1 time.Time, date2 time.Time) time.Time {
	if date1.Before(date2) {
		return date2
	}
	if date1.Equal(date2) {
		return date2
	}
	return date1
}

// Проверка пересечения периодов и дубликатов
func checkIntersection(rows []Exception, allExceptions []Exception) error {
	for _, row := range rows {
		var timeStart, _ = time.Parse("2006-01-02T15:04", row.DateStart)
		var timeEnd, _ = time.Parse("2006-01-02T15:04", row.DateEnd)
		for _, exception := range allExceptions {
			if row.Id == exception.Id {
				continue
			}
			var time2End, _ = time.Parse("2006-01-02T15:04", exception.DateEnd)
			var time2Start, _ = time.Parse("2006-01-02T15:04", exception.DateStart)
			if row.SensorType == exception.SensorType && row.PostId == exception.PostId {
				if timeStart.Before(time2End) && time2Start.Before(timeEnd) || timeStart.Equal(time2Start) || timeEnd.Equal(time2End) {
					return errors.New("периоды дат пересекаются")
				}
			}
		}
	}
	return nil
}

// Пересчитываем среднесуточные после изменений
func recalculationOfAverageDaily(convertedDeletedRows []exceptionDb, convertedAddedRows []exceptionDb, convertedUpdatingRows []exceptionDb) {
	var updatedSensors = make(map[int][]time.Time)
	var allUpdates = convertedDeletedRows
	allUpdates = append(allUpdates, convertedAddedRows...)
	allUpdates = append(allUpdates, convertedUpdatingRows...)

	for i := 0; i < len(allUpdates); i++ {
		dateStart := time.Date(allUpdates[i].dateStart.Year(), allUpdates[i].dateStart.Month(), allUpdates[i].dateStart.Day(), 0, 0, 0, 0, time.UTC)
		dateEnd := time.Date(allUpdates[i].dateEnd.Year(), allUpdates[i].dateEnd.Month(), allUpdates[i].dateEnd.Day(), 0, 0, 0, 0, time.UTC)
		today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

		for j := dateStart; j.Before(dateEnd.Add(24*time.Hour)) && j.Before(today); j = j.Add(24 * time.Hour) {
			if _, ok := updatedSensors[allUpdates[i].sensorId]; !ok {
				updatedSensors[allUpdates[i].sensorId] = make([]time.Time, 0)
			}
			if !exist(j, updatedSensors[allUpdates[i].sensorId]) {
				updatedSensors[allUpdates[i].sensorId] = append(updatedSensors[allUpdates[i].sensorId], j.Add(-time.Hour*7))
			}
		}
	}
	con := db.ConnectDB()
	defer con.Close()
	for sensorId, times := range updatedSensors {
		for _, timeAvg := range times {
			_, err := con.Exec("dbo.sp_air_calcAvg_day_withSensorAndDate", sensorId, timeAvg)
			if err != nil {
				log.Println("Ошибка пересчета среднесуточных", err)
			}
		}
	}
}

func exist(date time.Time, DateList []time.Time) bool {
	for i := 0; i < len(DateList); i++ {
		if date.Equal(DateList[i]) {
			return true
		}
	}
	return false
}

func deleteRows(exceptions []exceptionDb) {
	con := db.ConnectDB()
	defer con.Close()
	query := `DELETE FROM krasecology.dbo.Rpt_SKAT_Excludes WHERE Id in (` + sliceInString(exceptions) + `)`
	_, err := con.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
}

func updateRows(rows []exceptionDb) {
	con := db.ConnectDB()
	defer con.Close()
	query := `Update krasecology.dbo.Rpt_SKAT_Excludes Set SensorId = @p1, DateStart = @p2, DateEnd = @p3, Comment = @p4 WHERE id = @p5`
	for _, row := range rows {
		_, err := con.Exec(query, row.sensorId, row.dateStart, row.dateEnd, row.comment, row.id)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func addRows(rows []exceptionDb) {
	con := db.ConnectDB()
	defer con.Close()
	var query = `INSERT INTO krasecology.dbo.Rpt_SKAT_Excludes (SensorId, DateStart, DateEnd, Comment) VALUES (@p1, @p2, @p3, @p4) ;`
	for _, row := range rows {
		_, err := con.Exec(query, row.sensorId, row.dateStart, row.dateEnd, row.comment)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func ParserExcelExcludes(fileName string, year int) ([]Exception, error) {
	wb, err := xlsx.OpenFile(fileName) // Считываем файл
	if err != nil {
		return nil, errors.New("Ошибка считывания файла ")
	}
	sh, ok := wb.Sheet[strconv.Itoa(year)]
	if !ok {
		return nil, errors.New("Такой лист не найден ")
	}
	return getDataFromExcel(sh)
}

func getDataFromExcel(sh *xlsx.Sheet) ([]Exception, error) {
	var exceptions []Exception
	var err error
	posts := post.GetPosts()
	var cell *xlsx.Cell
	sensorTypeMaindevices := device.GetDevices()
	for i := 3; i < 1000; i++ {
		cell, err = sh.Cell(i, 3)
		if err != nil {
			log.Println(err)
			return exceptions, err
		}
		if cell.String() == "" {
			continue
		}
		if !strings.Contains(cell.String(), ".") {
			continue
		}
		ps := strings.Split(cell.String(), ".")
		var post post.Post
		if len(ps) > 1 {
			post, err = posts.FindPost(strings.Trim(ps[1], " "))
			if err != nil {
				return exceptions, errors.New("Ошибка при поиске поста: " + err.Error())
			}
		}
		cell, err = sh.Cell(i, 4)
		if err != nil {
			return exceptions, err
		}
		deviceName := cell.String()
		cell, err = sh.Cell(i, 6)
		if err != nil {
			return exceptions, err
		}
		if cell.String() == "" {
			continue
		}
		dateStart, err := cell.GetTime(false)
		if err != nil {
			return nil, errors.New("Ошибка получения даты начала простоя: " + err.Error())
		}
		cell, err = sh.Cell(i, 7)
		if err != nil {
			return exceptions, err
		}
		if cell.String() == "" {
			continue
		}

		dateEnd, err := cell.GetTime(false)
		if err != nil {
			return nil, errors.New("Ошибка получения даты конца простоя: " + err.Error())
		}
		for j := 0; j < len(sensorTypeMaindevices); j++ {
			if sensorTypeMaindevices[j].Name == deviceName {
				for _, sensorType := range sensorTypeMaindevices[j].SensorTypes {
					exceptions = append(exceptions, Exception{
						Id:         0,
						PostId:     post.Id,
						SensorType: sensorType.Id,
						DateStart:  dateStart.Add(1*time.Minute + 1*time.Second).Format("2006-01-02T15:04"),
						DateEnd:    dateEnd.Add(1*time.Minute + 1*time.Second).Format("2006-01-02T15:04"),
						Comment:    "",
					})
				}
				continue
			}
		}
	}
	return exceptions, err
}
