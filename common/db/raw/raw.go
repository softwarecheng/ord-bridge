package raw

import (
	badger "github.com/dgraph-io/badger/v4"
)

func Set(key, value []byte, db *badger.DB) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func BatchSet(key []byte, value []byte, wb *badger.WriteBatch) error {
	return wb.Set([]byte(key), value)
}

func Get(key []byte, db *badger.DB) ([]byte, error) {
	var ret []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		ret, err = item.ValueCopy(nil)
		return err
	})

	return ret, err
}

func TxnGet(key []byte, txn *badger.Txn) ([]byte, error) {
	item, err := txn.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

func TxnGetList(prefix []byte, txn *badger.Txn) (map[string][]byte, error) {
	result := make(map[string][]byte)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	for it.Seek(prefix); it.ValidForPrefix([]byte(prefix)); it.Next() {
		item := it.Item()
		value, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}
		key := string(item.KeyCopy(nil))
		result[key] = value
	}
	return result, nil
}
