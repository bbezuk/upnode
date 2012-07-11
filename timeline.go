package upnode

import (
	"github.com/andelf/go-curl/curl"
)

type Timeline struct {
	nameLookupTime    float64
	connectTime       float64
	appConnectTime    float64
	preTransferTime   float64
	startTransferTime float64
	totalTime         float64
	redirectTime      float64
}

func NewTimeline(easy *curl.CURL) *Timeline {
	t := new(Timeline)

	val, _ := easy.Getinfo(curl.INFO_NAMELOOKUP_TIME)
	t.nameLookupTime = val.(float64)
	val, _ = easy.Getinfo(curl.INFO_CONNECT_TIME)
	t.connectTime = val.(float64)
	val, _ = easy.Getinfo(curl.INFO_APPCONNECT_TIME)
	t.appConnectTime = val.(float64)
	val, _ = easy.Getinfo(curl.INFO_PRETRANSFER_TIME)
	t.preTransferTime = val.(float64)
	val, _ = easy.Getinfo(curl.INFO_STARTTRANSFER_TIME)
	t.startTransferTime = val.(float64)
	val, _ = easy.Getinfo(curl.INFO_TOTAL_TIME)
	t.totalTime = val.(float64)
	val, _ = easy.Getinfo(curl.INFO_REDIRECT_TIME)
	t.redirectTime = val.(float64)

	return t
}
