package main

import (
	"network/http"

	"config"
	"lib/btc"
	"lib/eth"
	"module/etherScan"
)

func main() {

	config.SetPATH("server")

	startETH()
	startBTC()

	http.Create(":8082")

}

func startBTC() {
	btc.Connect_bitcoind("simnet")
}

func startETH() {
	eth.Connect("http://localhost:8545")
	etherScan.Start()
}
