package ardor_curve25519

/*
#include "lib/curve25519.c"
*/
import "C"
import (
	"crypto/sha256"
	"encoding/hex"
)

func keyGen(secret string) (prikey, pubkey string) {
	prikey = C.GoString(C.GetSignPrikey(C.CString(secret)))
	pubkey = C.GoString(C.GetSignPubkey(C.CString(secret)))
	return
}

func clamp(k []byte) {
	k[31] &= 0x7F
	k[31] |= 0x40
	k[0] &= 0xF8
}

func Sign(prikey string, message []byte) []byte {
	// prikey secret sha256后的
	s, _ := hex.DecodeString(prikey)
	m := sha256hash(message)
	x := sha256hash(m, s)
	clamp(x)
	_, pub := keyGen(hex.EncodeToString(x)) // 协商公钥
	y, _ := hex.DecodeString(pub)
	h := sha256hash(m, y)
	ch := C.CString(hex.EncodeToString(h))
	cx := C.CString(hex.EncodeToString(x))
	cs := C.CString(hex.EncodeToString(s))
	v := C.GoString(C.Sign(ch, cx, cs))

	vs, _ := hex.DecodeString(v)
	return append(vs, h...)
}


func sha256hash(msg ...[]byte) []byte {
	hasher := sha256.New()
	for i := range msg {
		hasher.Write(msg[i])
	}
	return hasher.Sum(nil)
}

func KeyGen(secret string) (prikey, pubkey string) {
	return keyGen(secret)
}
