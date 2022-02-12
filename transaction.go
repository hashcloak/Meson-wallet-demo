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

func GenerateTx(from common.Address, to common.Address, value int64, data []byte, chainID int64, rpcEndpoint string) (*types.Transaction, error) {
	ethclient, err := ethclient.Dial(rpcEndpoint)
	if err != nil {
		return nil, err
	}
	recvChainID, err := ethclient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	if recvChainID.Int64() != chainID {
		return nil, fmt.Errorf("chain ID mismatch")
	}
	/*
	 * This is somewhere we need more privacy protection in the future
	 */
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
		Value: big.NewInt(value),
		Data:  data,
	})
	if err != nil {
		return nil, err
	}
	tx := types.NewTx(&types.AccessListTx{
		ChainID:  recvChainID,
		Nonce:    nonce,
		To:       &to,
		Value:    big.NewInt(value),
		Gas:      gasEstimate * 11 / 10,
		GasPrice: gasPrice,
		Data:     data,
	})

	return tx, nil
}
