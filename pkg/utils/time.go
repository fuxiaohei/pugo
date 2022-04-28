package utils

import (
	"time"
)

const (
	kDelay          = true
	kNoDelay        = false
	defalutInterval = 60
)

func Ticker(interval time.Duration, fn func()) {
	fn()
	tickerDone(interval, fn)
}

func TickerDelayed(interval time.Duration, function func()) {
	tickerDone(interval, function)
}

func AsyncTicker(interval time.Duration, function func()) {
	go function()
	tickerDone(interval, function)
}

func tickerDone(interval time.Duration, function func()) {

	if interval <= 0 {
		interval = defalutInterval
	}

	go func() {
		eventsTick := time.NewTicker(interval)
		defer eventsTick.Stop()

		for {
			select {
			case <-eventsTick.C:
				function()
			}
		}
	}()
}
