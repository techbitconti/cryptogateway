package btcScan

import (
	"fmt"
	"lib/btc"
	"time"
	//"db/redis"
	"github.com/btcsuite/btcd/wire"
)

func Start() {

	update()
}

func update() {

	go func() {
		for {

			time.Sleep(1 * time.Millisecond)

		}
	}()
}

func getBlock() {

	blockHeight, _ := btc.GetBlockCount()
	blockHash, _ := btc.GetBlockHash(blockHeight)

	block, _ := btc.GetBlock(blockHash)

	fmt.Println(".............block.................", block)

	parse(block)

}

func parse(block *wire.MsgBlock) {

	hashList, _ := block.TxHashes()

	fmt.Println("hashList :  ", hashList)

}
