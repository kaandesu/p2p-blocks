package blockchain

import (
	"bytes"
	"encoding/gob"
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
		Data:      []byte{},
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
	if err := encoder.Encode(res); err != nil {
		slog.Error("couldt not Serialize the block", "ERR", err)
	}
	return res.Bytes()
}

func (b *Block) Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		slog.Error("could not Serialize the block", "ERR", err)
	}

	return &block
}
