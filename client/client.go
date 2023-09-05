package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type RemoteExchange struct {
	Bid string `json:"bid"`
}

const ExchangeEndpoint = "http://localhost:8080/cotacao"
const ExchangeRequestTimeout = 300 * time.Millisecond
const FileName = "cotacao.txt"

var ErrGetDolarExchange = fmt.Errorf("falha ao requisitar cotação do dolar")
var ErrServerTimedout = fmt.Errorf("servidor falhou em responder no tempo esperado")
var ErrReadResponseContent = fmt.Errorf("falha ao ler resposta do servidor")
var ErrCreateFile = fmt.Errorf("falha ao criar arquivo de cotação")
var ErrSaveExchange = fmt.Errorf("falha ao salvar cotação")

func main() {
	ctx := context.Background()

	exchange, err := getExchange(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = saveExchange(*exchange)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Cotação salva com sucesso!")
}

func getExchange(ctx context.Context) (*RemoteExchange, error) {
	ctx, cancel := context.WithTimeout(ctx, ExchangeRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ExchangeEndpoint, nil)
	if err != nil {
		return nil, ErrGetDolarExchange
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrGetDolarExchange
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusGatewayTimeout {
		return nil, ErrServerTimedout
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ErrReadResponseContent
	}

	var exchange RemoteExchange
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		return nil, ErrReadResponseContent
	}

	return &exchange, nil
}

func saveExchange(exchange RemoteExchange) error {
	f, err := os.Create(FileName)
	if err != nil {
		return ErrCreateFile
	}

	_, err = fmt.Fprintf(f, "Dólar: %s\n", exchange.Bid)
	if err != nil {
		return ErrSaveExchange
	}

	return nil
}
