// Package crypto 提供对称加密工具，封装 AES-256-GCM 模式。
//
// 主要用途：加密用户上传的 LLM API Key。
// 密钥派生路径：SHA256(JWT_SECRET + ".model-setting-salt.v1")，
// 与 JWT 签名的密钥完全分离，但仅复用 JWT_SECRET 单源（避免引入额外配置）。
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// ModelSettingSalt 是模型配置密钥派生盐，与 JWT_SECRET 配合使用。
const ModelSettingSalt = ".model-setting-salt.v1"

// ErrCipher 加密/解密错误（对外统一包装）。
var ErrCipher = errors.New("cipher operation failed")

// AESGCM 封装 AES-256-GCM 加解密。
type AESGCM struct {
	gcm cipher.AEAD
}

// NewAESGCM 派生密钥并构造 AESGCM 实例。
// secret 是上游传入的 JWT_SECRET；不允许为空。
func NewAESGCM(secret string) (*AESGCM, error) {
	if secret == "" {
		return nil, errors.New("crypto: empty secret")
	}
	sum := sha256.Sum256([]byte(secret + ModelSettingSalt))
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &AESGCM{gcm: gcm}, nil
}

// Encrypt 用随机 nonce 加密 plain，返回 ciphertext 和 nonce。
// 每次调用都生成新的随机 nonce，nonce 大小由 gcm.NonceSize() 决定（标准 GCM 为 12 字节）。
func (a *AESGCM) Encrypt(plain []byte) (ct []byte, nonce []byte, err error) {
	nonce = make([]byte, a.gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, errors.Join(ErrCipher, err)
	}
	// Seal(dst, nonce, plaintext, additionalData) —— 不使用 additionalData
	ct = a.gcm.Seal(nil, nonce, plain, nil)
	return ct, nonce, nil
}

// Decrypt 用给定 nonce 解密 ciphertext，返回明文。
// nonce 长度必须等于 gcm.NonceSize()，否则返回 ErrCipher。
func (a *AESGCM) Decrypt(ct, nonce []byte) ([]byte, error) {
	if len(nonce) != a.gcm.NonceSize() {
		return nil, errors.Join(ErrCipher, errors.New("bad nonce size"))
	}
	plain, err := a.gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, errors.Join(ErrCipher, err)
	}
	return plain, nil
}