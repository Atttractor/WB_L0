package main

// Кароче нахуй там в файле postresql комменты перед функциями и ещё тут у тебя всё нахуй упало. А почему? А я не ебу почему.
// Тебе надо сделать чтобы сервак работал вместе с подпиской на канал nats-streaming и добавить запись и чтение из бд, а ещё посмотреть
// как проводить валидацию данных и где вообще её провести, тут в мейне или уже в storage или и там и там. А ищо надо штобы функция New() из
// postgresql.go работала нармальна, а не как залупа(Сейчас почему-то не заполняется список items в структуре model).

import (
	"Subscriber/internal/model"
	"Subscriber/internal/storage/postgresql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/nats-io/stan.go"
)

func main() {
	var (
		clusterID string
		clientID  string
		URL       string
		myModel   model.MyModel
	)

	// Storage
	mapa := postgresql.New()

	clusterID = "test-cluster"
	clientID = "myID"

	// Server
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("uid")

		m1 := mapa[uid]

		a, err := json.MarshalIndent(m1, "", "	")
		if err != nil {
			fmt.Println(err)
		}

		w.Write(a)
	})
	http.ListenAndServe(":80", nil)
	log.Println("Server started")

	// Subscriber
	sc, err := stan.Connect(clusterID, clientID,
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	log.Printf("Connected to clusterID: [%s] clientID: [%s]\n", clusterID, clientID)

	subj := "foo"
	sub, err := sc.Subscribe(subj, func(m *stan.Msg) {
		json.Unmarshal(m.Data, &myModel)
		fmt.Println(myModel)
		mapa[myModel.Order_uid] = myModel
	})
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s]\n", subj, clientID)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sub.Unsubscribe()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
