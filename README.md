# Meson-wallet-demo
A simple wallet that demonstrates how to use the Meson client library


## Dependecies

```
docker-compose
make
```

## Usage

Before starting to use this repo make sure there is an existing katzenpost mixnet running. You will need the ip address of one of the mixnet Authorities. A dockerized testnet can be found at: [hashcloak/Meson](https://github.com/hashcloak/Meson)


```
Usage of ./main:
  -c string
        Path to the server config file (default "katzenpost.toml")
  -chain int
        Chain ID for specific ETH-based chain (default 1)
  -pk string
        Private key used to sign the txn
  -rpc string
        Ethereum rpc endpoint (default "http://172.28.1.10:9545")
  -rt string
        Raw Transaction blob to send over the network
  -s string
        Service Name
  -t string
        Ticker
```

You can use either a private key for generating signed raw transactions or you can provide an already signed raw transaction with the `-rt` flag.


### Running:

Preflight checks:

- The dockerized testnet is running
- Make sure that a `./pk` file is present with the contents of your private key. 

```
make send-txn
```

### More details

If you would like to create your own command instead of using the `make` rule here is an example:


```bash
make build

RAW_TXN='0xf8640c843b9aca0083030d409400b1c66f34d680cb8bf82c64dcc1f39be5d6e77501802da03b274f8e63ce753e1ccdd03ac2d5e2595ef605335ed4962fe058eb667dbf9e6ba07c91420f9cb9805b18c6f25f84e530b35fca9eb45e4c3f6e6d624f53a3a76c40'

docker run \
  --mount type=bind,source=`pwd`/client.toml,target=/client.toml \
  --network nonvoting_testnet_nonvoting_test_net \
  hashcloak/meson-wallet /wallet \
  -c /client.toml \
  -t gor \
  -s gor \
  -chain 5 \
  -rt $RAW_TXN
```
