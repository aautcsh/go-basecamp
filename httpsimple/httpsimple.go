package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	bearer := "Bearer 00D6A000000tiQR!AQEAQFLIzcR18VwOaOvZHgnEb9tBewf9b5tnSwmMsdpdrn7Xp6oWmAbyUVsvmhLV85FgUy6M4ENhxSkKCFBKqsTXar3wMg9S"

	url := fmt.Sprintf("https://dvelop-documents-dev-ed.lightning.force.com/services/data/v41.0")

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	if err != nil {
		log.Fatal("err::request: ", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("err:response: ", err)
		return
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("err::body: ", err)
		return
	}

	fmt.Println(string(data))
}
