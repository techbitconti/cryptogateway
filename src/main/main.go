package main

import (
	"fmt"
	"network/http"

	"config"
	"lib/bch"
	"lib/btc"
	"lib/eth"
	"lib/ltc"
	"module/bchScan"
	"module/btcScan"
	"module/dbScan"
	"module/etherScan"
	"module/ltcScan"
	"module/xlmScan"
	//"db/mgodb"
	"db/redis"
)

func main2() {

	report()
}

func main1() {

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
	startLTC()
	startBCH()
	startXLM()

	// GO-3 : dbScan
	startDBScan()

	// Go-4 : start http server
	http.Create(":8082")

}

func startconfig() {
	config.SetPATH("server")
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

func startLTC() {
	ltc.Connect("simnet")
	ltcScan.Start()
}

func startBCH() {
	bch.Connect("simnet", "14.161.40.26")
	bchScan.Start()
}

func startETH() {
	eth.Connect("http://localhost:8545")
	etherScan.Start()
}

func startXLM() {
	xlmScan.Start()
}
