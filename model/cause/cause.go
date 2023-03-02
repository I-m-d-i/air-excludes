package cause

import (
	"AirExcludes/db"
	"database/sql"
	"github.com/pkg/errors"
	"log"
)

type Cause struct {
	Id        int
	Name      string
	SortOrder sql.NullInt32
}
type Causes []Cause

func GetCauses() (causes Causes) {
	con := db.ConnectDB()
	defer con.Close()
	query := `select Id, Name, SortOrder from Rpt_SKAT_Maintenance_Cause`
	rows, err := con.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var cause Cause
	for rows.Next() {
		if err = rows.Scan(&cause.Id, &cause.Name, &cause.SortOrder); err != nil {
			log.Println(err)
			return
		}
		causes = append(causes, cause)
	}
	return
}

func (c Causes) GetId(name string) (int, error) {
	if len(c) > 0 {
		for _, cause := range c {
			if cause.Name == name {
				return cause.Id, nil
			}
		}
		return 0, errors.New("Причина несправности " + name + " не найден")
	} else {
		return 0, errors.New("Список неисправностей пуст")
	}
}
