package wallet

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/hashcloak/go-ethereum/common"
)

type TransactionRequest struct {
	ChainID int64
	To      string
	Value   big.Int
	Data    string
	// for bitcoin transaction
	UnspentTxHash string
	UnspentIndex  uint
	Left          int64
	Chain         *chaincfg.Params
}

func ProcessRequest(w *Wallet, request TransactionRequest) (reply string, err error) {
	fmt.Println("Select account...")
	account := w.UiSelectAccount()
	if request.Chain != nil {
		// bitcoin transaction
		fmt.Println("Got bitcoin transaction")
		passphrase, err := w.uiPassphrase()
		if err != nil {
			return "", err
		}
		key, err := w.UnlockUnsafe(account, passphrase)
		if err != nil {
			return "", err
		}
		pvBytes := key.PrivateKey.D.Bytes()
		bPrivKey, bPubKey := btcec.PrivKeyFromBytes(pvBytes)

		// create redeem script for 1 of 1 multi-sig
		builder := txscript.NewScriptBuilder()
		builder.AddOp(txscript.OP_1)
		builder.AddData(bPubKey.SerializeCompressed())
		builder.AddOp(txscript.OP_1)
		builder.AddOp(txscript.OP_CHECKMULTISIG)
		redeemScript, err := builder.Script()
		if err != nil {
			return "", err
		}
		userAddress, err := btcutil.NewAddressScriptHash(
			redeemScript, request.Chain,
		)
		decUserAddress, err := btcutil.DecodeAddress(userAddress.String(), request.Chain)
		if err != nil {
			return "", err
		}

		decUserAddressByte, err := txscript.PayToAddrScript(decUserAddress)
		if err != nil {
			return "", err
		}
		value := request.Value.Int64()
		left := request.Left
		fmt.Printf("Using account %v\n", userAddress)
		fmt.Println("====================================")
		fmt.Printf("To: %v\n", request.To)
		fmt.Printf("Value: %v\n", value)
		fmt.Printf("Left: %v\n", left)
		fmt.Println("====================================")

		// prepare transaction
		tx := wire.NewMsgTx(wire.TxVersion)
		utxoHash, err := chainhash.NewHashFromStr(request.UnspentTxHash)
		if err != nil {
			return "", err
		}

		// and add the index of the UTXO
		inPoint := wire.NewOutPoint(utxoHash, uint32(request.UnspentIndex))
		txIn := wire.NewTxIn(inPoint, nil, nil)

		tx.AddTxIn(txIn)

		// adding the output to tx
		targetAddress, err := btcutil.DecodeAddress(request.To, request.Chain)
		if err != nil {
			return "", err
		}
		destinationAddrByte, err := txscript.PayToAddrScript(targetAddress)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(value, destinationAddrByte)
		tx.AddTxOut(txOut)
		txOut2 := wire.NewTxOut(left, decUserAddressByte)
		tx.AddTxOut(txOut2)

		// signing the tx
		sig1, err := txscript.RawTxInSignature(tx, 0, redeemScript, txscript.SigHashAll, bPrivKey)
		if err != nil {
			return "", err
		}

		signature := txscript.NewScriptBuilder()
		signature.AddOp(txscript.OP_FALSE).AddData(sig1)
		signature.AddData(redeemScript)
		signatureScript, err := signature.Script()
		if err != nil {
			return "", err
		}

		tx.TxIn[0].SignatureScript = signatureScript

		var signedTx bytes.Buffer
		tx.Serialize(&signedTx)

		hexSignedTx := hex.EncodeToString(signedTx.Bytes())
		reply, err = w.SendHexSignedTx(hexSignedTx, request.ChainID)
		if err != nil {
			return "", err
		}
		return reply, nil
	}
	fmt.Println("Selected account: ", account)
	tx, err := w.FillTx(
		account.Address,
		common.HexToAddress(request.To),
		&request.Value,
		request.Data,
		request.ChainID,
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
	reply, err = w.SendTx(signedTx)
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
