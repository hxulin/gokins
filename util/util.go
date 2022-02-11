package util

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"gokins/rsa"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

const (
	workspace          = "~/.gokins"
	publicKeyFilePath  = workspace + "/public.bin"
	privateKeyFilePath = workspace + "/private.bin"
	configFilePath     = workspace + "/config.yaml"
)

// FileIsExist 判断文件是否存在，存在返回true，不存在返回false
func FileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// KeyFileIsExist 判断密钥文件是否存在，存在返回true，不存在返回false
func KeyFileIsExist() (bool, error) {
	publicKeyFile, err := homedir.Expand(publicKeyFilePath)
	if err != nil {
		return false, err
	}
	privateKeyFile, err := homedir.Expand(privateKeyFilePath)
	if err != nil {
		return false, err
	}
	return FileIsExist(publicKeyFile) && FileIsExist(privateKeyFile), nil
}

// InitWorkspace 初始化工作目录
func InitWorkspace() {
	dir, err := homedir.Expand(workspace)
	if err == nil && !FileIsExist(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
}

// CreateSecretKey 创建密钥并保存到文件
func CreateSecretKey() {
	publicKeyFile, err1 := homedir.Expand(publicKeyFilePath)
	privateKeyFile, err2 := homedir.Expand(privateKeyFilePath)
	if err1 != nil || err2 != nil {
		fmt.Println("密钥文件存储位置初始化失败")
		return
	}
	// 初始化工作空间目录
	InitWorkspace()
	publicKey, privateKey, err := rsa.GenerateRsaKey(256)
	if err != nil {
		fmt.Println("密钥生成失败")
		return
	}
	err1 = ioutil.WriteFile(publicKeyFile, publicKey, os.ModePerm)
	err2 = ioutil.WriteFile(privateKeyFile, privateKey, os.ModePerm)
	if err1 != nil || err2 != nil {
		fmt.Println("密钥文件创建失败")
	} else {
		fmt.Println("密钥初始化完成，已保存到文件：")
		fmt.Println(publicKeyFile)
		fmt.Println(privateKeyFile)
	}
}

// ReadPublicKey 读取公钥内容
func ReadPublicKey() (publicKey []byte, err error) {
	keyFile, err := homedir.Expand(publicKeyFilePath)
	if err != nil {
		return nil, err
	}
	keyBytes, err := ioutil.ReadFile(keyFile)
	return keyBytes, err
}

// EncodeURIComponent URL 参数编码，实现和 JS 通用
func EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	return strings.Replace(r, "+", "%20", -1)
}
