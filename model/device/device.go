package device

import (
	"AirExcludes/db"
	"AirExcludes/model/sensorType"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type Device struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	SensorTypes []sensorType.SensorType
}
type Devices []Device

func GetDevices() (devices Devices) {
	con := db.ConnectDB()
	defer con.Close()
	var query = `SELECT 
    	rsmd.Id,
    	rsmd.Name,
		sstmdl.Id_Svc_SensorTypes,
		st.Name
		FROM krasecology.dbo.Rpt_SKAT_Maint_Devices rsmd 
    	inner join dbo.Svc_SensorTypes_Maint_Devices_l sstmdl on rsmd.Id = sstmdl.Id_Rpt_SKAT_Maint_Devices 
		inner join dbo.Svc_SensorTypes st on st.Id = sstmdl.Id_Svc_SensorTypes
		order by rsmd.Id`
	rows, err := con.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	type deviceDB struct {
		deviceId       int
		deviceName     string
		sensorTypeId   int
		sensorTypeName string
	}
	var deviceDBRow deviceDB
	var deviceDBList []deviceDB
	for rows.Next() {
		if err = rows.Scan(&deviceDBRow.deviceId, &deviceDBRow.deviceName, &deviceDBRow.sensorTypeId, &deviceDBRow.sensorTypeName); err != nil {
			log.Println(err)
			return
		}
		deviceDBList = append(deviceDBList, deviceDBRow)
	}
	var device Device
	for _, deviceDBRow := range deviceDBList {
		if device.Id == 0 {
			device.Id = deviceDBRow.deviceId
			device.Name = strings.Trim(deviceDBRow.deviceName, " ")
			device.SensorTypes = []sensorType.SensorType{{Id: deviceDBRow.sensorTypeId, Name: deviceDBRow.sensorTypeName}}
		}
		if device.Id != deviceDBRow.deviceId {
			devices = append(devices, device)
			device.Id = deviceDBRow.deviceId
			device.Name = strings.Trim(deviceDBRow.deviceName, " ")
			device.SensorTypes = []sensorType.SensorType{{Id: deviceDBRow.sensorTypeId, Name: deviceDBRow.sensorTypeName}}

		}
	}
	devices = append(devices, device)
	return
}

func (d Devices) GetId(name string) (int, error) {
	if len(d) > 0 {
		for _, device := range d {
			if device.Name == name {
				return device.Id, nil
			}
		}
		return 0, errors.New("Прибор \"" + name + "\" не найден")
	} else {
		return 0, errors.New("Список приборов пуст")
	}
}
