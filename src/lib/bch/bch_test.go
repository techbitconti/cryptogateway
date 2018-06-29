package bch

import (
	"fmt"
	"testing"
)

func Test_bitcash(t *testing.T) {

	Connect("simnet")

	addr, err := GetNewAddress("add")
	fmt.Println(addr, err)
}
