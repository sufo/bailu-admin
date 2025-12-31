/**
* Create by sufo
* @Email ouamour@Gmail.com
*
* @Desc aes
 */

package aes

import (
	"bytes"
	cryptoAes "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//加密过程：
//  1、处理数据，对数据进行填充，采用PKCS7（当密钥长度不够时，缺几位补几个几）的方式。
//  2、对数据进行加密，采用AES加密方法中CBC加密模式
//  3、对得到的加密数据，进行base64加密，得到字符串
// 解密过程相反

type Aes interface {
	i()
	// Encrypt 加密
	Encrypt(encryptStr string) (string, error)

	//Decrypt 解密
	Decrypt(decryptStr string) (string, error)
}

type aes struct {
	key string
	iv  string
}

func New(key, iv string) Aes {
	return &aes{
		key: key,
		iv:  iv,
	}
}

func (a *aes) i() {}

//AES加密
func (a *aes) Encrypt(encryptStr string) (string, error) {
	encryptBytes := []byte(encryptStr)
	//创建加密实例
	block, err := cryptoAes.NewCipher([]byte(a.key))
	if err != nil {
		return "", err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes = pkcs7Padding(encryptBytes, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(a.iv))
	//blockMode := cipher.NewCBCEncrypter(block, a.key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return base64.URLEncoding.EncodeToString(crypted), nil
}

func (a *aes) Decrypt(decryptStr string) (string, error) {
	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)
	if err != nil {
		return "", err
	}

	block, err := cryptoAes.NewCipher([]byte(a.key))
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(a.iv))
	decrypted := make([]byte, len(decryptBytes))

	blockMode.CryptBlocks(decrypted, decryptBytes)
	decrypted, err = pkcs7UnPadding(decrypted)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pcks5UnPadding(decrypted []byte) []byte {
	length := len(decrypted)
	unPadding := int(decrypted[length-1])
	return decrypted[:(length - unPadding)]
}

//pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

//pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}
