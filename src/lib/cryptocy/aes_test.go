package cryptocy

import (
	//"encoding/json"
	"log"
	"testing"
)

/*
func Test_1(t *testing.T) {
	CIPHER_KEY := []byte("01234567890123450123456789012345")

	js := map[string]interface{}{
		"amount": 100,
		"to":     "0x5e531bb813994f27f65d2d5f7dc7c51dbe5406f7",
	}
	bb, _ := json.MarshalIndent(js, " ", "")
	msg := string(bb)

	if encrypted, err := Encrypt(CIPHER_KEY, msg); err != nil {
		log.Println(err)
	} else {
		log.Printf("CIPHER KEY: %s\n", string(CIPHER_KEY))
		log.Printf("ENCRYPTED: %s\n", encrypted)

		if decrypted, err := Decrypt(CIPHER_KEY, encrypted); err != nil {
			log.Println(err)
		} else {
			log.Printf("DECRYPTED: %s\n", decrypted)
		}
	}
}
*/

func Test_2(t *testing.T) {

	CIPHER_KEY := GenKey([]byte("47de15108b35169c4aff4826d5c413fe117e361a900325f6d3df1f0e04cbd706"))

	msg := "SCS4IITATPILLHQDGMWWZPJROJLWTZVPYTNEOG57ZPJXGBNXGUTRWMG3"

	if encrypted, err := Encrypt(CIPHER_KEY, msg); err != nil {
		log.Println(err)
	} else {
		log.Printf("CIPHER KEY: %s\n", string(CIPHER_KEY))
		log.Printf("ENCRYPTED: %s\n", encrypted)

		if decrypted, err := Decrypt(CIPHER_KEY, encrypted); err != nil {
			log.Println(err)
		} else {
			log.Printf("DECRYPTED: %s\n", decrypted)
		}
	}
}
