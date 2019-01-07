
# lndate

A sandbox project about lnd client

# Environments
```
macOS v10.14.2
go version go1.11.2 darwin/amd64
btcd version 0.12.0-beta
lnd version 0.5.1-beta commit=v0.5.1-beta-74-gc9a5593a6468bdc06e0af3792017e0c798a5b37b-dirty
```

# How to use
1) run lnd in background
```sh
$ btcd --txindex --simnet --rpcuser=sonokko --rpcpass=1qazxsw23edcvfr4
$ lnd --rpclisten=localhost:10009 --listen=localhost:10011 --restlisten=localhost:8001 --datadir=data  --debuglevel=info --bitcoin.simnet --bitcoin.active --bitcoin.node=btcd --btcd.rpcuser=sonokko --btcd.rpcpass=1qazxsw23edcvfr4 --no-macaroons
```

2) Run
```
$ go run main.go
```