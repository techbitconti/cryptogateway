package network

import (
	//"fmt"
	"time"
)

func Threading() {

	//now_rank := time.Now().Unix()

	go func() {
		for {

			//fmt.Println(time.Now().Second())
			time.Sleep(1 * time.Second)
		}
	}()
}
