package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/stan.go"
)

type Delivery struct {
	Delivery_uuid string `json:"delivery_uuid"`
	Name          string `json:"name,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Zip           string `json:"zip,omitempty"`
	City          string `json:"city,omitempty"`
	Address       string `json:"address,omitempty"`
	Region        string `json:"region,omitempty"`
	Email         string `json:"email,omitempty"`
}

type Payment struct {
	Transaction   string `json:"transaction"`
	Request_id    string `json:"request_id,omitempty"`
	Currency      string `json:"currency,omitempty"`
	Provider      string `json:"provider,omitempty"`
	Amount        int    `json:"amount,omitempty"`
	Payment_dt    int    `json:"payment_dt,omitempty"`
	Bank          string `json:"bank,omitempty"`
	Delivery_cost int    `json:"delivery_cost,omitempty"`
	Goods_total   int    `json:"goods_total,omitempty"`
	Custom_fee    int    `json:"custom_fee,omitempty"`
}

type Item struct {
	Chrt_id      int    `json:"chrt_id,omitempty"`
	Track_number string `json:"track_number"`
	Price        int    `json:"price,omitempty"`
	Rid          string `json:"rid,omitempty"`
	Name         string `json:"name,omitempty"`
	Sale         int    `json:"sale,omitempty"`
	Size         string `json:"size,omitempty"`
	Total_price  int    `json:"total_price,omitempty"`
	Nm_id        int    `json:"nm_id,omitempty"`
	Brand        string `json:"brand,omitempty"`
	Status       int    `json:"status,omitempty"`
}

type MyModel struct {
	Order_uid          string   `json:"order_uid"`
	Track_number       string   `json:"track_number"`
	Entry              string   `json:"entry,omitempty"`
	Delivery           Delivery `json:"delivery,omitempty"`
	Payment            Payment  `json:"payment,omitempty"`
	Items              []Item   `json:"items,omitempty"`
	Locale             string   `json:"locale,omitempty"`
	Internal_signature string   `json:"internal_signature,omitempty"`
	Customer_id        string   `json:"customer_id,omitempty"`
	Delivery_service   string   `json:"delivery_service,omitempty"`
	Shardkey           string   `json:"shardkey,omitempty"`
	Sm_id              int      `json:"sm_id,omitempty"`
	Date_created       string   `json:"date_created,omitempty"`
	Oof_shard          string   `json:"oof_shard,omitempty"`
}

func printMsg(m *stan.Msg) {
	var myModel MyModel
	json.Unmarshal(m.Data, &myModel)
	fmt.Println(myModel)
}

func main() {
	var (
		clusterID, clientID string
		URL                 string
		qgroup              string
		durable             string
	)

	clusterID = "test-cluster"
	clientID = "myID"

	sc, err := stan.Connect(clusterID, clientID,
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	log.Printf("Connected to clusterID: [%s] clientID: [%s]\n", clusterID, clientID)


	subj := "foo"
	mcb := func(msg *stan.Msg) {
		printMsg(msg)
	}

	sub1, _ := sc.Subscribe(subj, func(m *stan.Msg) {
	    fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	sub1.Unsubscribe()

	sub, err := sc.QueueSubscribe(subj, qgroup, mcb, stan.DurableName(durable))
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
