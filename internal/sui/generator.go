package sui

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/say8hi/addrforge/internal/util"
	"golang.org/x/crypto/blake2b"
)

type Wallet struct {
	PrivateKey string
	Address    string
}

// GenerateWallet creates a new Sui wallet
func GenerateWallet() (*Wallet, error) {
	// Generate 32 bytes of entropy for seed
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return nil, err
	}

	// Create Ed25519 key from seed
	priv := ed25519.NewKeyFromSeed(seed)
	pub := priv.Public().(ed25519.PublicKey)

	// Add signature scheme flag (0x00 for Ed25519) BEFORE public key
	pubKeyWithScheme := append([]byte{0x00}, pub...)

	// Hash with BLAKE2b (256 bits = 32 bytes output)
	hasher, err := blake2b.New256(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create BLAKE2b hasher: %w", err)
	}
	hasher.Write(pubKeyWithScheme)
	hash := hasher.Sum(nil)

	address := "0x" + hex.EncodeToString(hash)

	// Sui private key format: 33 bytes = 1-byte scheme flag + 32-byte private key seed
	// For Ed25519: scheme flag is 0x00 (comes FIRST)
	privKeyWith33Bytes := append([]byte{0x00}, seed...)

	// Encode as Bech32 with 'suiprivkey' prefix
	privKeyBech32, err := util.EncodeBech32("suiprivkey", privKeyWith33Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key: %w", err)
	}

	return &Wallet{
		PrivateKey: privKeyBech32,
		Address:    address,
	}, nil
}

// ValidatePrefix validates that the prefix contains only valid hex characters
func ValidatePrefix(prefix string) error {
	// Remove 0x if present
	if len(prefix) > 2 && prefix[:2] == "0x" {
		prefix = prefix[2:]
	}

	// Check hex characters
	for _, char := range prefix {
		if !((char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'f') ||
			(char >= 'A' && char <= 'F')) {
			return fmt.Errorf("contains invalid hex character: %c", char)
		}
	}
	return nil
}
