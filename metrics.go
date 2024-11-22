package main

import "math"

// RSI calculation function
// 100 - (100 / (1 + RS))
// RS = Average Gain / Average Loss
func calculateRSI(prices []float64, period int) float64 {
	if len(prices) < period {
		return 50 // Neutral
	}

	var gains, losses float64
	start := len(prices) - period
	for i := start + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	rsi := 100.0 - (100.0 / (1 + rs))
	return rsi
}

const shortMomentumPeriod = 6 // Shorter, more sensitive period
const longMomentumPeriod = 30 // Longer, smoother period

func calculateMomentum(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0 // Not enough data
	}
	return prices[len(prices)-1] - prices[len(prices)-period]
}

// simpleMovingAverage calculates the average of the last windowSize elements.
func simpleMovingAverage(data []float64, windowSize int) float64 {
	start := len(data) - windowSize
	if start < 0 {
		start = 0
		windowSize = len(data)
	}

	sum := 0.0
	for _, value := range data[start:] {
		sum += value
	}
	return sum / float64(windowSize)
}

// float64ToInt rounds the float64 to the nearest integer.
func float64ToInt(f float64) int {
	return int(math.Round(f))
}
