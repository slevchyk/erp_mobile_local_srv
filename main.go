package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/slevchyk/erp_mobile_local_srv/dbase"
	"github.com/slevchyk/erp_mobile_local_srv/models"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var cfg models.Config
var db *sql.DB
var app *firebase.App

func init() {
	var err error

	cfg = models.Config{
		Auth: models.AuthConfig{
			User:     "mobile",
			Password: "Dq4fS^J&^nqQ(fg4",
		},
		DB:   models.DBConfig{
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
	http.HandleFunc("/test", testHandler)

	err := http.ListenAndServe(":7133", nil)
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
		Title:    "Title go",
		Body:     "New form go webserver",
		ImageURL: "",
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"page": "channel",
			"time":  "2:45",
			"title": notification.Title,
			"news": notification.Body,
		},
		Notification: &notification,
		Token: registrationToken,
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
		var c models.Channel
		var err error

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(bs, &c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = dbase.InsertChannel(db, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Obtain a messaging.Client from the App.
		ctx := context.Background()
		client, err := app.Messaging(ctx)
		if err != nil {
			log.Printf("error getting Messaging client: %v\n", err)
		}

		// This registration token comes from the client FCM SDKs.
		registrationToken := "fBdPzDgC9eM:APA91bFZ97xXxQr6jUZpY4Cn2VXp0BnlhNg9loAHlAEI_MLbdgP8H-txBEBZoLXRJ84CvQfwlnel7udFHXSxQssPzWkwwpq817euuYPR0gVMoPYCxG9qIN7JzCnmoKqKIElWcHQzHsbC"

		// See documentation on defining a message payload.
		message := &messaging.Message{
			Data: map[string]string{
				"key": "FLUTTER_NOTIFICATION_CLICK",
				"score": "850",
				"time":  "2:45",
			},
			Notification: &messaging.Notification{
				Title:    "Title",
				Body:     "New form go webserver",
				ImageURL: "",
			},
			Token: registrationToken,
		}

		// Send a message to the device corresponding to the provided
		// registration token.
		response, err := client.Send(ctx, message)
		if err != nil {
			log.Fatalln(err)
		}
		// Response is a message ID string.
		fmt.Println("Successfully sent message:", response)

		w.WriteHeader(http.StatusOK)
	}

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

		rows, err := dbase.SelectFirebaseTokenByUserIdToken(db, ft.UserID, ft.Token);
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