# Tier2Pool

![GitHub last commit](https://img.shields.io/github/last-commit/tier2pool/tier2pool?style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/tier2pool/tier2pool?style=flat-square)
![GitHub license](https://img.shields.io/github/license/tier2pool/tier2pool?style=flat-square)

A mining pool proxy tool, support BTC, ETH, ETC, XMR mining pool, etc.

## Build

I use Ubuntu as a demo.

```shell
sudo apt update
sudo apt install git make snapd -y
sudo snap install go --classic
git clone https://github.com/tier2pool/tier2pool
cd tier2pool
make build
# make build_windows_amd64
# make build_linux_amd64
cd build && ls
```

## Usage

- Server

```shell
./tier2pool server --ssl-certificate ./fullchain.pem --ssl-certificate-key ./privkey.pem --token password --redirect https://www.bing.com:443
```

- Client

The default address is `127.0.0.1:1234`.

```shell
# ETH
./tier2pool client --server wss://example.com --pool tls://us1.ethermine.org:5555 --token password

# XMR
./tier2pool client --server wss://example.com --pool tcp://pool.minexmr.com:4444 --token password
```


## TODO

- [x] Encryption and obfuscation of transmitted data.
- [ ] Display hash rate and submission information.
- [ ] Display connection status between miner and mining pool.

## Donate

### ETH

You can donate any amount to me in the Ethereum `Mainnet` or `Polygon` to support my work.

```diff
+ 0x000000A52a03835517E9d193B3c27626e1Bc96b1
```

### XMR

```diff
+ 84TZwzCfHhkZ43JzygNqaN5ke6t3uRSD32rofAhV19jB1VNzDnkaciWN7c7tfqFvKt95f4Y6jyEecWzsnUHi1koZNqBveJb
```

## License

[GNU General Public License v3.0](LICENSE).
