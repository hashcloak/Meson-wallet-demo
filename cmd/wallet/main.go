package main

import (
	"flag"
	"fmt"

	wallet "github.com/hashcloak/Meson-wallet-demo"
)

func main() {
	walletCfgFile := flag.String("w", "wallet.toml", "Wallet config file")
	flag.Parse()

	w, err := wallet.New(*walletCfgFile)
	if err != nil {
		panic(err)
	}
	DemoSetup(w)
	err = DemoSend(w)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done. Shutting down.")
	w.Close()
}
