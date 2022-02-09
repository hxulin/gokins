package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

// GenerateRsaKey 生成RSA私钥和公钥
// bits 密钥长度
func GenerateRsaKey(bits int) ([]byte, []byte, error) {
	// GenerateKey 函数使用随机数据生成器 random 生成一对具有指定字位数的 RSA 密钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	// 通过 x509 标准序列化 RSA 私钥
	x509privateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	// 获取公钥的数据
	publicKey := privateKey.PublicKey
	// 通过 x509 标准序列化公钥信息
	x509publicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, nil, err
	}
	return x509publicKey, x509privateKey, nil
}

// Encrypt RSA 加密
// publicKey 公钥
// plainText 需要加密的数据
func Encrypt(publicKey []byte, plainText string) (string, error) {
	// x509 解码公钥
	x509publicKey, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	// 类型断言
	key := x509publicKey.(*rsa.PublicKey)
	// 对明文进行加密
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(plainText))
	if err != nil {
		return "", err
	}
	// 返回密文
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

// Decrypt RSA 解密
// cipherText 需要解密的数据
// privateKey 私钥
func Decrypt(privateKey []byte, cipherText string) (string, error) {
	encryptedData, err := base64.URLEncoding.DecodeString(cipherText)
	if err != err {
		return "", err
	}
	// x509 解码私钥
	x509privateKey, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	// 对密文进行解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, x509privateKey, encryptedData)
	if err != nil {
		return "", err
	}
	// 返回明文
	return string(plainText), nil
}
