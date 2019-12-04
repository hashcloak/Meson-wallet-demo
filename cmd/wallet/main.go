package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/hashcloak/Meson-client/pkg/client"
	"github.com/hashcloak/Meson-wallet-demo/pkg/ethers"
	"github.com/katzenpost/client/config"
)

func main() {
	cfgFile := flag.String("c", "client.toml", "Path to the server config file")
	ticker := flag.String("t", "", "Ticker")
	chainID := flag.Int("chain", 1, "Chain ID for specific ETH-based chain")
	service := flag.String("s", "", "Service Name")
	rawTransactionBlob := flag.String("rt", "", "Raw Transaction blob to send over the network")
	privKey := flag.String("pk", "", "Private key used to sign the txn")
	rpcEndpoint := flag.String("rpc", "http://172.28.1.10:9545", "Ethereum rpc endpoint")
	flag.Parse()

	cfg, err := config.LoadFile(*cfgFile)
	if err != nil {
		panic(err)
	}

	if *rawTransactionBlob == "" {
		if *privKey == "" {
			panic("must specify a transaction blob in hex or a private key to sign a txn")
		}
		rawTransactionBlob, err = produceSignedRawTxn(privKey, rpcEndpoint, chainID)
		if err != nil {
			panic("Raw txn erro: " + err.Error())
		}
	}

	c, err := client.New(cfg, *service)
	if err != nil {
		panic("Client error" + err.Error())
	}

	c.Start()
	reply, err := c.SendRawTransaction(rawTransactionBlob, chainID, ticker)
	if err != nil {
		panic("Meson Request Error" + err.Error())
	}

	fmt.Printf("Reply from the provider: %s\n", reply)
	c.Stop()

}

func produceSignedRawTxn(pk *string, rpcEndpoint *string, chainID *int) (*string, error) {
	ethers, err := ethers.SetURLAndChainID(*rpcEndpoint)
	if err != nil {
		return nil, err
	}

	if ethers.ChainID.Int64() != int64(*chainID) {
		return nil, errors.New("ChainIDs are not the same between rpcEndpoint and chainID flag")
	}
	rawTxn, err := ethers.GenerateSignedRawTxn(*pk)
	return rawTxn, nil
}
