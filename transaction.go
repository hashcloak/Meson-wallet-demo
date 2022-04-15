package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GenerateTransaction queries RPC directly to get nonce, gasprice, gaslimit
// We recommend using (*Wallet).FillTx() instead for more privacy
func GenerateTransaction(from common.Address, to common.Address, value *big.Int, data []byte, chainID int64, rpcEndpoint string) (*types.Transaction, error) {
	ethclient, err := ethclient.Dial(rpcEndpoint)
	if err != nil {
		return nil, err
	}
	defer ethclient.Close()
	recvChainID, err := ethclient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	if recvChainID.Int64() != chainID {
		return nil, fmt.Errorf("chain ID mismatch")
	}
	nonce, err := ethclient.PendingNonceAt(context.Background(), from)
	if err != nil {
		return nil, err
	}
	gasPrice, err := ethclient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	gasEstimate, err := ethclient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:    &to,
		Value: value,
		Data:  data,
	})
	if err != nil {
		return nil, err
	}
	tx := types.NewTx(&types.AccessListTx{
		ChainID:  recvChainID,
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasEstimate,
		GasPrice: gasPrice,
		Data:     data,
	})

	return tx, nil
}
