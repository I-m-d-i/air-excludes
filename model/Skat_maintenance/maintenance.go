package maintenance

import (
	"AirExcludes/db"
	"AirExcludes/model/cause"
	"AirExcludes/model/device"
	"AirExcludes/model/post"
	"database/sql"
	"errors"
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"log"
	"strconv"
	"strings"
)

type Maintenance struct {
	PostId        int
	Target        string
	DeviceName    string
	DeviceId      int
	DeviceModel   string
	FaultStart    sql.NullTime
	FaultEnd      sql.NullTime
	FaultDuration sql.NullFloat64
	FaultName1    string
	FaultName2    string
	MaintType     string
	MaintWork     string
	Comment       string
	WorkStart     sql.NullTime
	WorkEnd       sql.NullTime
	TripStart     sql.NullTime
	Staff         string
	Year          int
	Month         string
	PostName      string
	CauseId       int
}

type Maintenances []Maintenance

func ParserExcelMaintenance(fileName string, year int) (Maintenances, error) {
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

func AddMaintenance(maintenances Maintenances, year int) error {
	fmt.Println(len(maintenances))
	con := db.ConnectDB()
	defer con.Close()
	tx := db.InitTx(con)
	_, err := tx.Exec(`delete dbo.Rpt_skat_maintenance where Year = @p1`, year)
	if err != nil {
		log.Println("Ошибка во время удаления записей из базы данных: " + err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		return errors.New("Ошибка во время удаления записей из базы данных: " + err.Error())
	}
	var query = `INSERT INTO krasecology.dbo.Rpt_skat_maintenance 
    (PostId, DeviceName, Target, DeviceId ,	DeviceModel, FaultStart,	FaultEnd,	FaultDuration,	FaultName1 , FaultName2, MaintType,
	MaintWork,	Comment,	WorkEnd,	WorkStart,	TripStart ,Staff ,	[Year],	[Month], PostName,	CauseId) 
	values (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8,@p9,@p10,@p11,@p12,@p13,@p14,@p15,@p16,@p17,@p18,@p19,@p20,@p21)`
	for _, maintenance := range maintenances {
		_, err = tx.Exec(query, maintenance.PostId, maintenance.DeviceName, maintenance.Target, maintenance.DeviceId, maintenance.DeviceModel, maintenance.FaultStart,
			maintenance.FaultEnd, maintenance.FaultDuration, maintenance.FaultName1, maintenance.FaultName2, maintenance.MaintType,
			maintenance.MaintWork, maintenance.Comment, maintenance.WorkEnd, maintenance.WorkStart, maintenance.TripStart, maintenance.Staff,
			maintenance.Year, maintenance.Month, maintenance.PostName, maintenance.CauseId)
		if err != nil {
			log.Println("Ошибка при добавлении в базу данных: " + err.Error())
			err = tx.Rollback()
			if err != nil {
				log.Println("Ошибка отката: " + err.Error())
				return err
			}
			return errors.New("Ошибка при добавлении в базу данных: " + err.Error())
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println("Ошибка во время подтверждения транзакции: " + err.Error())
		return errors.New("Ошибка во время подтверждения транзакции: " + err.Error())
	}
	return nil
}

func getDataFromExcel(sh *xlsx.Sheet) (Maintenances, error) {
	var maintenance Maintenance
	var maintenances Maintenances
	var err error
	var cell *xlsx.Cell
	var row *xlsx.Row
	posts := post.GetPosts()
	devices := device.GetDevices()

	var post post.Post
	var year int
	var deviceId int
	var causeId int
	for i := 3; ; i++ {
		row, err = sh.Row(i)
		if err != nil {
			return maintenances, errors.New("Ошибка получения строки " + strconv.Itoa(i) + " : " + err.Error())
		}
		cell = row.GetCell(0)
		if cell.Value == "" {
			break
		}
		year, err = cell.Int()
		if err != nil {
			return maintenances, errors.New("Некоректный год в строке " + strconv.Itoa(i) + " : " + err.Error())
		}
		maintenance.Year = year
		maintenance.Month = row.GetCell(1).String()
		maintenance.Target = row.GetCell(2).String()
		maintenance.PostName = row.GetCell(3).String()
		ps := strings.Split(maintenance.PostName, ".")
		if len(ps) > 1 {
			post, err = posts.FindPost(strings.Trim(ps[1], " "))
			if err != nil {
				return maintenances, errors.New("Ошибка при поиске поста: " + err.Error())
			}
			maintenance.PostId = post.Id
		}
		maintenance.DeviceName = row.GetCell(4).String()
		deviceId, err = devices.GetId(strings.Trim(maintenance.DeviceName, " "))
		if err != nil {
			return maintenances, errors.New("Ошибка при нахождении прибора " + maintenance.DeviceName + ": " + err.Error())
		}
		maintenance.DeviceId = deviceId
		maintenance.DeviceModel = row.GetCell(5).String()
		faultStart, err := row.GetCell(6).GetTime(false)
		if err != nil {
			maintenance.FaultStart.Valid = false
		} else {
			maintenance.FaultStart.Valid = true
		}
		maintenance.FaultStart.Time = faultStart
		faultEnd, err := row.GetCell(6).GetTime(false)
		if err != nil {
			maintenance.FaultEnd.Valid = false
		} else {
			maintenance.FaultEnd.Valid = true
		}
		maintenance.FaultEnd.Time = faultEnd
		durationTemp, err := row.GetCell(10).Float()
		if err != nil {
			maintenance.FaultDuration.Valid = false
		} else {
			maintenance.FaultDuration.Valid = true
		}
		maintenance.FaultDuration.Float64 = durationTemp
		maintenance.FaultName1 = row.GetCell(11).String()
		maintenance.FaultName2 = row.GetCell(12).String()
		maintenance.MaintType = row.GetCell(13).String()
		maintenance.MaintWork = row.GetCell(14).String()
		maintenance.Comment = row.GetCell(15).String()
		maintenance.Staff = row.GetCell(22).String()
		causes := cause.GetCauses()
		causeId, err = causes.GetId(row.GetCell(23).String())
		if err != nil {
			return maintenances, err
		}
		maintenance.CauseId = causeId
		maintenances = append(maintenances, maintenance)
	}
	return maintenances, err
}
