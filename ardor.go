package main

import (
	"./ardor-curve25519"
	"encoding/hex"
	"math/big"
	"strings"
)

var (
	prefix   = "ARODR-"
	alphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	guess    = []string{}
	syndrome = []int{0, 0, 0, 0, 0}
	codeword = []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	cwmap    = []int{3, 2, 1, 0, 7, 6, 5, 4, 13, 14, 15, 16, 12, 8, 9, 10, 11}
	gexp     = []int{1, 2, 4, 8, 16, 5, 10, 20, 13, 26, 17, 7, 14, 28, 29, 31, 27, 19, 3, 6, 12, 24, 21, 15, 30, 25, 23, 11, 22, 9, 18, 1}
	glog     = []int{0, 0, 1, 18, 2, 5, 19, 11, 3, 29, 6, 27, 20, 8, 12, 23, 4, 10, 30, 17, 7, 22, 28, 26, 21, 25, 9, 16, 13, 14, 24, 15}
)

func reset() {
	codeword = []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	cwmap = []int{3, 2, 1, 0, 7, 6, 5, 4, 13, 14, 15, 16, 12, 8, 9, 10, 11}
	gexp = []int{1, 2, 4, 8, 16, 5, 10, 20, 13, 26, 17, 7, 14, 28, 29, 31, 27, 19, 3, 6, 12, 24, 21, 15, 30, 25, 23, 11, 22, 9, 18, 1}
	glog = []int{0, 0, 1, 18, 2, 5, 19, 11, 3, 29, 6, 27, 20, 8, 12, 23, 4, 10, 30, 17, 7, 22, 28, 26, 21, 25, 9, 16, 13, 14, 24, 15}
}

func ok() bool {
	sum := 0
	for i := 1; i < 5; i++ {
		t := 0
		for j, t := 0, 0; j < 31; j++ {
			if j > 12 && j < 27 {
				continue
			}

			pos := j
			if j > 26 {
				pos -= 14
			}

			t ^= gmult(codeword[pos], gexp[(i*j)%31])
		}
		sum |= t
		syndrome[i] = t
	}
	return sum == 0

}

func addGuess() {
	s := getAccount()
	length := len(guess)
	if length > 2 {
		return
	}
	for i := 0; i < length; i++ {
		if guess[i] == s {
			return
		}
	}

	guess[length] = s
}

func gmult(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	idx := (glog[a] + glog[b]) % 31
	return gexp[idx]
}

func fromAccId(accountId string) {
	inp := []int{}
	out := []int{}
	pos := 0
	length := len(accountId)

	for i := 0; i < length; i++ {
		inp = append(inp, int(rune(accountId[i]))-48) // 48 rune('0')
	}

	for {
		divede := 0
		newLen := 0
		for i := 0; i < length; i++ {
			divede = divede*10 + inp[i]
			if divede >= 32 {
				inp[newLen] = divede >> 5
				newLen += 1
				divede &= 31
			} else if newLen > 0 {
				inp[newLen] = 0
				newLen += 1
			}
		}

		length = newLen
		out = append(out, divede)
		pos += 1

		if newLen == 0 {
			break
		}
	}

	for i := 0; i < 13; i++ {
		pos -= 1
		if pos >= 0 {
			codeword[i] = out[i]
		} else {
			codeword[i] = 0
		}
	}

	p := []int{0, 0, 0, 0}

	for i := 12; i >= 0; i-- {
		fb := codeword[i] ^ p[3]

		p[3] = p[2] ^ gmult(30, fb)
		p[2] = p[1] ^ gmult(6, fb)
		p[1] = p[0] ^ gmult(9, fb)
		p[0] = gmult(17, fb)
	}

	codeword[13] = p[0]
	codeword[14] = p[1]
	codeword[15] = p[2]
	codeword[16] = p[3]
}

func setCodeword(cw []int, length, skip int) {

	if length == 0 {
		length = 17
	}

	if skip == 0 {
		skip = -1
	}

	for i, j := 0, 0; i < length; i++ {
		if i != skip {
			codeword[cwmap[j]] = cw[i]
			j += 1
		}
	}

}

func toAccId(account string) string {
	account = account[6:]
	clean := []int{}
	length := 0
	guess = []string{}

	for i := 0; i < len(account); i++ {
		pos := strings.Index(alphabet, string(account[i]))
		if pos >= 0 {
			clean = append(clean, pos)
			length += 1
		}
	}

	if length == 17 {
		setCodeword(clean, 0, 0)
		if !ok() {
			return ""
		}
	} else {
		return ""
	}

	out := ""
	inp := make([]int, 13)
	length = 13

	for i := 0; i < 13; i++ {
		inp[i] = codeword[12-i]
	}

	for {
		divide := 0
		newLen := 0

		for i := 0; i < length; i++ {
			divide = divide*32 + inp[i]

			if divide >= 10 {
				inp[newLen] = divide / 10
				newLen += 1
				divide %= 10
			} else if newLen > 0 {
				inp[newLen] = 0
				newLen += 1
			}
		}

		length = newLen
		out += string(int(rune(divide)) + 48)

		if newLen == 0 {
			break
		}
	}

	// out 反序
	// return reverseStr(out)
	return reverseStr(out)
}

func getAccount() string {
	out := prefix
	for i := 0; i < 17; i++ {
		out += string(alphabet[codeword[cwmap[i]]])
		if i&3 == 3 && i < 13 {
			out += "-"
		}
	}
	return out
}

func accountIdToAccount(accountId string) string {
	reset()
	fromAccId(accountId)
	return getAccount()
}

func accountToAccountId(account string) string {
	reset()
	return toAccId(account)
}

func pubkeyToAccountId(pubkey string) (accountId string, err error) {
	binPubkey, err := hex.DecodeString(pubkey)
	if err != nil {
		return
	}
	hash_ := simpleHash(binPubkey)
	slice := hash_[:8]
	hexAccountId := hex.EncodeToString(reverse(slice))
	bigInt := big.NewInt(0)
	bigInt, _ = bigInt.SetString(hexAccountId, 16)
	return bigInt.String(), nil
}

func pubkeyToAccount(pubkey string) (account string, err error) {
	accId, err := pubkeyToAccountId(pubkey)
	if err != nil {
		return
	}

	return accountIdToAccount(accId), nil
}

func seedToKey(seed string) (prikey, pubkey string) {
	hash_ := simpleHash([]byte(seed))
	prikey, pubkey = ardor_curve25519.KeyGen(hex.EncodeToString(hash_))
	return
}

func AccountIdToAccount(accountId string) string {
	return accountIdToAccount(accountId)
}

func AccountToAccountId(account string) string {
	return accountToAccountId(account)
}

func PubkeyToAccountId(pubkey string) (accountId string, err error) {
	return pubkeyToAccountId(pubkey)
}

func PubkeyToAccount(pubkey string) (account string, err error) {
	return pubkeyToAccount(pubkey)
}

func SeedToKey(seed string) (prikey, pubkey string) {
	return seedToKey(seed)
}
