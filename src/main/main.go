package main

import (
	"network/http"

	"config"
	"lib/btc"
	"lib/eth"
	"module/dbScan"
)

func main() {

	config.SetPATH("local")

	startETH()
	startBTC()
	dbScan.Start()

	http.Create(":8083")

}

func startBTC() {
	btc.Connect_bitcoind("testnet")
}

func startETH() {
	eth.Connect("http://localhost:8545")
}
