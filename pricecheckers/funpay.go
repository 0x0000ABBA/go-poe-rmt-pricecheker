package pricecheckers

import (
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

const FUNPAY_URL = "https://funpay.com/en/chips/173/"
const LEAGUE = "(PC) The Forbidden Sanctum"
const MIN_AMOUNT = 20
const TOLERANCE = 0.005

func GetFunpayPrices(c *colly.Collector) ([]float64, error) {

	prices := []float64{}
	var er error

	c.OnHTML("a.tc-item", func(e *colly.HTMLElement) {
		server := e.ChildText("div.tc-server")
		quantity, err := strconv.Atoi(strings.ReplaceAll(e.ChildText("div.tc-amount"), " ", ""))
		if er != nil {
			er = err
		}
		if server == LEAGUE && quantity >= MIN_AMOUNT {
			price, err := strconv.ParseFloat(strings.Trim(e.ChildText("div.tc-price"), " €"), 64)
			if err != nil {
				er = err
			}
			prices = append(prices, price)
		}
	})

	if er != nil {
		return nil, er
	}

	er = c.Visit(FUNPAY_URL)

	if er != nil {
		return nil, er
	}

	return prices, nil
}
