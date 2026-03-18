package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Response struct {
	BaseCurr   string  `json:"base_code"`
	TargetCurr string  `json:"target_code"`
	Rate       float64 `json:"conversion_rate"`
	Result     float64 `json:"conversion_result"`
}

func main() {
	baseCurr := "UAH"
	targetCurr := "EUR"
	amount := 100

	exchangeResult, err := Exchange(baseCurr, targetCurr, amount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(exchangeResult)

}

func Exchange(baseCurr string, targetCurr string, amount int) (Response, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")

	reqUrl := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s/%d", apiKey, baseCurr, targetCurr, amount)

	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(reqUrl)
	if err != nil {
		log.Fatalf("GET request failed: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("cannot read body: %v", err)
	}

	var r Response

	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalf("read json error: %v", err)
	}

	return r, nil

}
