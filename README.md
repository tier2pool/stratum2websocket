# Tier2Pool

A mining pool proxy tool.

## Build

### Build from source

```bash
make build
```

## Usage

- Server

```sh
./tier2pool server --ssl-certificate ./fullchain.pem --ssl-certificate-key ./privkey.pem --token password --redirect https://www.bing.com:443
```

- Client

```shell
./tier2pool client --server wss://example.com --pool pool.minexmr.com:4444 --token password
```

## License

[GNU General Public License v3.0](LICENSE).
