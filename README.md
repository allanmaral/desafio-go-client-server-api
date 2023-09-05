# Desafio Go - Cotação do Dólar

Esse repositório, contem a solução para o desafio proposto, que envolve a criação de dois sistemas em Go: `client.go` e `server.go`. O objetivo é criar um cliente que solicita a cotação do dólar a partir de um servidor, que consome uma API de câmbio e retorna o resultado no formato JSON. Além disso, o servidor deve registrar as cotações em um banco de dados SQLite e os dois sistemas devem usar contextos para gerenciar timeouts e retornar erros em caso de execução insuficiente.

## Pré-requisitos

Antes de executar os sistemas, certifique-se de ter o Go instalado em sua máquina. Você pode baixá-lo em [https://go.dev/dl/](https://go.dev/dl/).

## Como executar

Siga as instruções abaixo para executar o cliente e o servidor:

### Server (server.go)

1. Navegue até o diretório do servidor:

   ```bash
   cd server
   ```

2. Instale as depdendências:
  
   ```bash
   go mod tidy
   ```

3. Execute o servidor:

   ```bash
   go run server.go
   ```

O servidor estará em execução na porta 8080 e responderá às solicitações na rota `http://localhost:8080/cotacao`.

### Client (client.go)

1. Navegue até o diretório do cliente:

   ```bash
   cd client
   ```

2. Execute o cliente:

   ```bash
   go run client.go
   ```

O cliente enviará uma solicitação ao servidor e salvará a cotação do dólar atual em um arquivo chamado `cotacao.txt`.

## Timeout

O cliente e o servidor usam contextos para gerenciar timeouts. O cliente tem um timeout máximo de 300ms para receber a resposta do servidor. O servidor tem um timeout máximo de 200ms para chamar a API de cotação do dólar e 10ms para persistir os dados no banco de dados SQLite. Se algum desses timeouts for excedido, o sistema retornará um erro nos logs.

## Banco de Dados

O servidor utiliza um banco de dados SQLite para registrar as cotações recebidas. O banco de dados é criado automaticamente e os registros são armazenados na tabela "exchanges". Você pode consultar o banco de dados para obter o histórico de cotações.

## Estrutura do Projeto

- `client/client.go`: Implementação do cliente que faz a solicitação ao servidor e salva a cotação em um arquivo.
- `server/server.go`: Implementação do servidor que consome a API de câmbio, registra as cotações no banco de dados e responde às solicitações do cliente.

