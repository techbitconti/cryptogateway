package main

import (
	"module/dbScan"
)

func report() {

	// GO-0 : load config
	startconfig()

	// GO-1 : start mongod-redis
	startDB()

	// G0-2: start module
	startETH()
	startBTC()

	dbScan.LoadReport_BTC()
	dbScan.LoadReport_ETH()
}
