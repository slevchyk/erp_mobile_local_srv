package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/slevchyk/erp_mobile_local_srv/dbase"
	"github.com/slevchyk/erp_mobile_local_srv/models"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var cfg models.Config
var db *sql.DB
var app *firebase.App
var eLog debug.Log

type myService struct{}

const (
	//Mobile is type of data source from mobile device
	Mobile = "mobile"
	//Accounting is type of data source from main DB
	Accounting = "accounting"
)

func init() {
	var err error
	var dir string

	for k, v := range os.Args {
		if v == "-dir" && len(os.Args) > k {
			dir = os.Args[k+1]
			dir, _ = strconv.Unquote(dir)
			dir += "/"
		}
	}

	cfg, err = loadConfiguration(fmt.Sprintf("%sconfig.json", dir))
	if err != nil {
		log.Fatal("Can't load configuration file config.json", err.Error())
	}

	db, _ = dbase.ConnectDB(cfg.DB)
	dbase.InitDB(db)

	opt := option.WithCredentialsFile(fmt.Sprintf("%sfirebase-adminsdk.json", dir))
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}
	if !isIntSess {
		runService(cfg.WinService.Name, false)
		return
	}

	if len(os.Args) == 2 {

		cmd := strings.ToLower(os.Args[1])
		switch cmd {
		case "debug":
			runService(cfg.WinService.Name, true)
			return
		case "install":
			err = installService(cfg.WinService.Name, cfg.WinService.LongName, cfg.WinService.Description)
		case "remove":
			err = removeService(cfg.WinService.Name)
		case "start":
			err = startService(cfg.WinService.Name)
			log.Printf("failed to  %v", err)
		case "stop":
			err = controlService(cfg.WinService.Name, svc.Stop, svc.Stopped)
		case "pause":
			err = controlService(cfg.WinService.Name, svc.Pause, svc.Paused)
		case "continue":
			err = controlService(cfg.WinService.Name, svc.Continue, svc.Running)
		default:
			log.Fatalf("unknown command %s", cmd)
		}

		if err != nil {
			log.Fatalf("failed to %s %s: %v", cmd, cfg.WinService.Name, err)
		}

		return
	}

	webApp()
}

func webApp() {
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/api/channel", basicAuth(channelsHandler))
	http.HandleFunc("/api/token", basicAuth(tokenHandler))
	http.HandleFunc("/api/timing", basicAuth(timingHandler))
	http.HandleFunc("/api/profile", basicAuth(profileHandler))
	http.HandleFunc("/api/helpdesk", basicAuth(helpDeskHandler))
	http.HandleFunc("/api/helpdesk/processed", basicAuth(helpDeskProcessedHandler))
	http.HandleFunc("/api/paydesk", basicAuth(payDeskHandler))
	http.HandleFunc("/api/paydesk/processed", basicAuth(payDeskProcessedHandler))
	http.HandleFunc("/api/costitems", basicAuth(costItemsHandler))
	http.HandleFunc("/api/incomeitems", basicAuth(incomeItemsHandler))
	http.HandleFunc("/api/payoffices", basicAuth(payOfficesHandler))
	http.HandleFunc("/api/currency", basicAuth(currencyHandler))
	http.HandleFunc("/api/usergrants", basicAuth(userGrantsHandler))

	http.HandleFunc("/test", testHandler)

	err := http.ListenAndServe(":8822", nil)
	if err != nil {
		log.Printf("err is %s", err)
		panic(err)
	}
}

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func installService(name, lname, desc string) error {
	exepath, err := exePath()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}

	wd, err := os.Getwd()
	if err != nil {
		s.Close()
		return err
	}
	wd = strconv.Quote(wd)
	log.Println(wd)

	s, err = m.CreateService(name, exepath, mgr.Config{DisplayName: lname, Description: desc}, "-dir", wd, "is", "auto-started")
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

func removeService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

func startService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

func controlService(name string, c svc.Cmd, to svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		eLog = debug.New(name)
	} else {
		eLog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer eLog.Close()

	eLog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &myService{})
	if err != nil {
		log.Printf("%s service failed: %v", name, err)
		eLog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	eLog.Info(1, fmt.Sprintf("%s service stopped", name))
}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	go webApp()

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				eLog.Info(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	return
}

func loadConfiguration(file string) (models.Config, error) {
	var config models.Config

	cfgFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return config, err
	}

	err = json.Unmarshal(cfgFile, &config)
	if err != nil {
		log.Println(err)
		return config, err
	}

	return config, nil
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "eSoFVzw-C9s:APA91bF0IGHuiCymYX--l4lJDbwDiQw7XaSoyDHIkBRfz4v-5tzdeSBhXqFS0bmYNRP61J5w3kGRlf_A8-OiyaSVoKsW5_69p6_zC2MA4ypufNYXqMxBxtbROB-STv7LCfqAPoXlwrhN"

	notification := messaging.Notification{
		Title:    "From Go webserver",
		Body:     "Hello world! Happy Valentines Day",
		ImageURL: "https://previews.123rf.com/images/worldofvector/worldofvector1902/worldofvector190200013/117068035-happy-valentines-day-sign-symbol-red-heart-icon-cute-graphic-object-flat-design-style-love-greeting-.jpg",
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"click_action":      "FLUTTER_NOTIFICATION_CLICK",
			"notification_type": "channel",
			"id":                "6",
			"user_id":           "",
			"type":              "message",
			"update_id":         "6",
			"title":             "From Go webserver",
			"news":              "Hello world! Happy Valentines Day",
		},
		Notification: &notification,
		Token:        registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Good job"))

}

func basicAuth(pass func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !validate(pair[0], pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func validate(username, password string) bool {
	if username == cfg.Auth.User && password == cfg.Auth.Password {
		return true
	}
	return false
}

func channelsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		channelPost(w, r)
	} else if r.Method == http.MethodGet {
		channelGet(w, r)
	}
}

func channelPost(w http.ResponseWriter, r *http.Request) {

	var cs []models.Channel
	var err error

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bs, &cs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, c := range cs {

		rows, err := dbase.SelectChannelById(db, c.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}

		if rows.Next() {
			_, err = dbase.UpdateChannel(db, c)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				continue
			}
		} else {
			_, err = dbase.InsertChannel(db, c)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				continue
			}
		}
		rows.Close()

		if c.UserID != "" {

			rows, err := dbase.SelectFirebaseTokenByUserId(db, c.UserID)
			if err != nil {
				continue
			}

			if !rows.Next() {
				rows.Close()
				continue
			}

			var ft models.FirebaseTokens
			err = dbase.ScanFirebaseToken(rows, &ft)
			if err != nil {
				rows.Close()
				continue
			}

			rows.Close()

			sendChannelNotification(c, ft.Token)
		} else {
			//to all users
			sendChannelNotification(c, "")
		}
	}

	w.WriteHeader(http.StatusOK)
}

func sendChannelNotification(c models.Channel, token string) {
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("error getting Messaging client: %v\n", err)
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"click_action":      "FLUTTER_NOTIFICATION_CLICK",
			"notification_type": "channel",
			"id":                strconv.Itoa(c.ID),
			"user_id":           c.UserID,
			"type":              c.Type,
			"update_id":         strconv.Itoa(c.UpdateID),
			"title":             c.Title,
			"news":              c.News,
		},
		Notification: &messaging.Notification{
			Title:    c.Title,
			Body:     c.News,
			ImageURL: "",
		},
		Token: token,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}

func channelGet(w http.ResponseWriter, r *http.Request) {

	var cs []models.Channel
	var c models.Channel
	var err error

	fvUserID := r.FormValue("userid")
	fvUpdateID := r.FormValue("offset")

	updateID, err := strconv.Atoi(fvUpdateID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := dbase.SelectChannelsByUserIdUpdateId(db, fvUserID, updateID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err = dbase.ScanChannel(rows, &c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cs = append(cs, c)
	}
	rows.Close()

	result, err := json.Marshal(cs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var ft models.FirebaseTokens
		var err error

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(bs, &ft)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := dbase.SelectFirebaseTokenByUserIdToken(db, ft.UserID, ft.Token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			rows.Close()
			w.WriteHeader(http.StatusOK)
			return
		}

		rows, err = dbase.SelectFirebaseTokenByUserId(db, ft.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			_, err = dbase.UpdateFirebaseTokens(db, ft)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			_, err = dbase.InsertToken(db, ft)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		rows.Close()

		w.WriteHeader(http.StatusOK)
	}
}

func timingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		timingPost(w, r)
	} else if r.Method == http.MethodGet {
		timingGet(w, r)
	}
}

func timingPost(w http.ResponseWriter, r *http.Request) {

	var ts []models.Timing
	var err error

	// параметр який вказує наам з якого об'єкта прийшли дані
	// це може бути або мобільний пристрій або облікова система
	fvFrom := r.FormValue("from")

	// перевіримо чи параметр from відповідає одному з двох дозволених значень
	if fvFrom != Mobile && fvFrom != Accounting {
		http.Error(w, "incorrect \"from\" param", http.StatusBadRequest)
		return
	}

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в масив об'єктів типу Timing
	err = json.Unmarshal(bs, &ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// переребираємо всі елементи масиву об'єктів типу Timing
	// проводимо синхронізацію з базою даних\
	for k := range ts {

		// отримуємо вказівник на елемент масиву
		// щоб була можливість модифікувати його під час синхроназації
		// t := &v
		t := &ts[k]
		if t.ID != 0 {
			rows, err := dbase.SelectTimingById(db, t.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !rows.Next() {
				rows.Close()
				http.Error(w, fmt.Sprintf("timing with id=%v not found\n", t.ID), http.StatusBadRequest)
				return
			}

			var et models.Timing
			err = dbase.ScanTiming(rows, &et)
			rows.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			t.MobID = et.MobID
			t.AccID = et.AccID
			if t.UpdatedAt.Time.After(et.UpdatedAt.Time) {
				_, err = dbase.UpdateTiming(db, *t)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				t.Status = et.Status
				t.IsTurnstile = et.IsTurnstile
				t.StartedAt = et.StartedAt
				t.EndedAt = et.EndedAt
				t.CreatedAt = et.CreatedAt
				t.UpdatedAt = et.UpdatedAt
				t.DeletedAt = et.DeletedAt
			}
		} else if fvFrom == Mobile {
			rows, err := dbase.SelectTimingByMobIdUserIdDate(db, t.MobID, t.UserID, t.Date.Time)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if rows.Next() {
				var et models.Timing
				err = dbase.ScanTiming(rows, &et)
				rows.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				t.ID = et.ID
				t.AccID = et.AccID
				if t.UpdatedAt.Time.After(et.UpdatedAt.Time) {
					_, err = dbase.UpdateTiming(db, *t)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else {
					t.Status = et.Status
					t.IsTurnstile = et.IsTurnstile
					t.StartedAt = et.StartedAt
					t.EndedAt = et.EndedAt
					t.CreatedAt = et.CreatedAt
					t.UpdatedAt = et.UpdatedAt
					t.DeletedAt = et.DeletedAt
				}
			} else {
				rows.Close()
				t.ID, err = dbase.InsertTiming(db, *t)
				if err != nil {

					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		} else if fvFrom == Accounting {
			rows, err := dbase.SelectTimingByAccIdUerIdDate(db, t.AccID, t.UserID, t.Date.Time)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if rows.Next() {
				var et models.Timing
				err = dbase.ScanTiming(rows, &et)
				rows.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				t.ID = et.ID
				t.MobID = et.MobID
				if t.UpdatedAt.Time.After(et.UpdatedAt.Time) {
					_, err = dbase.UpdateTiming(db, *t)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else {
					t.Status = et.Status
					t.IsTurnstile = et.IsTurnstile
					t.StartedAt = et.StartedAt
					t.EndedAt = et.EndedAt
					t.CreatedAt = et.CreatedAt
					t.UpdatedAt = et.UpdatedAt
					t.DeletedAt = et.DeletedAt
				}
			} else {
				rows.Close()
				t.ID, err = dbase.InsertTiming(db, *t)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	response := map[string][]models.Timing{"timing": ts}
	bs, err = json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bs)
	w.WriteHeader(http.StatusOK)
}

func timingGet(w http.ResponseWriter, r *http.Request) {

	var ts []models.Timing
	var t models.Timing
	var err error

	fvType := r.FormValue("type")

	switch fvType {
	case "dateuser":
		{
			fvDate := r.FormValue("date")
			fvUserID := r.FormValue("userid")

			var errMessage string

			if fvDate == "" {
				errMessage += fmt.Sprintf("empty \"date\" param\n")
			}

			var date time.Time
			date, err := time.Parse("20060102", fvDate)
			if err != nil {
				errMessage += fmt.Sprintf("incorrect \"date\" param format\n")
			}

			if fvUserID == "" {
				errMessage += fmt.Sprintf("empty \"user id\" param\n")
			}

			if errMessage != "" {
				http.Error(w, errMessage, http.StatusBadRequest)
				return
			}

			rows, err := dbase.SelectTimingByUserIdDate(db, fvUserID, date)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			for rows.Next() {
				err = dbase.ScanTiming(rows, &t)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				ts = append(ts, t)
			}
			rows.Close()
		}
	case "updatedat":
		{
			fvDate := r.FormValue("date")

			var errMessage string

			if fvDate == "" {
				errMessage += fmt.Sprintf("empty \"date\" param\n")
			}

			var date time.Time
			date, err := time.Parse("2006-01-02T15:04:05", fvDate)
			if err != nil {
				errMessage += fmt.Sprintf("incorrect \"date\" param format\n")
			}

			if errMessage != "" {
				http.Error(w, errMessage, http.StatusBadRequest)
				return
			}

			rows, err := dbase.SelectTimingByUpdatedAt(db, date)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			for rows.Next() {
				err = dbase.ScanTiming(rows, &t)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				ts = append(ts, t)
			}
			rows.Close()
		}
	default:
		http.Error(w, "incorrect \"type\" param", http.StatusBadRequest)
		return
	}

	response := map[string][]models.Timing{"timing": ts}
	bs, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(bs)
	w.WriteHeader(http.StatusOK)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		profileGet(w, r)
	} else if r.Method == http.MethodPost {
		profilePost(w, r)
	}

}

func profileGet(w http.ResponseWriter, r *http.Request) {

	var err error
	var p models.Profile

	fvPhone := r.FormValue("phone")
	fvPin := r.FormValue("pin")

	var errMessage string

	if fvPhone == "" {
		errMessage += fmt.Sprintf("param \"phone\" is empty\n")
	}

	if fvPin == "" {
		errMessage += fmt.Sprintf("param \"pin\" is empty\n")
	}

	if errMessage != "" {
		http.Error(w, errMessage, http.StatusBadRequest)
		return
	}

	fvPhone = "+" + fvPhone

	rows, err := dbase.SelectProfileByPhonePin(db, fvPhone, fvPin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !rows.Next() {
		rows.Close()
		http.Error(w, "can't find any profile by this phone and pin", http.StatusNotFound)
		return
	}

	err = dbase.ScanProfile(rows, &p)
	rows.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func profilePost(w http.ResponseWriter, r *http.Request) {

	var p models.Profile
	var err error

	fvFrom := r.FormValue("from")

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в об'єкт типу Profile
	err = json.Unmarshal(bs, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.UserID == "" {
		http.Error(w, "user id can not be empty", http.StatusInternalServerError)
		return
	}

	rows, err := dbase.SelectProfileByUserID(db, p.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows.Next() {
		var existProfile models.Profile
		err := dbase.ScanProfile(rows, &existProfile)
		rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p.ID = existProfile.ID
		_, err = dbase.UpdateProfile(db, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		rows.Close()
		_, err = dbase.InsertProfile(db, p)
		if err != nil {
			rows.Close()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if fvFrom == "accounting" {

		cu := map[string]string{
			"phone": p.Phone,
			"pin":   p.Pin}

		bs, err := json.Marshal(cu)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hc := &http.Client{}
		URL := cfg.MainSrv + "/api/clouddbuser"
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(bs))
		req.SetBasicAuth(cfg.MainAuth.User, cfg.MainAuth.Password)
		resp, err := hc.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if resp.StatusCode != http.StatusOK {
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Error(w, string(bs), resp.StatusCode)
			return
		}

	}

	w.WriteHeader(http.StatusOK)
}

func helpDeskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		helpDeskPost(w, r)
	} else if r.Method == http.MethodGet {
		helpDeskGet(w, r)
	}

}

func helpDeskPost(w http.ResponseWriter, r *http.Request) {

	var err error
	var hd models.HelpDesk

	fvFrom := r.FormValue("from")

	// перевіримо чи параметр from відповідає одному з двох дозволених значень
	if fvFrom == Mobile {
		hd.IsModifiedByMob = true
		hd.IsModifiedByAcc = false
	} else if fvFrom == Accounting {
		hd.IsModifiedByAcc = true
	} else {
		http.Error(w, "incorrect \"from\" param", http.StatusBadRequest)
		return
	}

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в об'єкт типу HelpDesk
	err = json.Unmarshal(bs, &hd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if hd.ID == 0 {
		id, err := dbase.InsertHelpDesk(db, hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		answer := map[string]int64{
			"id": id}

		bs, err := json.Marshal(answer)

		_, err = w.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		return
	}

	if hd.IsModifiedByMob {
		_, err = dbase.UpdateHelpDesk(db, hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var existHelpDesk models.HelpDesk

		rows, err := dbase.SelectHelpDeskByID(db, hd.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !rows.Next() {
			rows.Close()
			http.Error(w, "no such helpdesk", http.StatusBadRequest)
			return
		}

		err = dbase.ScanHelpDesk(rows, &existHelpDesk)
		rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if existHelpDesk.IsModifiedByMob {
			http.Error(w, "object already has been updated by mobile database", http.StatusBadRequest)
			return
		}

		_, err = dbase.UpdateHelpDesk(db, hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

}

func helpDeskGet(w http.ResponseWriter, r *http.Request) {

	var err error
	var hd models.HelpDesk

	fvID := r.FormValue("id")
	fvFor := r.FormValue("for")

	if fvID == "" && fvFor == "" {
		http.Error(w, "incorrect params", http.StatusBadRequest)
		return
	}

	if fvID != "" {

		id, err := strconv.Atoi(fvID)
		if err != nil {
			http.Error(w, "incorrect \"id\" param", http.StatusBadRequest)
			return
		}

		rows, err := dbase.SelectHelpDeskByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !rows.Next() {
			rows.Close()
			http.Error(w, "cant find help_desk by this id", http.StatusBadRequest)
			return
		}

		err = dbase.ScanHelpDesk(rows, &hd)
		rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bs, err := json.Marshal(hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	} else if fvFor != "" {

		var hds []models.HelpDesk
		var rows *sql.Rows

		if fvFor == Mobile {

			fvUserID := r.FormValue("userid")
			if fvUserID == "" {
				http.Error(w, "incorrect \"userid\" param", http.StatusBadRequest)
				return
			}

			rows, err = dbase.SelectHelpDeskModifiedByAcc(db, fvUserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if fvFor == Accounting {
			rows, err = dbase.SelectHelpDeskModifiedByMob(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "incorrect \"for\" param", http.StatusBadRequest)
		}

		for rows.Next() {
			err = dbase.ScanHelpDesk(rows, &hd)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			hds = append(hds, hd)
		}
		rows.Close()

		bs, err := json.Marshal(hds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "bad params", http.StatusBadRequest)
}

func helpDeskProcessedHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	var hd models.HelpDesk

	fvFrom := r.FormValue("from")
	fvID := r.FormValue("id")

	// перевіримо чи параметр from відповідає одному з двох дозволених значень
	if fvFrom == Mobile {
		hd.IsModifiedByAcc = false
	} else if fvFrom == Accounting {
		hd.IsModifiedByMob = false
	} else {
		http.Error(w, "incorrect \"from\" param", http.StatusBadRequest)
		return
	}

	if fvID == "" {
		http.Error(w, "incorrect \"id\" param", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(fvID)
	if err != nil {
		http.Error(w, "non integer \"id\" param", http.StatusBadRequest)
		return
	}

	rows, err := dbase.SelectHelpDeskByID(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows.Next() {
		err = dbase.ScanHelpDesk(rows, &hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = dbase.UpdateHelpDesk(db, hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rows.Close()

	w.WriteHeader(http.StatusOK)
}
func payDeskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		payDeskPost(w, r)
	} else if r.Method == http.MethodGet {
		payDeskGet(w, r)
	}

}

func payDeskPost(w http.ResponseWriter, r *http.Request) {
	var err error
	var pd models.PayDesk

	fvFrom := r.FormValue("from")

	// перевіримо чи параметр from відповідає одному з двох дозволених значень
	if fvFrom == Mobile {
		pd.IsModifiedByMob = true
		pd.IsModifiedByAcc = false
	} else if fvFrom == Accounting {
		pd.IsModifiedByAcc = true
	} else {
		http.Error(w, "incorrect \"from\" param", http.StatusBadRequest)
		return
	}

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в об'єкт типу PayDesk
	err = json.Unmarshal(bs, &pd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if pd.ID == 0 {
		id, err := dbase.InsertPayDesk(db, pd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		answer := map[string]int64{
			"id": id}

		bs, err := json.Marshal(answer)

		_, err = w.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		return
	}

	if pd.IsModifiedByMob {
		_, err = dbase.UpdatePayDesk(db, pd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	var existPayDesk models.PayDesk

	rows, err := dbase.SelectPayDeskByID(db, pd.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !rows.Next() {
		rows.Close()
		http.Error(w, "no such paydesk", http.StatusBadRequest)
		return
	}

	err = dbase.ScanPayDesk(rows, &existPayDesk)
	rows.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existPayDesk.IsModifiedByMob {
		http.Error(w, "object already has been updated by mobile database", http.StatusNoContent)
		return
	}

	existPayDesk.IsDeleted = pd.IsDeleted
	existPayDesk.Amount = pd.Amount
	existPayDesk.Payment = pd.Payment
	existPayDesk.DocumentDate = pd.DocumentDate
	existPayDesk.DocumentNumber = pd.DocumentNumber
	existPayDesk.IsModifiedByMob = pd.IsModifiedByMob
	existPayDesk.IsModifiedByAcc = pd.IsModifiedByAcc

	_, err = dbase.UpdatePayDesk(db, existPayDesk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func payDeskGet(w http.ResponseWriter, r *http.Request) {

	var err error
	var pd models.PayDesk

	fvID := r.FormValue("id")
	fvFor := r.FormValue("for")

	if fvID == "" && fvFor == "" {
		http.Error(w, "incorrect params", http.StatusBadRequest)
		return
	}

	if fvID != "" {

		id, err := strconv.Atoi(fvID)
		if err != nil {
			http.Error(w, "incorrect \"id\" param", http.StatusBadRequest)
			return
		}

		rows, err := dbase.SelectPayDeskByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !rows.Next() {
			rows.Close()
			http.Error(w, "cant find pay_desk by this id", http.StatusBadRequest)
			return
		}

		err = dbase.ScanPayDesk(rows, &pd)
		rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bs, err := json.Marshal(pd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	} else if fvFor != "" {

		var hds []models.PayDesk
		var rows *sql.Rows

		if fvFor == Mobile {

			fvUserID := r.FormValue("userid")
			if fvUserID == "" {
				http.Error(w, "incorrect \"userid\" param", http.StatusBadRequest)
				return
			}

			rows, err = dbase.SelectPayDeskModifiedByAcc(db, fvUserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if fvFor == Accounting {
			rows, err = dbase.SelectPayDeskModifiedByMob(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "incorrect \"for\" param", http.StatusBadRequest)
		}

		for rows.Next() {
			err = dbase.ScanPayDesk(rows, &pd)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			hds = append(hds, pd)
		}
		rows.Close()

		bs, err := json.Marshal(hds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(bs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "bad params", http.StatusBadRequest)

}

func payDeskProcessedHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	var pd models.PayDesk

	fvFrom := r.FormValue("from")
	fvID := r.FormValue("id")

	// перевіримо чи параметр from відповідає одному з двох дозволених значень
	if fvFrom != Mobile && fvFrom != Accounting {
		http.Error(w, "incorrect \"from\" param", http.StatusBadRequest)
		return
	}

	if fvID == "" {
		http.Error(w, "incorrect \"id\" param", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(fvID)
	if err != nil {
		http.Error(w, "non integer \"id\" param", http.StatusBadRequest)
		return
	}

	rows, err := dbase.SelectPayDeskByID(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows.Next() {
		err = dbase.ScanPayDesk(rows, &pd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if fvFrom == Mobile {
			pd.IsModifiedByAcc = false
		} else if fvFrom == Accounting {
			pd.IsModifiedByMob = false
		}

		_, err = dbase.UpdatePayDesk(db, pd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rows.Close()

	w.WriteHeader(http.StatusOK)
}

func costItemsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		costItemsPost(w, r)
	} else if r.Method == http.MethodGet {
		costItemsGet(w, r)
	}

}

func costItemsPost(w http.ResponseWriter, r *http.Request) {
	var err error
	var ci models.CostItem
	var cis []models.CostItem

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в об'єкт типу []CostItem
	err = json.Unmarshal(bs, &cis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range cis {

		rows, err := dbase.SelectCostItemsByAccID(db, v.AccID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			err = dbase.ScanCostItem(rows, &ci)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			v.ID = ci.ID
			v.CreatedAt = ci.CreatedAt
			v.UpdatedAt.Valid = true
			v.UpdatedAt.Time = time.Now()

			_, err = dbase.UpdateCostItem(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			v.CreatedAt.Valid = true
			v.CreatedAt.Time = time.Now()
			_, err = dbase.InsertCostItem(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		rows.Close()
	}

	w.WriteHeader(http.StatusOK)
}

func costItemsGet(w http.ResponseWriter, r *http.Request) {

	var err error
	var ci models.CostItem
	var cis []models.CostItem

	rows, err := dbase.SelectCostItems(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err = dbase.ScanCostItem(rows, &ci)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cis = append(cis, ci)
	}
	rows.Close()

	bs, err := json.Marshal(cis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func incomeItemsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		incomeItemsPost(w, r)
	} else if r.Method == http.MethodGet {
		incomeItemsGet(w, r)
	}

}

func incomeItemsPost(w http.ResponseWriter, r *http.Request) {
	var err error
	var ii models.IncomeItem
	var iis []models.IncomeItem

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в об'єкт типу []IncomeItem
	err = json.Unmarshal(bs, &iis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range iis {

		rows, err := dbase.SelectIncomeItemsByAccID(db, v.AccID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			err = dbase.ScanIncomeItem(rows, &ii)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			v.ID = ii.ID
			v.CreatedAt = ii.CreatedAt
			v.UpdatedAt.Valid = true
			v.UpdatedAt.Time = time.Now()

			_, err = dbase.UpdateIncomeItem(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			v.CreatedAt.Valid = true
			v.CreatedAt.Time = time.Now()
			_, err = dbase.InsertIncomeItem(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		rows.Close()
	}

	w.WriteHeader(http.StatusOK)
}

func incomeItemsGet(w http.ResponseWriter, r *http.Request) {
	var err error
	var ii models.IncomeItem
	var iis []models.IncomeItem

	rows, err := dbase.SelectIncomeItems(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err = dbase.ScanIncomeItem(rows, &ii)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		iis = append(iis, ii)
	}
	rows.Close()

	bs, err := json.Marshal(iis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func payOfficesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		payOfficesPost(w, r)
	} else if r.Method == http.MethodGet {
		payOfficesGet(w, r)
	}

}

func payOfficesPost(w http.ResponseWriter, r *http.Request) {
	var err error
	var po models.PayOffice
	var pos []models.PayOffice

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в об'єкт типу []PayOffice
	err = json.Unmarshal(bs, &pos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range pos {

		rows, err := dbase.SelectPayOfficesByAccID(db, v.AccID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			err = dbase.ScanPayOffice(rows, &po)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			v.ID = po.ID
			v.CreatedAt = po.CreatedAt
			v.UpdatedAt.Valid = true
			v.UpdatedAt.Time = time.Now()

			_, err = dbase.UpdatePayOffice(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			v.CreatedAt.Valid = true
			v.CreatedAt.Time = time.Now()
			_, err = dbase.InsertPayOffice(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		rows.Close()
	}

	w.WriteHeader(http.StatusOK)
}

func payOfficesGet(w http.ResponseWriter, r *http.Request) {
	var err error
	var po models.PayOffice
	var pos []models.PayOffice

	rows, err := dbase.SelectPayOffices(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err = dbase.ScanPayOffice(rows, &po)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pos = append(pos, po)
	}
	rows.Close()

	bs, err := json.Marshal(pos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func currencyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		currencyPost(w, r)
	} else if r.Method == http.MethodGet {
		currencyGet(w, r)
	}

}

func currencyPost(w http.ResponseWriter, r *http.Request) {
	var err error
	var c models.Currency
	var cs []models.Currency

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в об'єкт типу []PayOffice
	err = json.Unmarshal(bs, &cs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range cs {

		rows, err := dbase.SelectCurrencyByAccID(db, v.AccID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			err = dbase.ScanCurrency(rows, &c)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			v.ID = c.ID
			v.CreatedAt = c.CreatedAt
			v.UpdatedAt.Valid = true
			v.UpdatedAt.Time = time.Now()

			_, err = dbase.UpdateCurrency(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			v.CreatedAt.Valid = true
			v.CreatedAt.Time = time.Now()
			_, err = dbase.InsertCurrency(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		rows.Close()
	}

	w.WriteHeader(http.StatusOK)
}

func currencyGet(w http.ResponseWriter, r *http.Request) {
	var err error
	var c models.Currency
	var cs []models.Currency

	rows, err := dbase.SelectCurrency(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err = dbase.ScanCurrency(rows, &c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cs = append(cs, c)
	}
	rows.Close()

	bs, err := json.Marshal(cs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func userGrantsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		userGrantsPost(w, r)
	} else if r.Method == http.MethodGet {
		userGrantsGet(w, r)
	}

}

func userGrantsPost(w http.ResponseWriter, r *http.Request) {
	var err error
	var ug models.UserGrants
	var ugs []models.UserGrants

	// зчитуємо тіло запиту
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// конвертуэмо масив байтів в об'єкт типу []UserGrants
	err = json.Unmarshal(bs, &ugs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, v := range ugs {

		rows, err := dbase.SelectUserGrantsByUserIDObject(db, v.UserID, v.ObjectType, v.ObjectAccID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows.Next() {
			err = dbase.ScanUserGrants(rows, &ug)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			v.UserID = ug.UserID
			v.ObjectType = ug.ObjectType

			_, err = dbase.UpdateUserGrants(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			_, err = dbase.InsertUserGrants(db, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		rows.Close()
	}

	w.WriteHeader(http.StatusOK)
}

func userGrantsGet(w http.ResponseWriter, r *http.Request) {
	var err error
	var ug models.UserGrants
	var ugs []models.UserGrants

	fvUserID := r.FormValue("userid")

	if fvUserID == "" {
		http.Error(w, "incorrect \"userid\" param", http.StatusBadRequest)
		return
	}

	rows, err := dbase.SelectUserGrantsByUserID(db, fvUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		err = dbase.ScanUserGrants(rows, &ug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ugs = append(ugs, ug)
	}
	rows.Close()

	bs, err := json.Marshal(ugs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
