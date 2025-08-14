# AddrForge

![addrforge-logo](./logo.png)

🔧 **High-performance vanity address generator for Ethereum, Solana, and Sui blockchains**

AddrForge is a fast, multi-threaded CLI tool that generates cryptocurrency addresses with custom prefixes. Built with Go for maximum performance and efficiency.

## Features

- ⚡ **Multi-threaded generation** - Configurable worker pools for optimal performance
- 🎯 **Multiple prefix search** - Search for up to 5 different prefixes simultaneously  
- 🔤 **Case-sensitive/insensitive matching** - Flexible matching options
- 📊 **Real-time statistics** - Live TPS (Transactions Per Second) monitoring
- 💾 **File output** - Save results to file instead of console
- 🛡️ **Input validation** - Validates prefixes for each blockchain format
- 🔧 **Clean architecture** - Well-structured, maintainable code

## Supported Blockchains

- **Ethereum** - Generate addresses with hex prefixes (0x...)
- **Solana** - Generate addresses with Base58 prefixes
- **Sui** - Generate addresses with hex prefixes (0x...)

## Installation

### From Source

```bash
git clone https://github.com/say8hi/addrforge-cli.git
cd addrforge
go build -o addrforge
```

### Prerequisites

- Go 1.21+
- Dependencies will be downloaded automatically via `go mod`

## Quick Start

### Generate Ethereum Address

```bash
# Generate address starting with "dead"
./addrforge eth 0xdead

# Multiple prefixes
./addrforge eth 0xdead 0xbeef 0xcafe

# Case-insensitive search with 8 workers
./addrforge eth -i -w 8 0xDEAD
```

### Generate Solana Address

```bash
# Generate address starting with "Sol"  
./addrforge sol Sol

# Save results to file
./addrforge sol -o results.txt Sol1 Sol2
```

### Generate Sui Address

```bash
# Generate address starting with "ace"
./addrforge sui 0xace

# Multiple prefixes with case-insensitive matching
./addrforge sui -i -w 12 0xACE 0xFACE
```

## Command Reference

### Global Commands

```bash
addrforge [command] [flags]
```

### Ethereum Generation

```bash
addrforge eth <prefix1> [prefix2] [prefix3] [prefix4] [prefix5] [flags]
```

**Flags:**
- `-w, --workers int` - Number of worker threads (default: 4)
- `-i, --ignore-case` - Case-insensitive prefix matching
- `-o, --output string` - Output file for results
- `-h, --help` - Help for eth command

### Solana Generation

```bash
addrforge sol <prefix1> [prefix2] [prefix3] [prefix4] [prefix5] [flags]
```

**Flags:**
- `-w, --workers int` - Number of worker threads (default: 4)
- `-i, --ignore-case` - Case-insensitive prefix matching  
- `-o, --output string` - Output file for results
- `-h, --help` - Help for sol command

### Sui Generation

```bash
addrforge sui <prefix1> [prefix2] [prefix3] [prefix4] [prefix5] [flags]
```

**Flags:**
- `-w, --workers int` - Number of worker threads (default: 4)
- `-i, --ignore-case` - Case-insensitive prefix matching
- `-o, --output string` - Output file for results
- `-h, --help` - Help for sui command

## Output Format

When a matching address is found, AddrForge outputs:

**Ethereum/Sui:**
```
MATCH FOUND!
Address: 0xdeadbeef1234567890abcdef1234567890abcdef
Private: a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef12345678
Worker: 3
```

**Solana:**
```
MATCH FOUND!
Address: SoLxxx1234567890abcdef1234567890abcdef1234
Private: 3xK9L4m7N8p2Q5r6S9t0U1v4W7x0Y3z6A9b2C5d8E1f4
Worker: 2
```

## Performance Tips

### Optimal Worker Count

- **CPU-bound**: Set workers = CPU cores × 1.5-2
- **Example**: 8-core CPU → use 12-16 workers
- **Test different values** to find your system's optimum

### Prefix Difficulty

The difficulty increases exponentially with prefix length:

| Prefix Length | Ethereum/Sui (hex) | Solana (base58) |
|---------------|---------------------|-----------------|
| 1 character   | ~16 attempts        | ~58 attempts    |
| 2 characters  | ~256 attempts       | ~3,364 attempts |
| 3 characters  | ~4,096 attempts     | ~195,112 attempts|
| 4 characters  | ~65,536 attempts    | ~11,316,496 attempts|

### System Requirements

- **RAM**: ~50MB base + ~1MB per worker
- **CPU**: More cores = better performance
- **Storage**: Minimal (only for output files)

## Example Sessions

### Ethereum Vanity Address

```bash
$ ./addrforge eth -w 8 0xcafe

Searching for Ethereum address with prefixes: [cafe]
TPS: 45231     | Total: 1847291   | Errors: 0

MATCH FOUND!
Address: 0xcafe123456789abcdef0123456789abcdef012345
Private: 1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
Worker: 5
```

### Solana Vanity Address

```bash
$ ./addrforge sol -i Sol

Searching for Solana address with prefixes: [Sol]
Using case-insensitive matching
TPS: 38456     | Total: 892341    | Errors: 0

MATCH FOUND!
Address: SoLxxx1234567890abcdef1234567890abcdef1234
Private: 3xK9L4m7N8p2Q5r6S9t0U1v4W7x0Y3z6A9b2C5d8E1f4
Worker: 2
```

### Sui Vanity Address

```bash
$ ./addrforge sui -w 12 0xface

Searching for Sui address with prefixes: [face]
TPS: 52847     | Total: 2156382   | Errors: 0

MATCH FOUND!
Address: 0xfaceb00c123456789abcdef0123456789abcdef01
Private: 5678901234abcdef5678901234abcdef5678901234abcdef5678901234abcdef
Worker: 8
```

## Project Structure

```
addrforge-cli/
├── cmd/addrforge/     # Main entry point
│   └── main.go
├── internal/
│   ├── cli/           # CLI commands and interface
│   │   ├── root.go    # Root command setup
│   │   ├── eth.go     # Ethereum command
│   │   ├── sol.go     # Solana command
│   │   └── sui.go     # Sui command
│   ├── eth/           # Ethereum wallet generation
│   │   └── generator.go
│   ├── sol/           # Solana wallet generation  
│   │   └── generator.go
│   ├── sui/           # Sui wallet generation
│   │   └── generator.go
│   ├── util/          # Utility functions
│   │   ├── base58.go  # Base58 encoding
│   │   ├── bech32.go  # Bech32 encoding
│   │   ├── matcher.go # Prefix matching
│   │   └── output.go  # Output formatting
│   └── worker/        # Worker pool implementation
│       └── pool.go
├── go.mod
├── go.sum
└── README.md
```

## Security Notes

⚠️ **Important Security Considerations:**

1. **Private Key Handling**: Generated private keys are displayed in terminal or saved to files. Ensure:
   - Terminal history is disabled or cleared
   - Output files are properly secured
   - Files are deleted after use if not needed

2. **Random Number Generation**: Uses cryptographically secure random number generators:
   - Ethereum: `crypto/ecdsa` with secure randomness
   - Solana: `crypto/ed25519` with `crypto/rand`
   - Sui: `crypto/ed25519` with `crypto/rand`

3. **Production Use**: For production wallets:
   - Run on secure, isolated machines
   - Clear terminal history: `history -c`
   - Secure deletion of output files
   - Consider hardware wallets for high-value storage

## Contributing

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)  
5. **Open** a Pull Request

### Development Guidelines

- Follow Go conventions and `gofmt` formatting
- Add tests for new functionality
- Update documentation for new features
- Ensure code comments are in English

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- **Ethereum Foundation** for go-ethereum library
- **Solana Labs** for cryptographic standards
- **Sui Foundation** for blockchain standards
- **Cobra CLI** for excellent command-line interface framework

## Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/say8hi/addrforge-cli/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/say8hi/addrforge-cli/discussions)
- ⭐ **Star** the project if you find it useful!

---

**⚠️ Disclaimer**: This tool is for educational and legitimate purposes only. Users are responsible for complying with all applicable laws and regulations in their jurisdiction.


