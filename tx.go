package main

import (
	"./ardor-curve25519"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

func serialize(pubkey, value, fee, blockHeight, blockId, recipient []byte) []byte {
	tx := bytes.NewBuffer(nil)
	tx.WriteByte(1)            // chain
	tx.Write(pad(3, 0))        // 共4字节
	tx.WriteByte(254)          // type 会 254-256 网页上显示-2
	tx.WriteByte(0)            // subtype
	tx.WriteByte(1)            // version
	tx.Write(nxtTimeBytes())   // timestamp 4字节
	tx.WriteByte(15)           // deadline
	tx.WriteByte(0)            // 空一个 暂时不知道干啥用 好像都是空
	tx.Write(pubkey)           // 发送者公钥
	tx.Write(recipient)        // 收款人id 反序 如果长度9取前8 如果不足9 往后补空字节 8字节
	tx.Write(value)            // 转账金额 * 100000000 反序  如果长度9取前8 如果不足8 往后补空字节
	tx.Write(fee)              // 矿工费 长度不足8补足8
	tx.Write(make([]byte, 32)) // ref full hash 好像没啥用
	tx.Write(make([]byte, 32)) // signature
	tx.Write(blockHeight)
	tx.Write(blockId)
	tx.Write(pad(4, 0))
	return tx.Bytes()
}

func nxtTimeBytes() []byte {
	return reverse(uint32ToBytes(uint32(time.Now().Unix()) - 1514296800))
}

func uint32ToBytes(i uint32) []byte {
	temp := make([]byte, 4)
	binary.BigEndian.PutUint32(temp, i)
	return temp
}

func uint64ToBytes(i uint64) []byte {
	temp := make([]byte, 8)
	binary.BigEndian.PutUint64(temp, i)
	return temp
}

func MakeTx(recipientAccount string, amount int, pubkey string, blockHeight, blockId int) []byte {
	recipientAccId := accountToAccountId(recipientAccount)
	bigInt := big.NewInt(0)
	bigInt.SetString(recipientAccId, 10)
	recipient := bigInt.Bytes()
	recipient = reverse(recipient)
	for {
		if len(recipient) < 8 {
			recipient = append(recipient, byte(0))
		} else if len(recipient) > 8 {
			recipient = recipient[:8]
		} else {
			break
		}
	}

	amt := reverse(uint64ToBytes(uint64(amount)))
	fee := uint64ToBytes(uint64(100000000))
	fee = reverse(fee)
	pub, _ := hex.DecodeString(pubkey)

	height := reverse(uint32ToBytes(uint32(blockHeight)))
	blkId := reverse(uint64ToBytes(uint64(blockId)))

	return serialize(pub, amt, fee, height, blkId, recipient)
}

func SignTx(prikey string, unsignedTx []byte) []byte {
	sig := ardor_curve25519.Sign(prikey, unsignedTx)
	signedTx := append(unsignedTx[:69], sig...)
	if cap(signedTx) > 149 {
		return signedTx[:149]
	}
	return append(signedTx, make([]byte, 149)...)[:149]
}


func pad(length, val int) []byte {
	t := make([]byte, length)
	for i := 0; i < length; i++ {
		t[i] = uint8(val)
	}
	return t
}
