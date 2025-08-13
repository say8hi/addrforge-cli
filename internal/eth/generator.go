package eth

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	PrivateKey string
	Address    string
}

func GenerateWallet() (*Wallet, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey).Hex()

	privHex := hex.EncodeToString(crypto.FromECDSA(privKey))

	return &Wallet{
		PrivateKey: privHex,
		Address:    address,
	}, nil
}
