package basic

import (
	badger "github.com/dgraph-io/badger/v4"
)

func Delete(key []byte, db *badger.DB) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func Set(key []byte, value []byte, db *badger.DB) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}
