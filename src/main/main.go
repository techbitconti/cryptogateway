package main

import (
	"network/http"

	"config"
	"lib/btc"
	"lib/eth"
)

func main() {

	config.SetPATH("server")

	startETH()
	//startBTC()

	http.Create(":8082")

}

func startBTC() {
	btc.Connect_bitcoind("testnet")
}

func startETH() {
	eth.Connect("http://localhost:8545")
}
