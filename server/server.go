package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const ExchangeEndpoint = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
const ExchangeRequestTimeout = 200 * time.Millisecond
const ExchangeInsertTimeout = 10 * time.Millisecond

type RemoteExchangeResponse struct {
	Content RemoteExchange `json:"USDBRL"`
}

type RemoteExchange struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
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

type Exchange struct {
	ID         int
	Code       string
	CodeIn     string
	Name       string
	High       string
	Low        string
	VarBid     string
	PctChange  string
	Bid        string
	Ask        string
	Timestamp  string
	CreateDate string
}

type ExchangeHandler struct {
	DB *sql.DB
}

func main() {
	db, err := openConn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.Handle("/cotacao", &ExchangeHandler{DB: db})

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func (h *ExchangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	re, err := getDollarExchangeRate(ctx)
	if err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	exchange := mapExchange(re)

	err = insertExchange(ctx, h.DB, exchange)
	if err != nil {
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(re)
}

func getDollarExchangeRate(ctx context.Context) (*RemoteExchange, error) {
	ctx, cancel := context.WithTimeout(ctx, ExchangeRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ExchangeEndpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var re RemoteExchangeResponse
	err = json.Unmarshal(body, &re)
	if err != nil {
		return nil, err
	}

	return &re.Content, nil
}

func mapExchange(re *RemoteExchange) *Exchange {
	return &Exchange{
		Code:       re.Code,
		CodeIn:     re.CodeIn,
		Name:       re.Name,
		High:       re.High,
		Low:        re.Low,
		VarBid:     re.VarBid,
		PctChange:  re.PctChange,
		Bid:        re.Bid,
		Ask:        re.Ask,
		Timestamp:  re.Timestamp,
		CreateDate: re.CreateDate,
	}
}

func openConn() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "exchange.db")
	if err != nil {
		return nil, err
	}

	err = ensureTableExists(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func ensureTableExists(db *sql.DB) error {
	stmt := `
        CREATE TABLE IF NOT EXISTS exchanges(
            id INTEGER PRIMARY KEY,
            code TEXT,
            code_in TEXT,
            name TEXT,
            high TEXT,
            low TEXT,
            var_bid TEXT,
            pct_change TEXT,
            bid TEXT,
            ask TEXT,
            timestamp TEXT,
            create_date TEXT
        );
    `
	_, err := db.Exec(stmt)
	if err != nil {
		return errors.New("failed to ensure the exchanges table exists")
	}

	return nil
}

func insertExchange(ctx context.Context, db *sql.DB, e *Exchange) error {
	ctx, cancel := context.WithTimeout(ctx, ExchangeInsertTimeout)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, "INSERT INTO exchanges(code, code_in, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) VALUES(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(e.Code, e.CodeIn, e.Name, e.High, e.Low, e.VarBid, e.PctChange, e.Bid, e.Ask, e.Timestamp, e.CreateDate)
	if err != nil {
		return err
	}

	return nil
}
