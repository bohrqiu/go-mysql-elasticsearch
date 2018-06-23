package river

import (
	"fmt"
	"testing"
)
func Test_en(t *testing.T) {
	//http://www.seacha.com/tools/aes.html
	//算法模式：	CBC
	//密钥长度：	128
	//密钥 和 密钥偏移量填一样
	//补码方式：PKCS7Padding
	//加密结果编码方式:base64
	//配置文件配置数据为e{6leGPIac1R6Nz1Tkk8QdRg==}
	orig := "123456"
	key := "CZW4sPJTAKgcftis"
	fmt.Println("原文：", orig)

	encryptCode := AesEncrypt(orig, key)
	fmt.Println("密文：" , encryptCode)
	decryptCode := AesDecrypt(encryptCode, key)
	fmt.Println("解密结果：", decryptCode)

	decryptCode=DecryptIfPossible("e{"+encryptCode+"}",key);
	fmt.Println("解密结果：", decryptCode)

}
