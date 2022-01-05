package wallet

import (
	"errors"

	"github.com/hashcloak/Meson-wallet-demo/pkg/ethers"
)

func GenerateTransaction(w *Wallet) (*string, error) {
	ethers, err := ethers.SetURLAndChainID(w.Config.RpcEndpoint)
	if err != nil {
		return nil, err
	}

	if ethers.ChainID.Int64() != int64(w.Config.ChainID) {
		return nil, errors.New("ChainIDs are not the same between rpcEndpoint and chainID flag")
	}
	rawTxn, err := ethers.GenerateSignedRawTxn(w.Config.PrivKey)
	if err != nil {
		return nil, err
	}
	return rawTxn, nil
}
