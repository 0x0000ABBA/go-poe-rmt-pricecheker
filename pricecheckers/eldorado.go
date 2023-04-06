package pricecheckers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
)

const ELDORADO_URL = "https://www.eldorado.gg/api/predefinedOffers/augmentedItem/4d6bc251-6643-4bea-b9b1-542e21f106a4/?pageIndex=1&pageSize=20"

type EldoradoResponse struct {
	Results []EldoradoItem
}

type EldoradoItem struct {
	Offer EldoradoOffer `json:"offer"`
}

type EldoradoOffer struct {
	Quantity int           `json:"quantity"`
	Price    EldoradoPrice `json:"pricePerUnit"`
}

type EldoradoPrice struct {
	Amount float64 `json:"amount"`
}

func GetEldoradoPrices() ([]float64, error) {

	var items EldoradoResponse
	prices := []float64{}
	resp, err := http.Get(ELDORADO_URL)
	if err != nil {
		return nil, err 
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err 
	}
	json.Unmarshal(body, &items)
	for _, v := range items.Results {
		if v.Offer.Quantity > MIN_AMOUNT {
			prices = append(prices, v.Offer.Price.Amount)
		}
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i] < prices[j]
	})
	return prices, nil
}