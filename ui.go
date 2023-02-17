package wallet

import (
	"encoding/hex"
	"fmt"
	"log"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/term"
)

func (w *Wallet) uiSetup() error {
	var passphrase string

	if len(w.Accounts()) == 0 {
		fmt.Print("Init passphrase: ")
		_, err := fmt.Scanf("%s", &passphrase)
		if err != nil {
			return err
		}
		if len(passphrase) < 4 {
			return fmt.Errorf("passphrase has to be at least 4 characters")
		}
		ac, err := w.NewAccount(passphrase)
		if err != nil {
			return err
		}
		fmt.Printf("Account created with address %v\n", ac.Address)
	}
	return nil
}

func (w *Wallet) UiSelectAccount() accounts.Account {
	if len(w.Accounts()) == 0 {
		panic("No accounts being set")
	}
	return w.Accounts()[0]
}

func (w *Wallet) uiConfirm(fromAddress common.Address, tx *types.Transaction) {
	fmt.Printf("Using account %v\n", fromAddress)
	log.Println("====================================")
	fmt.Printf("To: %v\n", tx.To())
	fmt.Printf("Value: %v\n", tx.Value())
	fmt.Printf("Data: %v\n", hex.EncodeToString(tx.Data()))
	fmt.Printf("Chain: %v\n", tx.ChainId())
	fmt.Printf("GasLimit: %v\n", tx.Gas())
	fmt.Printf("GasPrice: %v\n", tx.GasPrice())
	log.Println("====================================")
}

func (w *Wallet) uiPassphrase() (string, error) {
	fmt.Print("Enter passphrase: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	log.Println("")
	return string(password), nil
}
