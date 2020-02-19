package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

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

const (
	Mobile     = "mobile"
	Accounting = "accounting"
)

func init() {
	var err error

	cfg = models.Config{
		Auth: models.AuthConfig{
			User:     "mobile",
			Password: "Dq4fS^J&^nqQ(fg4",
		},
		DB: models.DBConfig{
			Name:     "worker_local",
			User:     "worker",
			Password: "worker",
		},
	}

	db, _ = dbase.ConnectDB(cfg.DB)
	dbase.InitDB(db)

	opt := option.WithCredentialsFile("willingwork-43b10-firebase-adminsdk-2sf2v-7600960d26.json")
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}
}

func main() {

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/api/channel", basicAuth(channelsHandler))
	http.HandleFunc("/api/token", basicAuth(tokenHandler))
	http.HandleFunc("/api/timing", basicAuth(timingHandler))
	http.HandleFunc("/test", testHandler)

	err := http.ListenAndServe(":8822", nil)
	if err != nil {
		panic(err)
	}

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

	if r.Method == "POST" {
		channelPost(w, r)
	} else if r.Method == "GET" {
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
		defer rows.Close()

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

		if c.UserID != "" {

			rows, err := dbase.SelectFirebaseTokenByUserId(db, c.UserID)
			if err != nil {
				continue
			}
			defer rows.Close()

			if !rows.Next() {
				continue
			}

			var ft models.FirebaseTokens
			err = dbase.ScanFirebaseToken(rows, &ft)
			if err != nil {
				continue
			}

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
	defer rows.Close()

	for rows.Next() {
		err = dbase.ScanChannel(rows, &c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cs = append(cs, c)
	}

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
		defer rows.Close()

		if rows.Next() {
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

		w.WriteHeader(http.StatusOK)

	}
}

func timingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		timingPost(w, r)
	}
}

func timingPost(w http.ResponseWriter, r *http.Request) {

	var ts []models.Timing
	//var t models.Timing
	var err error

	fvFrom := r.FormValue("from")

	if fvFrom != Mobile && fvFrom != Accounting {
		http.Error(w, "wrong \"from\" param", http.StatusBadRequest)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bs, &ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, t := range ts {
		if t.ID != 0 {
			rows, err := dbase.SelectTimingById(db, t.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !rows.Next() {
				http.Error(w, fmt.Sprintf("timing with id=%v not found\n", t.ID), http.StatusBadRequest)
				return
			}

			_, err = dbase.UpdateTiming(db, t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if fvFrom == Mobile {
			rows, err := dbase.SelectTimingByMobIdUserIdDate(db, t.MobID, t.UserID, t.Date)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if rows.Next() {
				var et models.Timing
				err = dbase.ScanTiming(rows, &et)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				t.ID = et.ID
				t.AccID = et.AccID
				if t.UpdatedAt > et.UpdatedAt {
					_, err = dbase.UpdateTiming(db, t)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			} else {
				t.ID, err = dbase.InsertTiming(db, t)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

		} else if fvFrom == Accounting {
			rows, err := dbase.SelectTimingByAccIdUerIdDate(db, t.AccID, t.UserID, t.Date)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !rows.Next() {
				var et models.Timing
				err = dbase.ScanTiming(rows, &et)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				t.ID = et.ID
				t.MobID = et.MobID
				if t.UpdatedAt > et.UpdatedAt {
					_, err = dbase.UpdateTiming(db, t)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			} else {
				t.ID, err = dbase.InsertTiming(db, t)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	bs, err = json.Marshal(ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bs)
	w.WriteHeader(http.StatusOK)
}
