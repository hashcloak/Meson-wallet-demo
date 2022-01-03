package wallet

import (
	"errors"

	"github.com/hashcloak/Meson-wallet-demo/pkg/ethers"
)

func GenerateTransaction(wallet *Wallet) (*string, error) {
	w := wallet.Config.Optional
	ethers, err := ethers.SetURLAndChainID(w.RpcEndpoint)
	if err != nil {
		return nil, err
	}

	if ethers.ChainID.Int64() != int64(w.ChainID) {
		return nil, errors.New("ChainIDs are not the same between rpcEndpoint and chainID flag")
	}
	rawTxn, err := ethers.GenerateSignedRawTxn(w.PrivKey)
	if err != nil {
		return nil, err
	}
	return rawTxn, nil
}
