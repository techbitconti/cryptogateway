package eth

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_SendPriv(t *testing.T) {

	Connect("http://localhost:8545")

	amETH, _ := strconv.ParseFloat("1000", 64)
	valueWei := ToWei(amETH, "ether") // amETH * math.Pow10(18)

	s := strconv.FormatFloat(valueWei, 'f', -1, 64)
	bi := ToBig256(s)
	fmt.Println(s, bi)

	SendTransactionRaw("47de15108b35169c4aff4826d5c413fe117e361a900325f6d3df1f0e04cbd706", "0x56bfed926439bb905fca945f47d4e4ffff9a77b8", bi, []byte{})

}
