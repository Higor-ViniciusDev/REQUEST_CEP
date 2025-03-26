package main

import (
	"context"
	"encoding/json"
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

	go fazRequest(ctx, ch, "https://viacep.com.br/ws/15771000/json/", 1)
	go fazRequest(ctx, ch, "https://brasilapi.com.br/api/cep/v1/15771000", 2)

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
		endereco.TipoResquest = 1
	} else {
		endereco.Bairro = data["neighborhood"].(string)
		endereco.Cep = data["cep"].(string)
		endereco.Uf = data["state"].(string)
		endereco.Rua = data["street"].(string)
		endereco.Cidade = data["city"].(string)
		endereco.TipoResquest = 2
	}

	return &endereco
}

// cep
// uf or state
// city or localidade
// neighborhood or bairro
// street or logradouro
