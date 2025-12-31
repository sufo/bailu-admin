/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

// 公钥加密
// path 公钥文件路径
func PublicEncryptByFile(encryptStr string, path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 读取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	// pem 解码
	block, _ := pem.Decode(buf)

	return rsaEncrypt(encryptStr, block.Bytes)
}

//func PublicEncrypt(encryptStr string, path string) (string, error) {
//	// 打开文件
//	file, err := os.Open(path)
//	if err != nil {
//		return "", err
//	}
//	defer file.Close()
//
//	// 读取文件内容
//	info, _ := file.Stat()
//	buf := make([]byte, info.Size())
//	file.Read(buf)
//
//	// pem 解码
//	block, _ := pem.Decode(buf)
//
//	// x509 解码
//	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		return "", err
//	}
//
//	// 类型断言
//	publicKey := publicKeyInterface.(*rsa.PublicKey)
//
//	//对明文进行加密
//	encryptedStr, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(encryptStr))
//	if err != nil {
//		return "", err
//	}
//
//	//返回密文
//	return base64.URLEncoding.EncodeToString(encryptedStr), nil
//
//}

func PublicEncrypt(encryptStr string, base64PubKey string) (base64CipherText string, err error) {
	pub, err := base64.StdEncoding.DecodeString(base64PubKey)
	if err != nil {
		return "", err
	}
	return rsaEncrypt(encryptStr, pub)
}

func rsaEncrypt(encryptStr string, publicKey []byte) (base64CipherText string, err error) {
	// x509 解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	// 类型断言
	rsaPublicKey := publicKeyInterface.(*rsa.PublicKey)

	//对明文进行加密
	encryptedStr, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(encryptStr))
	if err != nil {
		return "", err
	}

	//返回密文
	return base64.StdEncoding.EncodeToString(encryptedStr), nil
}

//// 私钥解密
//func PrivateDecrypt(decryptStr string, path string) (string, error) {
//	// 打开文件
//	file, err := os.Open(path)
//	if err != nil {
//		return "", err
//	}
//	defer file.Close()
//
//	// 获取文件内容
//	info, _ := file.Stat()
//	buf := make([]byte, info.Size())
//	file.Read(buf)
//
//	// pem 解码
//	block, _ := pem.Decode(buf)
//
//	// X509 解码
//	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		return "", err
//	}
//	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)
//
//	//对密文进行解密
//	decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decryptBytes)
//
//	//返回明文
//	return string(decrypted), nil
//}

// 私钥解密
// path 私钥文件路径
func PrivateDecryptByFile(decryptStr string, path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)

	// pem 解码
	block, _ := pem.Decode(buf)

	return rsaDecrypt(decryptStr, block.Bytes)
}

// 私钥解密
// base64PriKey
func PrivateDecrypt(base64CipherText string, base64PriKey string) (string, error) {
	privateBytes, err := base64.StdEncoding.DecodeString(base64PriKey)
	if err != nil {
		return "", err
	}

	return rsaDecrypt(base64CipherText, privateBytes)
}

func rsaDecrypt(decryptStr string, PrivKey []byte) (string, error) {
	// X509 解码
	privateKey, err := x509.ParsePKCS1PrivateKey(PrivKey)
	if err != nil {
		return "", err
	}
	decryptBytes, err := base64.StdEncoding.DecodeString(decryptStr)

	//对密文进行解密
	decrypted, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decryptBytes)

	//返回明文
	return string(decrypted), nil
}
