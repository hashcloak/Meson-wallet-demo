package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	wallet "github.com/hashcloak/Meson-wallet-demo"
)

const listenPort = ":18545"

func main() {
	walletCfgFile := flag.String("w", "wallet.toml", "Wallet config file")
	setListen := flag.Bool("l", false, "Listen and serve")
	setChainID := flag.Int64("c", 5, "Chain ID")
	setReceiver := flag.String("a", "", "Address of the receiver")
	setValue := flag.String("v", "10", "Value to transfer")
	setData := flag.String("d", "", "Data to append")
	// for bitcoin transaction
	unspentTxHash := flag.String("utxhash", "", "unspent tx hash (for bitcoin)")
	unspentIndex := flag.Uint("uindex", 0, "unspent index (for bitcoin)")
	left := flag.Int64("left", 0, "left money (should deduct miner fee, for bitcoin)")
	flag.Parse()

	w, err := wallet.New(*walletCfgFile)
	if err != nil {
		panic(err)
	}
	isBitcoin := (len(*unspentTxHash) > 0)
	defer w.Close()
	if *setListen {
		fmt.Println("Listening on port", listenPort[1:])
		http.HandleFunc("/tx", func(resp http.ResponseWriter, req *http.Request) {
			wallet.TransactionHandler(w, resp, req)
		})
		log.Fatal(http.ListenAndServe(listenPort, nil))

	} else {
		fmt.Println("Testing ...")
		value := big.Int{}
		if *setReceiver == "" {
			if isBitcoin {
				panic("please set receiver")
			}
			*setReceiver = w.UiSelectAccount().Address.Hex()
		}
		fmt.Println(".")
		if _, ok := value.SetString(*setValue, 10); !ok {
			fmt.Println("value is invalid")
			os.Exit(0)
		}
		fmt.Println("..")
		request := wallet.TransactionRequest{
			ChainID: *setChainID,
			To:      *setReceiver,
			Value:   value,
			Data:    *setData,
		}
		if isBitcoin {
			request.UnspentTxHash = *unspentTxHash
			request.UnspentIndex = *unspentIndex
			request.Left = *left
			request.Chain = &chaincfg.TestNet3Params
		}
		fmt.Println("...")
		reply, err := wallet.ProcessRequest(w, request)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("....")
		fmt.Println(reply)
	}
}
