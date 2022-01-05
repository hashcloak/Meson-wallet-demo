package ethers

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Ethers basic ethereum node interface
type Ethers struct {
	ChainID   *big.Int
	ethclient *ethclient.Client
}

// SetURLAndChainID init for ethers
func SetURLAndChainID(rpcEndpoint string) (Ethers, error) {
	ethclient, err := ethclient.Dial(rpcEndpoint)
	if err != nil {
		return Ethers{}, err
	}
	chainID, err := ethclient.ChainID(context.Background())
	return Ethers{chainID, ethclient}, err
}

// GenerateSignedRawTxn just signs a txn with
func (e *Ethers) GenerateSignedRawTxn(privKey string) (*string, error) {

	key, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	nonce, err := e.ethclient.PendingNonceAt(
		context.Background(),
		crypto.PubkeyToAddress(key.PublicKey),
	)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.ethclient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	to := common.HexToAddress(crypto.PubkeyToAddress(key.PublicKey).Hex())
	tx := types.NewTransaction(
		nonce,
		to,
		big.NewInt(123),
		uint64(21000),
		gasPrice,
		[]byte(""),
	)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(e.ChainID), key)
	if err != nil {
		return nil, err
	}
	raw, _ := signedTx.MarshalBinary()
	txn := "0x" + hex.EncodeToString(raw)
	return &txn, nil
}
