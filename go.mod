module github.com/hashcloak/Meson-wallet-demo

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.8
	github.com/hashcloak/Meson-client v0.0.0-20191203233537-dc6d4bdd3049
	github.com/hashcloak/Meson-plugin v0.0.0-00010101000000-000000000000 // indirect
	github.com/hashcloak/Meson/plugin v0.0.0-20191130194144-a50c00894e10
	github.com/katzenpost/client v0.0.3-0.20191109165001-aa02bb21ca21
)

replace github.com/hashcloak/Meson-plugin => ../Meson-plugin
