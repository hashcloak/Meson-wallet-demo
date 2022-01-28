package main

import (
	"flag"
	"fmt"
	"math/big"

	wallet "github.com/hashcloak/Meson-wallet-demo"
)

func main() {
	var passphrase string
	walletCfgFile := flag.String("w", "wallet.toml", "Wallet config file")
	flag.Parse()

	w, err := wallet.New(*walletCfgFile)
	if err != nil {
		panic(err)
	}
	if len(w.Accounts()) == 0 {
		fmt.Print("Init passphrase: ")
		_, err := fmt.Scanf("%s", &passphrase)
		if err != nil {
			panic(err)
		}
		if len(passphrase) < 4 {
			panic("passphrase has to be at least 4 characters")
		}
		ac, err := w.NewAccount(passphrase)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Account created with address %v\n", ac.Address)
	}
	ac := w.Accounts()[0]
	addr := ac.Address
	fmt.Printf("Using account %v\n", addr)
	tx, err := wallet.GenerateTx(addr, addr, w.ChainID(), w.Endpoint())
	if err != nil {
		panic(err)
	}
	fmt.Print("Enter passphrase: ")
	_, err = fmt.Scanf("%s", &passphrase)
	if err != nil {
		panic(err)
	}
	err = w.Unlock(ac, passphrase)
	if err != nil {
		panic(err)
	}
	signedTx, err := w.SignTx(ac, tx, big.NewInt(w.ChainID()))
	if err != nil {
		panic(err)
	}
	reply, err := w.Send(signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
	fmt.Println("Done. Shutting down.")
	w.Close()
}
