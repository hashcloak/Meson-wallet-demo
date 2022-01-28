package main

import (
	"fmt"
	"math/big"

	wallet "github.com/hashcloak/Meson-wallet-demo"
)

func DemoSetup(w *wallet.Wallet) {
	var passphrase string

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
}

func DemoSend(w *wallet.Wallet) (err error) {
	var passphrase string

	if len(w.Accounts()) == 0 {
		return fmt.Errorf("no accounts being set")
	}
	ac := w.Accounts()[0]
	addr := ac.Address
	fmt.Printf("Using account %v\n", addr)
	tx, err := wallet.GenerateTx(addr, addr, w.ChainID(), w.Endpoint())
	if err != nil {
		return
	}
	fmt.Print("Enter passphrase: ")
	_, err = fmt.Scanf("%s", &passphrase)
	if err != nil {
		return
	}
	signedTx, err := w.SignTxWithPassphrase(ac, passphrase, tx, big.NewInt(w.ChainID()))
	if err != nil {
		return
	}
	reply, err := w.Send(signedTx)
	if err != nil {
		return
	}
	fmt.Println(reply)
	return nil
}
