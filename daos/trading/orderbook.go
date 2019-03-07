package trading

import (
	"fmt"
	"github.com/linkedin/goavro"
	"github.com/shopspring/decimal"
)

var codecOrderBookZeroMQ *goavro.Codec

type Orderbook struct {
	Type 		string 						`json:"type"`
	Payload		OrderBookZeroMQ				`json:"payload"`
}

type OrderBookZeroMQ struct {
	Rate 					string
	BuyOrderPriceList 		[]PriceAmount
	SellOrderPriceList 		[]PriceAmount
}

func init()  {
	recordSchemaJSONZeroMQ := `{
		"name": "OrderBook",
		"type": "record",
		"fields": [
			{
				"name": "rate",
				"type": "string"
			},
			{
				"name": "buyOrders",
				"type": {
					"type": "array",
					"items": {
						"name": "Order",
						"type": "record",
						"fields": [
							{
								"name": "price",
								"type": "string"
							},
							{
								"name": "amount",
								"type": "string"
							}
						]
					}
				}
			},
			{
				"name": "sellOrders",
				"type": {
					"type": "array",
					"items": {
						"name": "Order",
						"type": "record",
						"fields": [
							{
								"name": "price",
								"type": "string"
							},
							{
								"name": "amount",
								"type": "string"
							}
						]
					}
				}
			}
		]
	}`
	codecOrderBookZeroMQ, _ = goavro.NewCodec(recordSchemaJSONZeroMQ)
}

func OrderBookFromJSONZeroMQ(msg []byte) OrderBookZeroMQ {
	native, _, err := codecOrderBookZeroMQ.NativeFromTextual(msg)
	if err != nil {
		fmt.Println(err)
	}
	OrderBookMap := native.(map[string]interface{})
	rate := OrderBookMap["rate"].(string)
	buyOrders := OrderBookMap["buyOrders"].([]interface{})
	sellOrders := OrderBookMap["sellOrders"].([]interface{})

	var orderBook OrderBookZeroMQ
	buyOrderPriceList := make([]PriceAmount, 0, 1)
	sellOrderPriceList := make([]PriceAmount, 0, 1)

	for _, buyOrder := range buyOrders {
		buyOrderMap := buyOrder.(map[string]interface{})

		var priceAmount PriceAmount
		priceAmount.Price, _ = decimal.NewFromString(buyOrderMap["price"].(string))
		priceAmount.Amount, _ = decimal.NewFromString(buyOrderMap["amount"].(string))

		buyOrderPriceList = append(buyOrderPriceList, priceAmount)
	}

	for _, sellOrder := range sellOrders {
		buyOrderMap := sellOrder.(map[string]interface{})

		var priceAmount PriceAmount
		priceAmount.Price, _ = decimal.NewFromString(buyOrderMap["price"].(string))
		priceAmount.Amount, _ = decimal.NewFromString(buyOrderMap["amount"].(string))

		sellOrderPriceList = append(sellOrderPriceList, priceAmount)
	}

	orderBook.Rate = rate
	orderBook.BuyOrderPriceList = buyOrderPriceList
	orderBook.SellOrderPriceList = sellOrderPriceList

	return orderBook
}


