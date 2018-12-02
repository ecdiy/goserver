package pay

import (
	"net/url"
	"sort"
	"fmt"
	"encoding/base64"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/rand"
	"crypto"
)

func aliSign(m url.Values, privateKey string) string {
	//对url.values进行排序
	sign := ""
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if m.Get(k) != "" {
			if i == 0 {
				sign = k + "=" + m.Get(k)
			} else {
				sign = sign + "&" + k + "=" + m.Get(k)
			}
		}
	}
	fmt.Println(sign)
	//对排序后的数据进行rsa2加密，获得sign
	b, ree := rsaEncrypt([]byte(sign), privateKey)
	if ree != nil {
		fmt.Println("加密失败：", ree)
		return ""
	} else {
		res := base64.StdEncoding.EncodeToString(b)
		fmt.Println("base加密：", res)
		return res
	}
}

func rsaEncrypt(origData []byte, private_key string) ([]byte, error) {
	block2, _ := pem.Decode([]byte(private_key)) //PiravteKeyData为私钥文件的字节数组
	if block2 == nil {
		fmt.Println("block空")
		return nil, nil
	}
	//priv即私钥对象,block2.Bytes是私钥的字节流
	priv, err := x509.ParsePKCS8PrivateKey(block2.Bytes)
	if err != nil {
		fmt.Println("无法还原私钥")
		return nil, nil
	}
	p := priv.(*rsa.PrivateKey)
	h2 := sha256.New()
	h2.Write(origData)
	hashed := h2.Sum(nil)
	signature2, err := rsa.SignPKCS1v15(rand.Reader, p, crypto.SHA256, hashed) //签名
	return signature2, err
}
