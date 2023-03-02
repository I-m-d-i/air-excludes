package sensorType

import (
	"AirExcludes/db"
)

type SensorType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type SensorTypes []SensorType

func GetSensorTypes() SensorTypes {
	var sensorTypes SensorTypes
	con := db.ConnectDB()
	defer con.Close()
	var query = `SELECT Id, Name FROM krasecology.dbo.Svc_SensorTypes where Visible=1 order by SortOrder`
	rows, err := con.Query(query)
	if err != nil {
		return sensorTypes
	}
	defer rows.Close()
	for rows.Next() {
		var sensorType SensorType
		if err := rows.Scan(&sensorType.Id, &sensorType.Name); err != nil {
			return sensorTypes
		}
		sensorTypes = append(sensorTypes, sensorType)
	}
	return sensorTypes
}
