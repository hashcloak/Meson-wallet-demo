package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

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
	flag.Parse()

	w, err := wallet.New(*walletCfgFile)
	if err != nil {
		panic(err)
	}
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
