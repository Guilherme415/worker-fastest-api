package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func getCepDetails(cep string, url string, channel chan<- CepResponse) (CepResponse, error) {
	response := CepResponse{}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	channel <- response

	return response, err
}

func main() {
	viaCepChannel := make(chan CepResponse)
	apicepChannel := make(chan CepResponse)

	cep := "14403720"

	viacepUrl := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	apicepUrl := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)

	go getCepDetails(cep, viacepUrl, viaCepChannel)
	go getCepDetails(cep, apicepUrl, apicepChannel)

	select {
	case msg := <-viaCepChannel:
		fmt.Printf("Viacep channel response has arrived before, msg: %s", msg)

	case msg := <-apicepChannel:
		fmt.Printf("Apicep channel response has arrived before, msg: %s", msg)

	case <-time.Tick(time.Second * 1):
		fmt.Print("Timeout has occured")
	}

}
