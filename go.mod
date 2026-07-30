module github.com/sodiqscript111/cNGN-Go

go 1.25.3

require (
	github.com/joho/godotenv v1.5.1
	// golang.org/x/crypto is required for nacl/box (Curve25519)
	// precomputation + OpenAfterPrecomputation in utils/ed25519.go.
	// The cNGN API encrypts response payloads with NaCl box using an
	// Ed25519-derived Curve25519 keypair. stdlib has no replacement.
	golang.org/x/crypto v0.54.0
)

require golang.org/x/sys v0.47.0 // indirect
