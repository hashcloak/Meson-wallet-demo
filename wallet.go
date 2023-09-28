package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	ks "github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	config "github.com/hashcloak/Meson-wallet-demo/config"
	client "github.com/hashcloak/Meson/client"
	"github.com/hashcloak/Meson/plugin/pkg/command"
	"github.com/hashcloak/Meson/plugin/pkg/common"
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

// An json request abstraction.
type jsonrpcRequest struct {
	ID uint `json:"id"`
	// Indicates which version of JSON RPC to use
	// Since all networks support JSON RPC 2.0, 1.0
	// this attribute is a constant
	JSONRPC string `json:"jsonrpc"`
	// Which method you want to call
	METHOD string `json:"method"`
	// Params for the method you want to call
	Params interface{} `json:"params"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RPCResponse struct {
	Version string    `json:"jsonrpc,omitempty"`
	ID      uint      `json:"id,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
	Result  string    `json:"result,omitempty"`
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

	// uncomment to use EthQuery command
	// payload, err := json.Marshal(command.EthQueryRequest{
	// 	From:  from.Hex(),
	// 	To:    to.Hex(),
	// 	Value: value,
	// 	Data:  data,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	//req := common.NewRequest(command.EthQuery, w.Ticker(chainID), payload).ToJson()
	// resp, err := w.throughMeson(req)
	// if err != nil {
	// 	return nil, err
	// }
	// response := new(command.EthQueryResponse)
	// err = json.Unmarshal([]byte(resp), response)
	// if err != nil {
	// 	return nil, err
	// }

	// Query gas with DirectPost command
	var dir_response []RPCResponse
	nonceRequest := jsonrpcRequest{
		ID:      1,
		JSONRPC: "2.0",
		METHOD:  "eth_getTransactionCount",
		Params:  []string{from.Hex(), "pending"},
	}
	gasPriceRequest := jsonrpcRequest{
		ID:      2,
		JSONRPC: "2.0",
		METHOD:  "eth_gasPrice",
	}
	param := map[string]interface{}{
		"from":  from.Hex(),
		"to":    to.Hex(),
		"value": fmt.Sprintf("0x%x", value),
	}
	if data != "" {
		param["data"] = data
	}
	gasEstimateRequest := jsonrpcRequest{
		ID:      3,
		JSONRPC: "2.0",
		METHOD:  "eth_estimateGas",
		Params:  []interface{}{param},
	}
	payload, err := json.Marshal([]jsonrpcRequest{
		nonceRequest,
		gasPriceRequest,
		gasEstimateRequest,
	})
	if err != nil {
		return nil, err
	}
	req := common.NewRequest(0x01, w.Ticker(chainID), payload).ToJson()
	resp, err := w.throughMeson(req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(resp), &dir_response)
	if err != nil {
		return nil, err
	}
	// Check if response type is error
	for _, pl := range dir_response {
		if pl.Error != nil {
			return nil, fmt.Errorf("error code: %d, msg: %s", pl.Error.Code, pl.Error.Message)
		}
	}
	response := command.EthQueryResponse{
		Nonce:    dir_response[0].Result,
		GasPrice: dir_response[1].Result,
		GasLimit: dir_response[2].Result,
	}
	// Comment above to use EthQuery command

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

	// Uncomment to use PostTransaction command
	// payload, err := json.Marshal(command.PostTransactionRequest{TxHex: blobString})
	// if err != nil {
	// 	return "", err
	// }
	// req := common.NewRequest(command.PostTransaction, w.Ticker(tx.ChainId().Int64()), payload).ToJson()

	// Send tx with DirectPost command
	payload, err := json.Marshal(jsonrpcRequest{
		ID:      1,
		JSONRPC: "2.0",
		METHOD:  "eth_sendRawTransaction",
		Params:  []string{blobString},
	})
	if err != nil {
		return "", err
	}
	req := common.NewRequest(0x01, w.Ticker(tx.ChainId().Int64()), payload).ToJson()
	// Comment above to use PostTransaction command

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
