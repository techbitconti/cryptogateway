package bch

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func Test_bitcash(t *testing.T) {

	Connect("simnet", "14.161.40.26")

	addr, priv := genAddressBCH()
	fmt.Println(addr, priv)
}

func genAddressBCH() (string, string) {

	utc := time.Now().Unix()
	decode := strconv.FormatInt(utc, 10)
	encode := base64.StdEncoding.EncodeToString([]byte(decode))

	address, err1 := GetNewAddress(encode)
	if err1 != nil {
		return "", ""
	}
	privKey, err2 := DumpPrivKey(address.String())
	if err2 != nil {
		return "", ""
	}

	return address.String(), privKey.String()
}
