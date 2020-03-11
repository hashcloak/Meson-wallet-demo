package main

import (
	"flag"
	"fmt"

	client "github.com/hashcloak/Meson-client"
	"github.com/hashcloak/Meson-wallet-demo/pkg/ethers"
)

func main() {
	cfgFile := flag.String("c", "client.toml", "Path to the server config file")
	ticker := flag.String("t", "", "Ticker")
	//service := flag.String("s", "", "Service Name")
	rawTransactionBlob := flag.String("rt", "", "Raw Transaction blob to send over the network")
	privKey := flag.String("pk", "", "Private key used to sign the txn")
	rpcEndpoint := flag.String("rpc", "https://goerli.hashcloak.com", "Ethereum rpc endpoint")
	flag.Parse()

	client, err := client.New(*cfgFile, *ticker)
	if err != nil {
		panic("ERROR In creating new client: " + err.Error())
	}

	if *rawTransactionBlob == "" {
		if *privKey == "" {
			panic("must specify a transaction blob in hex or a private key to sign a txn")
		}
		rawTransactionBlob, err = produceSignedRawTxn(privKey, rpcEndpoint)
		if err != nil {
			panic("Raw txn error: " + err.Error())
		}
	}

	client.Start()
	reply, err := client.SendRawTransaction(rawTransactionBlob, ticker)
	if err != nil {
		panic("ERROR Send raw transaction: " + err.Error())
	}

	fmt.Printf("reply: %s\n", reply)
	fmt.Println("Done. Shutting down.")
	client.Shutdown()
}

func produceSignedRawTxn(pk *string, rpcEndpoint *string) (*string, error) {
	ethers, err := ethers.SetURLAndChainID(*rpcEndpoint)
	if err != nil {
		return nil, err
	}

	rawTxn, err := ethers.GenerateSignedRawTxn(*pk)
	return rawTxn, nil
}
