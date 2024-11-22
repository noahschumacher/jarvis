package main

import "fmt"

const (
	rsiPeriod = 14

	actionBuy  = "BUY"
	actionSell = "SELL"
)

type priceModel struct {
	wallet  wallet
	fetcher *priceFetcher
	data    chan modelData
}

type wallet struct {
	usd  float64 // USD balance
	coin float64 // Coin balance (e.g. BTC)
}

type modelData struct {
	price       float64
	rsi         float64
	longMoment  float64
	shortMoment float64
	buySignal   int
	sellSignal  int
	action      string
}

func newPriceModel(name, addr string) *priceModel {
	return &priceModel{
		wallet:  wallet{usd: 1000, coin: 0},
		fetcher: newPriceFetcher(name, addr),
		data:    make(chan modelData),
	}
}

func (pm *priceModel) run() {
	go func() {
		if err := pm.fetcher.fetchPrice(); err != nil {
			fmt.Println("error fetching price:", err)
		}
	}()

	prices := make([]float64, 0, rsiPeriod)

	var (
		RSIs          []float64
		shortMomentum []float64
		longMomentum  []float64

		shortAboveLong = false
	)

	for p := range pm.fetcher.priceChan {
		prices = append(prices, p.priceUSD)

		rsi := calculateRSI(prices, rsiPeriod)
		shortMoment := calculateMomentum(prices, shortMomentumPeriod)
		longMoment := calculateMomentum(prices, longMomentumPeriod)

		RSIs = append(RSIs, rsi)
		shortMomentum = append(shortMomentum, shortMoment)
		longMomentum = append(longMomentum, longMoment)

		smoothedRSI := simpleMovingAverage(RSIs, 5)
		smoothedShort := simpleMovingAverage(shortMomentum, 5)
		smoothedLong := simpleMovingAverage(longMomentum, 5)

		buy, sell, action := pm.buyer(p.priceUSD, smoothedRSI, smoothedShort, smoothedLong, shortAboveLong)

		md := modelData{
			price:       p.priceUSD,
			rsi:         smoothedRSI,
			shortMoment: smoothedShort,
			longMoment:  smoothedLong,
			buySignal:   buy,
			sellSignal:  sell,
			action:      action,
		}

		pm.data <- md

		fmt.Printf("%s, RSI: %.1f, Short: %.4f, Long: %.4f, BuyS: %d, SellS: %d\n",
			p, smoothedRSI, smoothedShort, smoothedLong, buy, sell,
		)

		shortAboveLong = shortMoment > longMoment
	}
}

// buyer uses the RSI indictar and the short, long momentum to determine if
// it should buy or sell the asset. The logic is as follows:
func (pm *priceModel) buyer(price, rsi, short, long float64, shortAboveLong bool) (int, int, string) {
	var (
		momentumBuy  float64
		momentumSell float64
	)
	if shortAboveLong && short < long {
		momentumSell = 1 // short moment moving below long indicates sell
	} else if !shortAboveLong && short > long {
		momentumBuy = 1 // short moment moving above long indicates buy
	}

	buySignal := float64ToInt((100 - rsi) + 30*momentumBuy)
	sellSignal := float64ToInt(rsi + 30*momentumSell)

	if pm.wallet.usd > 0 && buySignal > 90 {
		pm.wallet.coin = pm.wallet.usd / price
		pm.wallet.usd = 0
		fmt.Println("BUY 🤑", price, "RSI:", rsi, "Coin:", pm.wallet.coin)
		return buySignal, sellSignal, actionBuy
	}

	if pm.wallet.coin > 0 && sellSignal > 90 {
		pm.wallet.usd = pm.wallet.coin * price
		pm.wallet.coin = 0
		fmt.Println("SELL 💸", price, "RSI:", rsi, "USD:", pm.wallet.usd)
		return buySignal, sellSignal, actionSell
	}

	return buySignal, sellSignal, ""
}
