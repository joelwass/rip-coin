# Rip Coin
The coin for friends ripping friends

## Setup
* install and setup GoLang https://golang.org/doc/install

## Setup Wallet
* `$ go get github.com/nathanjohnson320/rip-coin && go get -u github.com/ipfs/ipfs-update && ipfs-update install latest && ipfs init`
* `$ cd $GOPATH && cd src/github.com/nathanjohnson320/rip-coin && go run main.go`
* your wallet should have been created. your public and private keys can now be accessed with `$ cat ~/.rip/wallet.dat`

## For front - end, see rip-coin/README.md
