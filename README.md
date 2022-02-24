# Tier2Pool

![GitHub last commit](https://img.shields.io/github/last-commit/tier2pool/tier2pool?style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/tier2pool/tier2pool?style=flat-square)
![GitHub license](https://img.shields.io/github/license/tier2pool/tier2pool?style=flat-square)

A mining pool proxy tool, support BTC, ETH, ETC, XMR mining pool, etc.

## Build

```shell
# Ubuntu or Debian
sudo apt update
sudo apt install git build-essential snapd -y
sudo snap install go --classic
git clone https://github.com/tier2pool/tier2pool
cd tier2pool
make build
# make build_windows_amd64
# make build_linux_amd64
cd build && ls
```

## Usage

### Screen

The client is listening on `127.0.0.1:1234` by default.

```shell
# Linux
./build/tier2pool --help

# Server
screen -S tier2pool-server
./build/tier2pool server --ssl-certificate ./fullchain.pem --ssl-certificate-key ./privkey.pem --token password --redirect https://www.bing.com:443

# Client
screen -S tier2pool-client
# ETH
./build/tier2pool client --server wss://example.com --pool tls://us1.ethermine.org:5555 --token password
# XMR
./build/tier2pool client --server wss://example.com --pool tcp://pool.minexmr.com:4444 --token password
```

In addition, you can use the `--allowed-pool` parameter to limit the available mining pools.

```shell
./build/tier2pool server --allowed-pool tls://us1.ethermine.org:5555 --allowed tcp://pool.minexmr.com:4444 ...
```

### Docker

```shell
make build_image

# Server
docker run \
  --name tier2pool-server \
  -p 443:443 \
  --restart=on-failure:3 \
  -dit tier2pool/tier2pool:latest \
  server \
  --ssl-certificate ./fullchain.pem \
  --ssl-certificate-key ./privkey.pem \
  --token password \
  --redirect https://www.bing.com:443
  
# Client
docker run \
  --name tier2pool-client \
  -p 1234:1234 \
  --restart=on-failure:3 \
  -dit tier2pool/tier2pool:latest \
  client \
  --server wss://example.com \
  --pool tls://us1.ethermine.org:5555 \
  --token password
```

## TODO

- [x] Encryption and obfuscation of transmitted data.
- [ ] Display hash rate and submission information.
- [ ] Display connection status between miner and mining pool.

## Donate

### ETH

You can donate any amount to me in the Ethereum `Mainnet`, `Polygon` or `BEP20` to support my work.

```diff
+ 0x000000A52a03835517E9d193B3c27626e1Bc96b1
```

### XMR

```diff
+ 84TZwzCfHhkZ43JzygNqaN5ke6t3uRSD32rofAhV19jB1VNzDnkaciWN7c7tfqFvKt95f4Y6jyEecWzsnUHi1koZNqBveJb
```

## License

[GNU General Public License v3.0](LICENSE).
