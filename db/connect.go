package db

import (
	"AirExcludes/configs"
	"context"
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	"gitlab.com/gbh007/gojlog"
)

/*func (dat *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := dat.DB.Query(query, args...)
	return rows, err
}*/
// Подключаемся к основной базе MS SQL
func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlserver", configs.GetConfig().ConnStr)
	if err != nil {
		gojlog.Error(err)
		return db
	}
	if db.Ping() != nil {
		gojlog.Error(err)
		return db
	}
	return db
}

func InitTx(db *sql.DB) *sql.Tx {
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		gojlog.Error(err)
	}
	return tx
}
