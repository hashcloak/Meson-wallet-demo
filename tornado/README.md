# Support of Tornado Cash

To use tornado cash in this demo, just replace the `cli.js` of tornado-cli with our `cli.js`.

## Install
1. Clone the tornado-cli repository
```bash
git clone https://github.com/tornadocash/tornado-cli
```

2. Replace cli.js
```bash
cp cli.js tornado-cli/cli.js
```

3. Install libraries
```bash
cd tornado-cli
npm install
```

## Usage
1. Execute meson-wallet-demo in a separate terminal (with `-l`).
```bash
./wallet -w wallet.toml -l
```
2. Setup `.env` file of tornado-cli. Note that there is no need to provide `PRIVATE_KEY`. You only need to provide `SENDER_ACCOUNT` that is used by the wallet.
3. Follow the usage of tornado-cli. For example,
```bash
node cli.js deposit ETH 0.1 --rpc https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161
```
4. Go to the wallet terminal and confirm the transaction by entering your password.
