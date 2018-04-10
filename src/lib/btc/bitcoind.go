package btc

import (
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

func Connect_bitcoind(net string) {

	if Btcd != nil {
		return
	}

	host := ""
	switch net {
	case "mainnet":
		host = "localhost:8332"
	case "testnet":
		host = "localhost:18332"
	case "simnet":
		host = "localhost:18554"
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
