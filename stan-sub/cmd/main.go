package main

import (
	"Subscriber/internal/model"
	"Subscriber/internal/storage/postgresql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

func main() {
	var (
		clusterID string
		clientID  string
		myModel   model.MyModel
	)

	// Validator
	v := validator.New()

	// Storage
	mapa := postgresql.New()

	clusterID = "test-cluster"
	clientID = "myID"

	// Subscriber
	sc, err := stan.Connect(clusterID, clientID,
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Не удалось подключится к nats-streming-server: %v", err)
	}
	log.Printf("Успешное подключение к nats-streming-server, clusterID: [%s] clientID: [%s]\n", clusterID, clientID)

	subj := "foo"
	sub, err := sc.Subscribe(subj, func(m *stan.Msg) {
		log.Println("Новое сообщение!")
		json.Unmarshal(m.Data, &myModel)
		err := v.Struct(myModel)
        if err != nil {
            log.Fatalf("Данные не прошли валидацию: %s", err)
        }

		log.Println(myModel)
		mapa[myModel.Order_uid] = myModel

		err = postgresql.SaveData(myModel)
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Канал: [%s], clientID=[%s]\n", subj, clientID)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Printf("\nОтписка и закрытие подключения\n\n")
			sub.Unsubscribe()
			sc.Close()
			cleanupDone <- true
		}
	}()

	// Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("uid")
    
		m1 := mapa[uid]

    	ts, err := template.ParseFiles(`C:\Users\Dooplik\Desktop\WB_L0\stan-sub\templates\index.html`)
    	if err != nil {
        	log.Println(err)
    	}

    	err = ts.Execute(w, m1)
    	if err != nil {
			log.Println(err)
    	}
	})
	http.ListenAndServe(":80", nil)

	<-cleanupDone
}