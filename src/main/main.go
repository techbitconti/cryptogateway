package main

import (
	"network/http"

	"lib/btc"
	"lib/eth"

	"module/dbScan"
)

func main() {

	startETH()
	startBTC()
	dbScan.Start()

	http.Create(":8082")

}

func startBTC() {
	btc.Connect_bitcoind("testnet")
}

func startETH() {
	eth.Connect("http://localhost:8545")
}
