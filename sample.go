package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GenerateTx(from, to common.Address, chainID int64, rpcEndpoint string) (*types.Transaction, error) {
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
	tx := types.NewTx(&types.AccessListTx{
		ChainID:  recvChainID,
		Nonce:    nonce,
		To:       &to,
		Value:    big.NewInt(10),
		Gas:      25000,
		GasPrice: gasPrice,
		Data:     common.FromHex(""),
	})

	return tx, nil
}
