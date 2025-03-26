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
	TipoResquest string
	Bairro       string
}

func fazRequest(ctx context.Context, ch chan<- Endereco, url string) {
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

	convertJsonRequest(body)
}

func main() {
	//Criação do time de 200 milisec
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan Endereco)

	fazRequest(ctx, ch, "https://viacep.com.br/ws/15771000/json/")
}

func convertJsonRequest(r *http.Response) {
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

	fmt.Println(data["cep"].(string))
}

// cep
// uf or state
// city or localidade
// neighborhood or bairro
// street or logradouro
