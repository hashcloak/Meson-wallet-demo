package main

import (
	"flag"
	"fmt"

	wallet "github.com/hashcloak/Meson-wallet-demo"
)

func main() {
	walletCfgFile := flag.String("c", "config.toml", "Path to the meson wallet config file")
	rawTransactionBlob := flag.String("r", "", "Raw Transaction blob to send over the network")
	flag.Parse()

	w, err := wallet.New(*walletCfgFile)
	if err != nil {
		panic(err)
	}
	if *rawTransactionBlob == "" {
		if w.Config.Optional == nil {
			panic("must specify a transaction blob in hex or optional configs to sign a txn")
		}
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
