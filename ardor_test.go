package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	rightSeed      = "society tumble goose deep dumb shook candle spirit gay aim muscle boat"
	rightAccount   = "ARDOR-HFP3-LUHG-TZ4D-HZHP6"
	rightAccountId = "17462144152575129249"
	rightPrikey    = "72919baa138b0411e61f9fa4108ad0c167159199eef1e6bac8538c012aeee303"
	rightPubkey    = "fc818b2f7f29df6573af46537cd7906b59499dadc4da5309cd6683cfe6dc333a"
)

func TestSeedToKey(t *testing.T) {
	pri, pub := SeedToKey(rightSeed)
	fmt.Println(pri, pub)
	// 72919baa138b0411e61f9fa4108ad0c167159199eef1e6bac8538c012aeee303 fc818b2f7f29df6573af46537cd7906b59499dadc4da5309cd6683cfe6dc333a
}

func TestAccountIdToAccount(t *testing.T) {
	fmt.Println(AccountIdToAccount(rightAccountId))
}

func TestAccountToAccountId(t *testing.T) {
	fmt.Println(AccountToAccountId(rightAccount))
}

func TestPubkeyToAccount(t *testing.T) {
	fmt.Println(PubkeyToAccount(rightPubkey))
}

func TestPubkeyToAccountId(t *testing.T) {
	fmt.Println(PubkeyToAccountId(rightPubkey))
}

func TestMakeTx(t *testing.T) {
	unsigned := MakeTx("ARDOR-RH8M-M566-684A-3A8ME", 200000000, "fc818b2f7f29df6573af46537cd7906b59499dadc4da5309cd6683cfe6dc333a", 2559714, 5327494388856501683)
	fmt.Println(hex.EncodeToString(unsigned))
	signed := SignTx("72919baa138b0411e61f9fa4108ad0c167159199eef1e6bac8538c012aeee303", unsigned)
	fmt.Println(hex.EncodeToString(signed))
}
