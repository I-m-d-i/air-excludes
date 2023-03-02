package main

import (
	"AirExcludes/configs"
	"AirExcludes/controller"
	"flag"
	"fmt"
	"gitlab.com/gbh007/gojlog"
	_ "gitlab.com/krasecology/go-lib"
	model "gitlab.com/krasecology/go-lib/db/model"
	"log"
)

func main() {
	staticDir := flag.String("s", "static", "папка с файлами для раздачи веб сервера")
	webPort := flag.Int("p", 80, "порт веб сервера")
	flag.Parse()
	if configs.Load() != nil {
		log.Println("ошибка загрузки конфигурации")
		return
	}
	cnf := configs.GetConfig()
	// подключаем библиотеку к базе авторизации
	if err := model.Connect(
		cnf.Auth.Mongo.URI,
		cnf.Auth.Mongo.DBName,
		cnf.Auth.Mongo.User,
		cnf.Auth.Mongo.Password,
	); err != nil {
		gojlog.Critical("ошибка соединения с базой библиотеки")
		return
	}
	done := controller.Run(fmt.Sprintf(":%d", *webPort), *staticDir)
	gojlog.Info("Сервер запущен")
	<-done
}
