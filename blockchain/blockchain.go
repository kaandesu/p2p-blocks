package blockchain

import (
	"log/slog"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./chainDB"
)

type BlockChain struct {
	Database *badger.DB
	LastHash []byte
}

type BlockChainiterator struct {
	Database    *badger.DB
	CurrentHash []byte
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	if err != nil {
		slog.Error("could not open the badger db", "err", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			slog.Info("no existing blockchain found")
			genesis := Genesis()
			slog.Info("Genesys block proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			lastHash = genesis.Hash
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				slog.Error("getting lh failed", "err", err)
			}
			err = item.Value(func(val []byte) error {
				lastHash = append(lastHash, val...)
				return nil
			})
			return err
		}
	})
	if err != nil {
		slog.Error("could not update the database", "error", err)
	}

	blockchain := BlockChain{LastHash: lastHash, Database: db}

	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			slog.Error("can't get lh item", "err", err)
		}
		err = item.Value(func(val []byte) error {
			lastHash = append(lastHash, val...)
			return nil
		})

		return err
	})
	if err != nil {
		slog.Error("[AddBlock] cant view databse", "err", err)
	}

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte("lh"), newBlock.Serialize())
		if err != nil {
			slog.Error("[AddBlock] can't Serialize", "err", err)
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)
		return err
	})
	if err != nil {
		slog.Error("[AddBlock] could not update the database", "err", err)
	}
}

func (chain *BlockChain) Iterator() *BlockChainiterator {
	return &BlockChainiterator{
		CurrentHash: chain.LastHash,
		Database:    chain.Database,
	}
}

func (iter *BlockChainiterator) Next() *Block {
	var block *Block
	var encodedBlock []byte
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			slog.Error("[Next] could not get the current item", "err", err)
		}
		err = item.Value(func(val []byte) error {
			encodedBlock = append(encodedBlock, val...)
			return nil
		})
		block = Deserialize(encodedBlock)
		return err
	})
	if err != nil {
		slog.Error("[Next] Could not view the database", "err", err)
	}

	iter.CurrentHash = block.PrevHash

	return block
}
