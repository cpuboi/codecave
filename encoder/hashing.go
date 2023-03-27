package encoder

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5HashBytesToString(inputByteArray []byte) string {
	m := md5.New()
	hash_message := hex.EncodeToString(m.Sum(inputByteArray))
	return hash_message
}

func Md5HashBytesReturnXBytes(inputByteArray []byte, returnBytes int) []byte {
	m := md5.New()
	hash_message := hex.EncodeToString(m.Sum(inputByteArray))
	return []byte(hash_message)[0:returnBytes]
}
