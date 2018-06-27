package btc

import (
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

func Connect_bitcoind(net string) {

	if Btcd != nil {
		return
	}

	host := ""
	NET = net

	switch net {
	case "mainnet":
		host = "localhost:8332"
		Chaincfg = chaincfg.MainNetParams
	case "testnet":
		host = "localhost:18332"
		Chaincfg = chaincfg.TestNet3Params
	case "simnet":
		host = "localhost:18443"
		Chaincfg = chaincfg.SimNetParams
	}

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         "123",
		Pass:         "123",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Println(err)
	}
	//defer client.Shutdown()

	Btcd = client

}
