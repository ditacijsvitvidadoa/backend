package app

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/cash"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"io/ioutil"
	"net/http"
	"os"
)

func (a *App) GetAllCities(w http.ResponseWriter, r *http.Request) {
	cities, err := cash.GetCitiesFromRedis(a.cash.Conn)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, cities)
}

func (a *App) GetPostalsFromCity(w http.ResponseWriter, r *http.Request) {
	cityRef := r.PathValue("city_ref")

	postals, err := GetPostals(cityRef)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, postals)
}

func GetPostals(cityRef string) ([]entities.Postal, error) {
	url := "https://api.novaposhta.ua/v2.0/json/"

	apiKey, ok := os.LookupEnv("NP_API_KEY")
	if !ok {
		return nil, nil
	}

	requestBody := map[string]interface{}{
		"apiKey":       apiKey,
		"modelName":    "AddressGeneral",
		"calledMethod": "getWarehouses",
		"methodProperties": map[string]string{
			"CityRef": cityRef,
		},
	}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error creating JSON request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Игнорирование проверки сертификата
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request to Nova Poshta API: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error response from Nova Poshta API: %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Error parsing JSON response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Error in response format")
	}

	var postals []entities.Postal
	for _, item := range data {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		postal := entities.Postal{
			Description: itemMap["Description"].(string),
			Ref:         itemMap["Ref"].(string),
		}
		postals = append(postals, postal)
	}

	return postals, nil
}
