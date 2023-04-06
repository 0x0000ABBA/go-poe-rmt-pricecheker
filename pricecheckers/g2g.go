package pricecheckers

import (
	"strconv"

	"github.com/gocolly/colly"
)

const G2G_URL = "https://www.g2g.com/offer/-PC--Sanctum-Standard---Divine-Orb?service_id=lgc_service_1&brand_id=lgc_game_19398&fa=lgc_19398_tier%3Algc_19398_tier_42692%7Clgc_19398_server%3Algc_19398_server_47225"

func GetG2GPrices(c *colly.Collector) ([]float64, error) {
	prices := []float64{}
	var err error
	
	c.OnHTML("div.offers_bottom-section", func(e *colly.HTMLElement) {
		price, er := strconv.ParseFloat(e.ChildText("span.offers-price-total"), 64)
		if er != nil {
			err = er
		}
		prices = append(prices, price)
	})
	err = c.Visit(G2G_URL)

	if err != nil {
		return nil, err
	}
	return prices, nil
}
