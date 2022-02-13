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

func (w *Wallet) Send(tx *types.Transaction) (string, error) {
	blobByte, _ := tx.MarshalBinary()
	blobString := "0x" + hex.EncodeToString(blobByte)

	// serialize our transaction inside a eth kaetzpost request message
	req := common.NewRequest(w.Ticker(tx.ChainId().Int64()), blobString).ToJson()
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
