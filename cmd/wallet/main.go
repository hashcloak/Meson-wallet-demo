package main

import (
	"flag"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	wallet "github.com/hashcloak/Meson-wallet-demo"
)

func main() {
	walletCfgFile := flag.String("w", "wallet.toml", "Wallet config file")
	setReceiver := flag.String("a", "", "Address of the receiver")
	setAmount := flag.Int64("v", 10, "Value transfered")
	setData := flag.String("d", "", "Data appended")
	flag.Parse()

	w, err := wallet.New(*walletCfgFile)
	if err != nil {
		panic(err)
	}
	sender := DemoSetup(w)
	receiver := sender // default receiver
	if *setReceiver != "" {
		receiver = common.HexToAddress(*setReceiver)
	}
	tx, err := wallet.GenerateTx(
		sender,
		receiver,
		*setAmount,
		common.FromHex(*setData),
		w.ChainID(),
		w.Endpoint(),
	)
	if err != nil {
		panic(err)
	}
	err = DemoSend(w, tx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done. Shutting down.")
	w.Close()
}
