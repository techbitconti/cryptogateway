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
	startLTC()
	startBCH()
	startXLM()

	dbScan.LoadReport_BTC()
	dbScan.LoadReport_LTC()
	dbScan.LoadReport_BCH()
	dbScan.LoadReport_ETH()
	dbScan.LoadReport_XLM()

}
