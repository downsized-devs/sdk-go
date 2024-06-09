package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/downsized-devs/sdk-go/codes"
	"github.com/downsized-devs/sdk-go/errors"
	"github.com/downsized-devs/sdk-go/logger"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

const (
	keyLen = 32 // length of the derived key
)

type Interface interface {
	Encrypt(ctx context.Context, passphrase string, timestamp int64, plaintext string) string
	Decrypt(ctx context.Context, passphrase string, timestamp int64, ciphertext string) (string, error)
	ScryptPassword(ctx context.Context, salt, password string) string
	CompareScryptPassword(ctx context.Context, passwordHash, salt, password string) bool
	HashPassword(ctx context.Context, secretKey, password string) string
	CompareHashPassword(ctx context.Context, secretKey, hashPassword, password string) bool
}

type security struct {
	log    logger.Interface
	scrypt ScryptConfig
}

type ScryptConfig struct {
	Base64SignerKey     string
	Base64SaltSeperator string
	Rounds              int
	MemoryCost          int
}

func Init(log logger.Interface, scrypt ScryptConfig) Interface {
	return &security{
		log:    log,
		scrypt: scrypt,
	}
}

func (s *security) Encrypt(ctx context.Context, passphrase string, timestamp int64, plaintext string) string {
	passphrase = fmt.Sprintf("%s%d", passphrase, timestamp)

	// get key by reverse salt with passphrase
	key, salt := s.deriveKey(ctx, passphrase, nil)

	// iv is nonce, number once use, ack as 1 time use key
	// this also random auto generate
	// unique for all time per given key
	// https://golang.org/pkg/crypto/cipher/
	iv := make([]byte, 12)

	// feed the iv into rand seed
	_, err := rand.Read(iv)
	if err != nil {
		s.log.Error(ctx, err)
	}

	// create the encryption
	b, err := aes.NewCipher(key)
	if err != nil {
		s.log.Error(ctx, err)
	}

	aesgcm, err := cipher.NewGCM(b)
	if err != nil {
		s.log.Error(ctx, err)
	}

	data := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return fmt.Sprintf("%s-%s-%s", hex.EncodeToString(salt), hex.EncodeToString(iv), hex.EncodeToString(data))
}

func (s *security) Decrypt(ctx context.Context, passphrase string, timestamp int64, ciphertext string) (string, error) {
	passphrase = fmt.Sprintf("%s%d", passphrase, timestamp)
	// split character by -
	arr := strings.Split(ciphertext, "-")
	if len(arr) < 3 {
		return "", errors.NewWithCode(codes.CodeErrorSecurityInvalidChipper, "chipper is empty string or invalid")
	}

	// salt is the first word when split by -
	salt, err := hex.DecodeString(arr[0])
	if err != nil {
		return "", err
	}

	// iv is nonce, number once use, ack as 1 time use key
	// this also random auto generate
	// unique for all time per given key
	// https://golang.org/pkg/crypto/cipher/
	iv, err := hex.DecodeString(arr[1])
	if err != nil {
		return "", err
	}

	// data is result of encrypted
	data, err := hex.DecodeString(arr[2])
	if err != nil {
		return "", err
	}

	// get key by reverse salt with passphrase
	key, _ := s.deriveKey(ctx, passphrase, salt)
	b, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(b)
	if err != nil {
		return "", err
	}

	data, err = aesgcm.Open(nil, iv, data, nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *security) deriveKey(ctx context.Context, passphrase string, salt []byte) ([]byte, []byte) {
	// if salt is null
	// random generate salt by 8 character
	if salt == nil {
		salt = make([]byte, 8)
		// http://www.ietf.org/rfc/rfc2898.txt
		// Salt.
		_, err := rand.Read(salt)
		if err != nil {
			s.log.Error(ctx, err)
		}
	}
	// make aes256 key by passphrase and salt, and return salt
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

func (s *security) HashPassword(ctx context.Context, secretKey, password string) string {
	computedHash := hmac.New(sha256.New, []byte(secretKey))
	computedHash.Write([]byte(password))
	return hex.EncodeToString(computedHash.Sum(nil))
}

func (s *security) CompareHashPassword(ctx context.Context, hashPassword, secretKey, password string) bool {
	return hashPassword == s.HashPassword(ctx, secretKey, password)
}

func (s *security) ScryptPassword(ctx context.Context, salt string, password string) string {
	// Derive the key using SCRYPT
	N := 1 << s.scrypt.MemoryCost // CPU/memory cost parameter (must be a power of 2)
	r := s.scrypt.Rounds          // block size parameter
	p := 1                        // parallelization parameter
	saltConcat := append(s.b64Stddecode(ctx, salt), s.b64Stddecode(ctx, s.scrypt.Base64SaltSeperator)...)
	derivedKey, err := scrypt.Key([]byte(password), saltConcat, N, r, p, keyLen)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("Error generating derived key: %v", err))
	}

	// Encrypt the derived key using AES-CTR and return encoded string
	return base64.StdEncoding.EncodeToString(s.encryptDerivedKey(ctx, derivedKey, s.b64Stddecode(ctx, s.scrypt.Base64SignerKey)))
}

func (s *security) encryptDerivedKey(ctx context.Context, derivedKey, plaintext []byte) []byte {
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("Error creating AES cipher: %v", err))
	}

	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext
}

func (s *security) CompareScryptPassword(ctx context.Context, passwordHash, salt, password string) bool {
	return s.areBytesEqual(s.b64Stddecode(ctx, s.ScryptPassword(ctx, salt, password)), s.b64Stddecode(ctx, passwordHash))
}

func (s *security) b64Stddecode(ctx context.Context, str string) []byte {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("Failed to decode string: %v", err))
	}

	return b
}

func (s *security) areBytesEqual(a, b []byte) bool {
	return md5.Sum(a) == md5.Sum(b)
}
