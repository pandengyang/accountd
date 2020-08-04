package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
)

const (
	PRIKEY = "accountd_prikey.pem"
	PUBKEY = "accountd_pubkey.pem"
)

func main() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	privateKeyData, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	privateBlock := pem.Block{
		Type:  "esdsa private key",
		Bytes: privateKeyData,
	}

	privateKeyFile, err := os.Create(PRIKEY)
	if err != nil {
		panic(err)
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, &privateBlock)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.PublicKey

	publicKeyData, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}

	publicBlock := pem.Block{
		Type:  "ecdsa public key",
		Bytes: publicKeyData,
	}

	publicKeyFile, err := os.Create(PUBKEY)
	if err != nil {
		panic(err)
	}
	defer publicKeyFile.Close()

	err = pem.Encode(publicKeyFile, &publicBlock)
	if err != nil {
		panic(err)
	}
}
