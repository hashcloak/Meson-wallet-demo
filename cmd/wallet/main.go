package main

import (
	"flag"
	"fmt"

	wallet "github.com/hashcloak/Meson-wallet-demo"
)

func main() {
	walletCfgFile := flag.String("w", "wallet.toml", "Wallet config file")
	rawTransactionBlob := flag.String("r", "", "Raw transaction blob to send over the network")
	flag.Parse()

	w, err := wallet.New(*walletCfgFile)
	if err != nil {
		panic(err)
	}
	if *rawTransactionBlob == "" {
		rawTransactionBlob, err = wallet.GenerateTransaction(w)
		if err != nil {
			panic(err)
		}
	}
	reply, err := w.Send(*rawTransactionBlob)
	if err != nil {
		panic(err)
	}
	fmt.Printf("reply: %s\n", reply)
	fmt.Println("Done. Shutting down.")
	w.Close()
}
