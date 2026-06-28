package auth

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const (
	maxTimestampDiffSeconds = 15 * 60
	randomTokenLength       = 22 // ~2^132 keyspace
)

var (
	errInvalidToken     = errors.New("invalid token")
	errCannotParseToken = errors.New("cannot parse token")
	errNoPrivateKey     = errors.New("no private key")
	tokenCharacters     = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
	pluginPrefix        = "P"
	enhancedTokenPrefix = "gtfy"

	randReader = rand.Reader
)

type EnhancedToken struct {
	ident        string
	pubOrPrivKey []byte
	timestamp    int64
	signature    []byte
}

// PublicForm returns the a canonicalized representation of the public key.
func (c *EnhancedToken) PublicForm() string {
	if c.timestamp != 0 || len(c.signature) != 0 {
		return enhancedTokenPrefix + c.ident + "." + base64.RawURLEncoding.EncodeToString(c.pubOrPrivKey)
	}
	privKey := ed25519.NewKeyFromSeed(c.pubOrPrivKey)
	return enhancedTokenPrefix + c.ident + "." + base64.RawURLEncoding.EncodeToString(privKey.Public().(ed25519.PublicKey))
}

// Sign signs the timestamp with the private key and returns a new EnhancedToken.
func (c *EnhancedToken) Sign(timestamp int64) (*EnhancedToken, error) {
	if c.timestamp != 0 || len(c.signature) != 0 || len(c.pubOrPrivKey) != ed25519.SeedSize {
		return nil, errNoPrivateKey
	}
	privKey := ed25519.NewKeyFromSeed(c.pubOrPrivKey)
	sha512 := sha512.New()
	sha512.Write([]byte("iat="))
	fmt.Fprintf(sha512, "%d", timestamp)
	sign, err := privKey.Sign(nil, sha512.Sum(nil), crypto.SHA512)
	if err != nil {
		return nil, err
	}
	return &EnhancedToken{
		ident:        c.ident,
		pubOrPrivKey: privKey.Public().(ed25519.PublicKey),
		timestamp:    timestamp,
		signature:    sign,
	}, nil
}

func (c *EnhancedToken) ValidateTimestamp(now int64) bool {
	if c.timestamp == 0 && len(c.signature) == 0 {
		return true
	}
	if c.timestamp < now-maxTimestampDiffSeconds {
		return false
	}
	if c.timestamp > now+maxTimestampDiffSeconds {
		return false
	}
	return true
}

// String marshals the token into a string.
func (c *EnhancedToken) String() string {
	var b strings.Builder
	b.WriteString(enhancedTokenPrefix)
	b.WriteString(c.ident)
	b.WriteByte('.')
	b.WriteString(base64.RawURLEncoding.EncodeToString(c.pubOrPrivKey))
	if c.timestamp != 0 || len(c.signature) != 0 {
		fmt.Fprintf(&b, ".%d.", c.timestamp)
		b.WriteString(base64.RawURLEncoding.EncodeToString(c.signature))
	}
	return b.String()
}

// NewEnhancedToken creates a new EnhancedToken.
func NewEnhancedToken(ident string) *EnhancedToken {
	ident = strings.ReplaceAll(ident, ".", "_")
	var seed [ed25519.SeedSize]byte
	_, err := rand.Read(seed[:])
	if err != nil {
		panic("unreachable: random source should never return an error")
	}
	return &EnhancedToken{ident: ident, pubOrPrivKey: seed[:]}
}

// ParseEnhancedToken parses a string into an EnhancedToken.
func ParseEnhancedToken(token string) (*EnhancedToken, error) {
	token, found := strings.CutPrefix(token, enhancedTokenPrefix)
	if !found {
		return nil, fmt.Errorf("%w: token must start with %s", errCannotParseToken, enhancedTokenPrefix)
	}

	// count number of dots, one dot -> ident then private key, three dots -> ident, public key, challenge then signature
	fields := strings.SplitN(token, ".", 4)
	if len(fields) != 2 && len(fields) != 4 {
		return nil, fmt.Errorf("%w: token must have 2 or 4 fields separated by dots", errCannotParseToken)
	}
	ident := fields[0]
	pkOrPubkeyB64 := fields[1]
	pkOrPubkeyBytesLen := base64.RawURLEncoding.DecodedLen(len(pkOrPubkeyB64))
	pkOrPubkey, err := base64.RawURLEncoding.DecodeString(pkOrPubkeyB64)
	if err != nil {
		return nil, fmt.Errorf("%w: base64 decode failed: %w", errCannotParseToken, err)
	}
	if len(fields) == 2 {
		if pkOrPubkeyBytesLen != ed25519.SeedSize {
			return nil, fmt.Errorf("%w: private key must be %d bytes", errCannotParseToken, ed25519.SeedSize)
		}
		return &EnhancedToken{
			ident:        ident,
			pubOrPrivKey: pkOrPubkey,
		}, nil
	}
	if pkOrPubkeyBytesLen != ed25519.PublicKeySize {
		return nil, fmt.Errorf("%w: public key must be %d bytes", errCannotParseToken, ed25519.PublicKeySize)
	}
	timestampStr := fields[2]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: timestamp must be an integer: %w", errCannotParseToken, err)
	}
	signatureB64 := fields[3]
	signatureBytesLen := base64.RawURLEncoding.DecodedLen(len(signatureB64))
	if signatureBytesLen != ed25519.SignatureSize {
		return nil, fmt.Errorf("%w: signature must be %d bytes", errCannotParseToken, ed25519.SignatureSize)
	}
	signature, err := base64.RawURLEncoding.DecodeString(signatureB64)
	if err != nil {
		return nil, fmt.Errorf("%w: base64 decode failed: %w", errCannotParseToken, err)
	}
	sha512 := sha512.New()
	sha512.Write([]byte("iat=")) // query-like encoding to give us some semantic headroom should we need more fields in the future
	fmt.Fprintf(sha512, "%d", timestamp)
	if err := ed25519.VerifyWithOptions(pkOrPubkey, sha512.Sum(nil), signature, &ed25519.Options{Hash: crypto.SHA512}); err != nil {
		return nil, errInvalidToken
	}
	return &EnhancedToken{
		ident:        ident,
		pubOrPrivKey: pkOrPubkey,
		timestamp:    timestamp,
		signature:    signature,
	}, nil
}

func randIntn(n int) int {
	max := big.NewInt(int64(n))
	res, err := rand.Int(randReader, max)
	if err != nil {
		panic("random source is not available")
	}
	return int(res.Int64())
}

// GenerateApplicationToken generates an application token.
func GenerateApplicationToken() (publicForm, privateForm string) {
	token := NewEnhancedToken("a")
	return token.PublicForm(), token.String()
}

// GenerateClientToken generates a client token.
func GenerateClientToken() (publicForm, privateForm string) {
	token := NewEnhancedToken("c")
	return token.PublicForm(), token.String()
}

// GeneratePluginToken generates a plugin token.
func GeneratePluginToken() string {
	return pluginPrefix + generateRandomString(randomTokenLength)
}

// GenerateImageName generates an image name.
func GenerateImageName() string {
	return generateRandomString(25)
}

func generateRandomString(length int) string {
	res := make([]byte, length)
	for i := range res {
		index := randIntn(len(tokenCharacters))
		res[i] = tokenCharacters[index]
	}
	return string(res)
}

func init() {
	randIntn(2)
}
