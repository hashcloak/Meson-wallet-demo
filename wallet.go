package wallet

import (
	"fmt"

	client "github.com/hashcloak/Meson-client"
	cConfig "github.com/hashcloak/Meson-client/config"
	"github.com/hashcloak/Meson-plugin/pkg/common"
	wConfig "github.com/hashcloak/Meson-wallet-demo/config"
)

type Wallet struct {
	// wallet config
	Config *wConfig.Config
	// meson client
	client *client.Client
	// meson session
	session *client.Session
}

func New(walletCfgFile string) (wallet *Wallet, err error) {
	wallet = new(Wallet)
	wallet.Config, err = wConfig.LoadFile(walletCfgFile)
	if err != nil {
		return nil, err
	}
	clientCfg, err := cConfig.LoadFile(wallet.Config.ClientCfgFile)
	if err != nil {
		return nil, err
	}
	err = cConfig.UpdateTrust(clientCfg)
	if err != nil {
		return nil, err
	}
	clientCfg, linkKey := client.AutoRegisterRandomClient(clientCfg)
	err = clientCfg.SaveConfig(wallet.Config.ClientCfgFile)
	if err != nil {
		return nil, err
	}
	wallet.client, err = client.New(wallet.Config.ClientCfgFile, wallet.Config.Service)
	if err != nil {
		return nil, err
	}
	wallet.session, err = wallet.client.NewSession(linkKey)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (wallet *Wallet) Send(rawTransactionBlob string) ([]byte, error) {
	// serialize our transaction inside a eth kaetzpost request message
	req := common.NewRequest(wallet.Config.Ticker, rawTransactionBlob)
	mesonRequest := req.ToJson()
	mesonService, err := wallet.session.GetService(wallet.Config.Service)
	if err != nil {
		return nil, fmt.Errorf("client error: %v", err)
	}
	return wallet.session.BlockingSendUnreliableMessage(mesonService.Name, mesonService.Provider, mesonRequest)
}

func (wallet *Wallet) Close() {
	wallet.session.Shutdown()
	wallet.client.Shutdown()
}
