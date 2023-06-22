package main

import (
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/mr-tron/base58"
	"github.com/notassigned/tweb/cmd"
)

func main() {
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "test":
			TestNetwork()
		case "genkey":
			genKey()
		default:
			cmd.StartClient(os.Args[1])
		}

	}
	if len(os.Args) == 3 {
		cmd.StartServer()
	}
}

func genKey() {
	k, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	marshalled, _ := crypto.MarshalPrivateKey(k)
	fmt.Println(base58.Encode(marshalled))
}
