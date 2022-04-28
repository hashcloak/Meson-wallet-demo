# Meson-wallet-demo
A simple wallet that sends transactions through Meson.


## Build

```
cd cmd/wallet
go build .
```

## Usage

Make sure there is an existing katzenmint and Meson mixnet running. You will need to check for the configuration in `cmd/wallet/wallet.toml`.

```
Usage of ./wallet:
  -w string
        Path to the wallet config file (default "wallet.toml")
  -l bool
        Listen over tcp and serve transaction requests (default false)
```

If listen mode is switched off, the wallet sends a single test transaction. The following extra commands can be used.

```
  -c int
        Chain ID for specific ETH-based chain (default 5)
  -a string
        Address of the receiver (blank means sender itself) (default "")
  -v string
        Value to be transfered (default 10)
  -d string
        Data in hex format to be included (default "")
```
