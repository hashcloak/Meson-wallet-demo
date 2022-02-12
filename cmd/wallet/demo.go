package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	wallet "github.com/hashcloak/Meson-wallet-demo"
)

func DemoSetup(w *wallet.Wallet) common.Address {
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
	} else {
		fmt.Printf("Using account %v\n", w.Accounts()[0].Address)
	}
	return w.Accounts()[0].Address
}

func DemoSend(w *wallet.Wallet, tx *types.Transaction) (err error) {
	var passphrase string

	ac := w.Accounts()[0]
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
