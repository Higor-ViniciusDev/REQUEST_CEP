package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Endereco struct {
	Cep          string
	Uf           string
	Cidade       string
	TipoResquest int
	Bairro       string
	Rua          string
}

func fazRequest(ctx context.Context, ch chan<- *Endereco, url string, tipoRequest int) {
	defer close(ch)

	//Por algum motivo o request do viacep está demorando muito para responder, por isso adicionado um time com 350 milisec para validar
	if tipoRequest == 1 {
		time.Sleep(800 * time.Millisecond)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		panic(err)
	}

	cliente := &http.Client{}
	body, err := cliente.Do(req)

	if err != nil {
		panic(err)
	}

	defer body.Body.Close()

	end := convertJsonRequest(body, tipoRequest)

	ch <- end
}

func main() {
	//Criação do time de 200 milisec
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan *Endereco)
	newCh := make(chan *Endereco)

	go fazRequest(ctx, ch, "https://viacep.com.br/ws/15771000/json/", 1)
	go fazRequest(ctx, newCh, "https://brasilapi.com.br/api/cep/v1/15771034", 2)

	select {
	case end := <-ch:
		fmt.Println("API ViaCep Respondeu primeiro")
		fmt.Printf("CEP: %s\nUF: %s\nCidade: %s\nBairro: %s\nRua: %s\n", end.Cep, end.Uf, end.Cidade, end.Bairro, end.Rua)
	case end := <-newCh:
		fmt.Println("API brasilapi Respondeu primeiro")
		fmt.Printf("CEP: %s\nUF: %s\nCidade: %s\nBairro: %s\nRua: %s\n", end.Cep, end.Uf, end.Cidade, end.Bairro, end.Rua)
	case <-ctx.Done():
		fmt.Println("Timeout")
	}

}

func convertJsonRequest(r *http.Response, tipoRequest int) *Endereco {
	// var textoJson byte

	read, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Criando um mapa genérico para armazenar os dados do JSON
	var data map[string]any
	if err := json.Unmarshal(read, &data); err != nil {
		panic("ERROR")
	}

	endereco := Endereco{}
	if tipoRequest == 1 {
		endereco.Bairro = data["bairro"].(string)
		endereco.Cep = data["cep"].(string)
		endereco.Uf = data["uf"].(string)
		endereco.Rua = data["logradouro"].(string)
		endereco.Cidade = data["localidade"].(string)
		endereco.TipoResquest = tipoRequest
	} else {
		endereco.Bairro = data["neighborhood"].(string)
		endereco.Cep = data["cep"].(string)
		endereco.Uf = data["state"].(string)
		endereco.Rua = data["street"].(string)
		endereco.Cidade = data["city"].(string)
		endereco.TipoResquest = tipoRequest
	}

	return &endereco
}

// cep
// uf or state
// city or localidade
// neighborhood or bairro
// street or logradouro
