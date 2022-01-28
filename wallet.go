package wallet

import (
	"encoding/hex"
	"fmt"

	ks "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashcloak/Meson-plugin/pkg/common"
	config "github.com/hashcloak/Meson-wallet-demo/config"
	client "github.com/hashcloak/Meson/client"
)

type Wallet struct {
	*ks.KeyStore

	// wallet config
	config *config.Config
	// meson client
	client *client.Client
	// meson session
	session *client.Session
}

func New(CfgFile string) (wallet *Wallet, err error) {
	wallet = new(Wallet)

	// Setup wallet config
	wallet.config, err = config.LoadFile(CfgFile)
	if err != nil {
		return nil, err
	}

	// Setup wallet keystore
	wallet.KeyStore = ks.NewKeyStore(wallet.config.KSLocation, ks.StandardScryptN, ks.StandardScryptP)

	// Setup wallet client
	err = wallet.config.Meson.UpdateTrust()
	if err != nil {
		return nil, err
	}
	linkKey := client.AutoRegisterRandomClient(wallet.config.Meson)
	wallet.client, err = client.New(wallet.config.Meson, wallet.config.Service)
	if err != nil {
		return nil, err
	}

	// Setup wallet session
	wallet.session, err = wallet.client.NewSession(linkKey)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (wallet *Wallet) Send(tx *types.Transaction) (string, error) {
	blobByte, _ := tx.MarshalBinary()
	blobString := "0x" + hex.EncodeToString(blobByte)

	// serialize our transaction inside a eth kaetzpost request message
	req := common.NewRequest(wallet.config.Ticker, blobString)
	mesonRequest := req.ToJson()
	mesonService, err := wallet.session.GetService(wallet.config.Service)
	if err != nil {
		return "", fmt.Errorf("client error: %v", err)
	}
	reply, err := wallet.session.BlockingSendUnreliableMessage(mesonService.Name, mesonService.Provider, mesonRequest)
	if err != nil {
		return "", fmt.Errorf("send error: %v", err)
	}
	return common.ResponseFromJson(reply)
}

func (wallet *Wallet) ChainID() int64 {
	return wallet.config.ChainID
}

func (wallet *Wallet) Endpoint() string {
	return wallet.config.Endpoint
}

func (wallet *Wallet) Close() {
	wallet.client.Shutdown()
	for _, account := range wallet.Accounts() {
		_ = wallet.Lock(account.Address)
	}
}
