package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"log/slog"
	"time"
)

type Block struct {
	Hash      []byte
	Data      []byte
	PrevHash  []byte
	Nonce     int
	Timestamp int64
}

func CreateBlock(data string, prevHash []byte) (b *Block) {
	b = &Block{
		Hash:      []byte{},
		Data:      []byte(data),
		PrevHash:  prevHash,
		Nonce:     0,
		Timestamp: time.Now().Unix(),
	}

	pow := NewProof(b)
	b.Nonce, b.Hash = pow.Run()
	return
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	if err := encoder.Encode(b); err != nil {
		log.Panic(err)
	}
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		slog.Error("could not Serialize the block", "ERR", err)
	}

	return &block
}
