package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApiCep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func getInputCep() string {
	inputCep := ""
	for _, cep := range os.Args[1:] {
		inputCep = cep
	}
	return inputCep
}

func getFromViaCep(ch chan ViaCEP, cep string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://viacep.com.br/ws/"+cep+"/json/", nil)

	//teste, erro := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintf(os.Stderr, "Timeout: o tempo limite foi excedido")
		} else {
			fmt.Fprintf(os.Stderr, "Erro ao fazer requisição, %v\n", err)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao processar requisição, %v\n", err)
	}

	var data ViaCEP
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}

	ch <- data
}

func getFromBrasilApi(ch chan BrasilApiCep, cep string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintf(os.Stderr, "Timeout: o tempo limite foi excedido")
		} else {

			fmt.Fprintf(os.Stderr, "Erro ao processar requisição, %v\n", err)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao processar requisição, %v\n", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível ler a resposta.")
	}

	var data BrasilApiCep
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Não foi possível processar a resposta.")
	}
	//	time.Sleep(time.Second * 10)

	ch <- data
}

func main() {

	ch1 := make(chan ViaCEP)
	ch2 := make(chan BrasilApiCep)
	go getFromViaCep(ch1, getInputCep())
	go getFromBrasilApi(ch2, getInputCep())
	select {
	case msg := <-ch1:
		fmt.Printf("Endereço: %s %s %s %s %s\nRecebido de ViaCep \n", msg.Logradouro, msg.Bairro, msg.Cep, msg.Localidade, msg.Uf)
	case msg := <-ch2:
		fmt.Printf("Endereço: %s %s %s %s %s\nRecebido de BrasilApi \n", msg.Street, msg.Neighborhood, msg.Cep, msg.City, msg.State)
	}

}
