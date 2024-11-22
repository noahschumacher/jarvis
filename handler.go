package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func chartHandler(p *priceModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		type pdata struct {
			Price         float64 `json:"price"`
			RSI           float64 `json:"rsi"`
			ShortMomentum float64 `json:"shortMomentum"`
			LongMomentum  float64 `json:"longMomentum"`
			BuySignal     int     `json:"buySignal"`
			SellSignal    int     `json:"sellSignal"`
			Action        string  `json:"action"`
		}

		for d := range p.data {
			pd := pdata{
				Price:         d.price,
				RSI:           d.rsi,
				ShortMomentum: d.shortMoment,
				LongMomentum:  d.longMoment,
				BuySignal:     d.buySignal,
				SellSignal:    d.sellSignal,
				Action:        d.action,
			}

			jsonData, _ := json.Marshal(pd)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			w.(http.Flusher).Flush()
		}
	}
}
