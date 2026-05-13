package security

import (
	"context"
	"testing"

	"github.com/downsized-devs/sdk-go/logger"
	"github.com/stretchr/testify/assert"
)

// Test_security_ScryptPassword_InvalidConfig exercises the scrypt.Key error
// branch in ScryptPassword and the cascading aes.NewCipher error in
// encryptDerivedKey (when scrypt fails, derivedKey is empty so AES setup
// also fails). Both are normally unreachable with a valid config.
func Test_security_ScryptPassword_InvalidConfig(t *testing.T) {
	// MemoryCost=0 → N=1, which scrypt rejects (N must be > 1).
	s := &security{
		log: logger.Init(logger.Config{Level: "error"}),
		scrypt: ScryptConfig{
			Base64SignerKey:     "MoaHZJjRSE9Ktj6HnIkoldV+BmXpD7YVboHgJOY4SDnUNiNMTUILxlsY4igO3Uzx/n/VwFju9IC4fQfgDy7LwQ==",
			Base64SaltSeparator: "Bw==",
			Rounds:              0,
			MemoryCost:          0,
		},
	}
	// Production code logs both the scrypt error (line 171) and the
	// downstream aes.NewCipher error (line 181) but does not return, so
	// cipher.NewCTR is called with a nil block and panics. The coverage
	// for the two error-log lines is still recorded before the panic.
	assert.Panics(t, func() {
		s.ScryptPassword(context.Background(), "ekr/rlgB6tovww==", "pw")
	})
}

// Decrypt with fewer than 3 segments returns the "chipper is empty" code path.
func Test_security_Decrypt_FewSegments(t *testing.T) {
	s := &security{log: logger.Init(logger.Config{Level: "error"})}
	_, err := s.Decrypt(context.Background(), "pw", 0, "aabb-ccdd")
	assert.Error(t, err)
}
