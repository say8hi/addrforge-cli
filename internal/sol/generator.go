package sol

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/say8hi/addrforge/internal/util"
)

type Wallet struct {
	PrivateKey string
	Address    string
}

func GenerateWallet() (*Wallet, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	privBytes := append(priv.Seed(), pub...)

	return &Wallet{
		PrivateKey: util.Base58Encode(privBytes),
		Address:    util.Base58Encode(pub),
	}, nil
}
