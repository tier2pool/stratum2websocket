# Tier2Pool

![GitHub last commit](https://img.shields.io/github/last-commit/tier2pool/tier2pool?style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/tier2pool/tier2pool?style=flat-square)
![GitHub license](https://img.shields.io/github/license/tier2pool/tier2pool?style=flat-square)

A mining pool proxy tool.

## Build

I use Ubuntu as a demo.

```shell
sudo update
sudo apt install git make snapd -y
sudo snap install --classic
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
./tier2pool client --server wss://example.com --pool pool.minexmr.com:4444 --token password
```

## License

[GNU General Public License v3.0](LICENSE).
