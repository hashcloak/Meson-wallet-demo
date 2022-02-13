package wallet

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionRequest struct {
	ChainID int64
	To      string
	Value   big.Int
	Data    string
}

func ProcessRequest(w *Wallet, request TransactionRequest) (reply string, err error) {
	account := w.UiSelectAccount()
	tx, err := GenerateTransaction(
		account.Address,
		common.HexToAddress(request.To),
		&request.Value,
		common.FromHex(request.Data),
		request.ChainID,
		w.Endpoint(request.ChainID),
	)
	if err != nil {
		return "", err
	}
	w.uiConfirm(account.Address, tx)
	passphrase, err := w.uiPassphrase()
	if err != nil {
		return "", err
	}
	signedTx, err := w.SignTxWithPassphrase(account, passphrase, tx, tx.ChainId())
	if err != nil {
		return "", err
	}
	reply, err = w.Send(signedTx)
	if err != nil {
		return "", err
	}
	return reply, nil
}

func TransactionHandler(w *Wallet, resp http.ResponseWriter, req *http.Request) {
	request := new(TransactionRequest)
	err := json.NewDecoder(req.Body).Decode(request)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	reply, err := ProcessRequest(w, *request)
	if err != nil {
		fmt.Printf("Failed to process: %v\n\n", err)
		http.Error(resp, "Failed to process transaction", http.StatusInternalServerError)
		return
	}
	fmt.Printf("%s\n\n", reply)
	_, _ = resp.Write([]byte(reply))
}
