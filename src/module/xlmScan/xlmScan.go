package xlmScan

import (
	"time"
)

func Start() {

	update()
}

func update() {

	go func() {
		for {

			getBlock()

			time.Sleep(1 * time.Second)

		}
	}()
}

func getBlock() {

	parse()
}

func parse() {

}
