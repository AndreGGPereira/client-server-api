package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ExchangeRate struct {
	Bid string `json:"bid"`
}

func main() {

	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	result, err := requestQuote(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	sql, err := openDB()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}
	if err := createTable(sql); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}

	result.ID = uuid.NewString()
	if err := insert(sql, result); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result.Bid)
}

func requestQuote(ctx context.Context) (*USDBRL, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var quote *Quote
	if err := json.Unmarshal(body, &quote); err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	return &quote.USDBRL, nil
}

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "db.sqlite")
}
func createTable(db *sql.DB) error {
	table := `CREATE TABLE IF NOT EXISTS USDBRL (
		id STRING PRIMARY KEY,
		code STRING,
		codein STRING,
		name STRING,
		high STRING,
		low STRING,
		varBid STRING,
		pctChange STRING,
		bid STRING,
		ask STRING,
		timestamp STRING,
		create_date TIMESTAMP
	);`

	_, err := db.Exec(table)
	if err != nil {
		return err
	}
	return nil
}
func insert(db *sql.DB, usdbrl *USDBRL) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	insert := `INSERT INTO USDBRL (id, code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);`

	_, err := db.ExecContext(
		ctx,
		insert,
		usdbrl.ID,
		usdbrl.Code,
		usdbrl.Codein,
		usdbrl.Name,
		usdbrl.High,
		usdbrl.Low,
		usdbrl.VarBid,
		usdbrl.PctChange,
		usdbrl.Bid,
		usdbrl.Ask,
		usdbrl.Timestamp,
	)
	return err
}
