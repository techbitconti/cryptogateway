package main

import (
	"fmt"
	"network/http"

	"config"
	"lib/btc"
	"lib/eth"
	"module/btcScan"
	"module/dbScan"
	"module/etherScan"
	//"db/mgodb"
	"db/redis"
)

func main1() {
	// GO-0 : load config
	startconfig()

	// GO-1 : start mongod-redis
	startDB()

	// G0-2: start module
	startETH()
	startBTC()

	//
	withdraw()
}

func main() {

	// GO-0 : load config
	startconfig()

	// GO-1 : start mongod-redis
	startDB()

	// G0-2: start module
	startETH()
	startBTC()
	startDBScan()

	// Go-3 : start http server
	http.Create(":8082")

}

func startconfig() {
	config.SetPATH("local")
}

func startDB() {

	//	mgodb.Connect()
	//	fmt.Println("MongoDB ................")

	redis.Connect()
	fmt.Println("Redis ................")
}

func startDBScan() {
	dbScan.Start()
}

func startBTC() {
	btc.Connect_bitcoind("simnet")
	btcScan.Start()
}

func startETH() {
	eth.Connect("http://localhost:8545")
	etherScan.Start()
}
