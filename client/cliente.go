package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Panic(err)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}

	if err := createFile(body); err != nil {
		log.Panic(err)
	}
}

func createFile(body []byte) error {

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte("Dólar: " + string(body)))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
