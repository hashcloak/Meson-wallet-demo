package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	config "github.com/hashcloak/Meson-wallet-demo/config"
	client "github.com/hashcloak/Meson/client"
	"github.com/hashcloak/Meson/plugin/pkg/command"
	"github.com/hashcloak/Meson/plugin/pkg/common"
	ks "github.com/hashcloak/go-ethereum/accounts/keystore"
	ethcommon "github.com/hashcloak/go-ethereum/common"
	"github.com/hashcloak/go-ethereum/common/hexutil"
	"github.com/hashcloak/go-ethereum/core/types"
)

const mesonService = "meson"

type Wallet struct {
	*ks.KeyStore

	// wallet config
	config *config.Config
	// meson client
	client *client.Client
	// meson session
	session *client.Session
}

func New(CfgFile string) (w *Wallet, err error) {
	w = new(Wallet)

	// Setup wallet config
	w.config, err = config.LoadFile(CfgFile)
	if err != nil {
		return nil, err
	}

	// Setup wallet keystore
	w.KeyStore = ks.NewKeyStore(w.config.KSLocation, ks.StandardScryptN, ks.StandardScryptP)

	// Setup wallet client
	err = w.config.Meson.UpdateTrust()
	if err != nil {
		return nil, err
	}
	linkKey := client.AutoRegisterRandomClient(w.config.Meson)
	w.client, err = client.New(w.config.Meson, mesonService)
	if err != nil {
		return nil, err
	}

	// Setup wallet session
	w.session, err = w.client.NewSession(linkKey)
	if err != nil {
		return nil, err
	}

	// Initialize with UI
	err = w.uiSetup()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Wallet) throughMeson(req []byte) (string, error) {
	service, err := w.session.GetService(mesonService)
	if err != nil {
		return "", fmt.Errorf("client error: %v", err)
	}
	reply, err := w.session.BlockingSendUnreliableMessage(service.Name, service.Provider, req)
	if err != nil {
		return "", fmt.Errorf("send error: %v", err)
	}
	return common.ResponseFromJson(reply)
}

func (w *Wallet) FillTx(from ethcommon.Address, to ethcommon.Address, value *big.Int, data string, chainID int64) (*types.Transaction, error) {
	payload, err := json.Marshal(command.EthQueryRequest{
		From:  from.Hex(),
		To:    to.Hex(),
		Value: value,
		Data:  data,
	})
	if err != nil {
		return nil, err
	}
	req := common.NewRequest(command.EthQuery, w.Ticker(chainID), payload).ToJson()
	resp, err := w.throughMeson(req)
	if err != nil {
		return nil, err
	}
	response := new(command.EthQueryResponse)
	err = json.Unmarshal([]byte(resp), response)
	if err != nil {
		return nil, err
	}
	nonce, err := hexutil.DecodeUint64(response.Nonce)
	if err != nil {
		return nil, err
	}
	gasLimit, err := hexutil.DecodeUint64(response.GasLimit)
	if err != nil {
		return nil, err
	}
	gasPrice, err := hexutil.DecodeBig(response.GasPrice)
	if err != nil {
		return nil, err
	}
	dataByte, err := hexutil.Decode(data)
	if err != nil && data != "" {
		return nil, err
	}
	tx := types.NewTx(&types.AccessListTx{
		ChainID:  big.NewInt(chainID),
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     dataByte,
	})
	return tx, nil
}

func (w *Wallet) SendTx(tx *types.Transaction) (string, error) {
	blobByte, _ := tx.MarshalBinary()
	blobString := "0x" + hex.EncodeToString(blobByte)
	payload, err := json.Marshal(command.PostTransactionRequest{TxHex: blobString})
	if err != nil {
		return "", err
	}
	req := common.NewRequest(command.PostTransaction, w.Ticker(tx.ChainId().Int64()), payload).ToJson()
	return w.throughMeson(req)
}

func (w *Wallet) SendHexSignedTx(signedTx string, chainId int64) (string, error) {
	payload, err := json.Marshal(command.PostTransactionRequest{TxHex: signedTx})
	if err != nil {
		return "", err
	}
	req := common.NewRequest(command.PostTransaction, w.Ticker(chainId), payload).ToJson()
	return w.throughMeson(req)
}

func (w *Wallet) Ticker(chainID int64) string {
	return w.config.Chain[fmt.Sprint(chainID)].Ticker
}

func (w *Wallet) Endpoint(chainID int64) string {
	return w.config.Chain[fmt.Sprint(chainID)].Endpoint
}

func (w *Wallet) Close() {
	w.client.Shutdown()
	for _, account := range w.Accounts() {
		_ = w.Lock(account.Address)
	}
}
