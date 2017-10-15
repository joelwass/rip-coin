package rip

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

// DiskWallet is the struct to use for encoding when we save
type DiskWallet struct {
	Pub  []byte
	Priv []byte
}

// Wallet is an individual user
type Wallet struct {
	Key ecdsa.PrivateKey
}

// New creates a public and private key
func (w *Wallet) New() {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.Key = *priv
}

// Save writes the pem files to disk
func (w *Wallet) Save(dir string) {
	// Check if the path exists fir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		e := os.Mkdir(dir, 0700)
		if e != nil {
			fmt.Println(e)
		}
	}

	priv, pub := w.Encode()
	contents, err := json.Marshal(DiskWallet{
		Pub:  pub,
		Priv: priv,
	})
	if err != nil {
		fmt.Println("Error encoding wallet.")
	}

	err = ioutil.WriteFile(dir+"wallet.dat", contents, 0700)
	if err != nil {
		fmt.Println(err)
	}
}

// Encode s a key so it can be read from file
func (w *Wallet) Encode() ([]byte, []byte) {
	x509Encoded, _ := x509.MarshalECPrivateKey(&w.Key)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(&w.Key.PublicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncoded, pemEncodedPub
}

// Decode s a key from file
func (w *Wallet) Decode(pemEncoded []byte, pemEncodedPub []byte) {
	w.Key = *DecodePrivate(pemEncoded)
	w.Key.PublicKey = *DecodePublic(pemEncodedPub)
}

// DecodePrivate decodes the private key
func DecodePrivate(key []byte) *ecdsa.PrivateKey {
	block, _ := pem.Decode(key)
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	return privateKey
}

// DecodePublic decodes a public pem key
func DecodePublic(key []byte) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode(key)
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey
}

// Vote votes on a transaction, returns a new Vote
func (w *Wallet) Vote(t *Tx, vote bool) (*Vote, error) {
	_, pub := w.Encode()
	if !t.Verify() {
		return &Vote{Approval: false, Address: pub}, errors.New("Invalid transaction!")
	}

	return &Vote{Approval: vote, Address: pub}, nil
}

// LoadWallet retrieves a wallet
func LoadWallet() (*DiskWallet, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir

	dat, err := ioutil.ReadFile(dir + "/.rip/wallet.dat")
	if err != nil {
		return nil, err
	}
	w := DiskWallet{}
	json.Unmarshal(dat, &w)

	return &w, nil
}
