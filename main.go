package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	req, err := http.NewRequest(http.MethodGet, "https://api.box.com/2.0/events", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authorization", "Bearer <ACCESS_TOKEN>")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	type data struct{
		Entries []json.RawMessage `json:"entries"`
	}
	var entries data 
	err = json.Unmarshal(body, &entries)
	if err != nil {
		log.Fatal(err)
	}

	for _,value := range entries.Entries {
		
	}

}

/*
1. Hit the API
2. Get the data
3. Put the data in NATS
*/
