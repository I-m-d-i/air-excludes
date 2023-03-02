package controller

import (
	"AirExcludes/configs"
	maintenance "AirExcludes/model/Skat_maintenance"
	"AirExcludes/model/excludes"
	"AirExcludes/model/post"
	"AirExcludes/model/sensorType"
	"encoding/json"
	"fmt"
	"gitlab.com/gbh007/gojlog"
	"gitlab.com/krasecology/go-lib/account"
	"gitlab.com/krasecology/go-lib/web"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Run(addr string, staticDir string) <-chan struct{} {
	gojlog.Info("Заходим в Run")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	web.AddHandler(mux, "/api/getExceptions", "admin-menu:excludes:r", getExceptions())
	web.AddHandler(mux, "/api/getPosts", "", getPosts())
	web.AddHandler(mux, "/api/getSensorTypes", "", getSensorTypes())
	web.AddHandler(mux, "/api/uploadFile", "admin-menu:excludes:e", uploadExel())
	web.AddHandler(mux, "/api/saveExceptions", "admin-menu:excludes:e", saveExceptions())
	web.AddHandler(mux, "/api/auth/login", "", login())
	web.AddHandler(mux, "/api/auth/logout", "", web.Logout())
	web.AddHandler(mux, "/api/auth/authRequire", "", web.AuthRequireHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {})))
	// создание объекта сервера
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  1 * time.Minute,
	}
	done := make(chan struct{})
	go func() {
		if err := server.ListenAndServe(); err != nil {
			gojlog.Error(err)
		}
		close(done)
	}()
	return done
}
func uploadExel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST, OPTIONS")
			return
		}
		w.Header().Set("Content-Type", "multipart/form-data")
		w.Header().Set("Content-Type", "application/json")
		contentType := strings.Split(r.Header.Get("Content-Type"), ";")
		if contentType[0] != "multipart/form-data" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if err := r.ParseMultipartForm(10000); err != nil {
			gojlog.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		year, err := strconv.Atoi(r.FormValue("year"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var handler *multipart.FileHeader
		var file multipart.File
		// get 'file'
		file, handler, err = r.FormFile("file")
		if err != nil {
			log.Println(err)
			err = file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// read file bytes
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// write bytes to a localfile
		err = ioutil.WriteFile(handler.Filename, fileBytes, 0644)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response, err := excludes.ParserExcelExcludes(handler.Filename, year)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		maintenances, err := maintenance.ParserExcelMaintenance(handler.Filename, year)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = maintenance.AddMaintenance(maintenances, year)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		web.SetResponse(r, response)
	})
}
func login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			RecaptchaToken string `json:"recaptchaToken"`
			Login          string `json:"login"`
			Pass           string `json:"pass"`
		}{}
		if web.ParseJSON(r, &request) != nil {
			web.SetError(r, web.ErrParseData)
			return
		}
		message := url.Values{
			"secret":   {configs.GetConfig().SecretKeyReCaptcha},
			"response": {request.RecaptchaToken},
			"remoteip": {web.GetIP(r)},
		}
		resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", message)
		if err != nil {
			log.Fatalln(err)
			return
		}
		responseReCaptcha := struct {
			Success      bool          `json:"success"`
			Challenge_ts string        `json:"challenge_ts"`
			Hostname     string        `json:"hostname"`
			Error_codes  []interface{} `json:"error-codes"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&responseReCaptcha)
		if err != nil {
			web.SetError(r, web.ErrParseData)
			log.Println(err)
			return
		}
		if !responseReCaptcha.Success {
			web.SetError(r, fmt.Errorf("ответ на капчу не действителен"))
			log.Println("Капча не действительна IP:", r.RemoteAddr)
			return
		}
		token, ok := account.Login(request.Login, request.Pass)
		if !ok {
			web.SetError(r, fmt.Errorf("неверный логин/пароль"))
			return
		}
		web.SetResponse(r, struct{}{})
		if strings.Index(r.Header.Get("Origin"), "krasecology.ru") != -1 {
			http.SetCookie(w, &http.Cookie{
				Name:     web.CoreTokenName,
				Value:    token.Hex(),
				Path:     "/",
				HttpOnly: true,
				Domain:   "krasecology.ru",
			})
			// удаляем корневой токен на поддомене если есть
			http.SetCookie(w, &http.Cookie{
				Name:     web.CoreTokenName,
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Unix(0, 0),
			})
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:     web.CoreTokenName,
				Value:    token.Hex(),
				Path:     "/",
				HttpOnly: true,
			})
		}
		gojlog.Infof("создана сессия %s для %s ip %s", token.Hex(), request.Login, web.GetIP(r))
	})
}
func getExceptions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := excludes.GetExceptions()
		web.SetResponse(r, response)
	})
}
func getPosts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := post.GetPosts()
		web.SetResponse(r, response)
	})
}
func getSensorTypes() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := sensorType.GetSensorTypes()
		web.SetResponse(r, response)
	})
}
func saveExceptions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var updates excludes.ModifiedExceptions
		err := json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			web.SetError(r, web.ErrParseData)
			log.Println(err)
			return
		}
		err = excludes.SaveExceptions(updates)
		if err != nil {
			log.Println(err)
			web.SetError(r, web.ErrCreateData)
			return
		}
	})
}
