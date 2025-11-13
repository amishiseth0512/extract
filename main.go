package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	streamName := "etl"
	subjects := []string{"etl.*"}

	cfg := &nats.StreamConfig{
		Name:      streamName,
		Subjects:  subjects,
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		Replicas:  1,
	}
	_, err = js.AddStream(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stream '%s' created successfully.\n", streamName)

	req, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/organization/audit_logs", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ADMIN_KEY")))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	type auditLogsResponse struct {
		Data []json.RawMessage `json:"data"`
	}
	var data auditLogsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	for _, value := range data.Data {
		msg, err := json.Marshal(value)
		if err != nil {
			log.Fatal(err)
		}
		_, err = js.Publish("etl.open-ai-audit-logs", msg)
		if err != nil {
			log.Fatalf("Error publishing message: %v", err)
		}
		fmt.Println("Message published synchronously.")
	}

}

/*
1. Hit the API
2. Get the data
3. Put the data in NATS
*/
