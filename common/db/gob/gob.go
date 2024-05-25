package gob

import (
	"bytes"
	"encoding/gob"

	badger "github.com/dgraph-io/badger/v4"
)

func Set(key []byte, value interface{}, db *badger.DB) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return err
	}
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, buf.Bytes())
	})
}

func BatchSet(key []byte, data interface{}, wb *badger.WriteBatch) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return wb.Set([]byte(key), buf.Bytes())
}

func TxnGet(key []byte, txn *badger.Txn, target interface{}) error {
	item, err := txn.Get([]byte(key))
	if err != nil {
		return err
	}
	return item.Value(func(v []byte) error {
		err := decode(v, target)
		if err != nil {
			return err
		}
		return nil
	})
}

func TxnGetWithType[T any](key []byte, txn *badger.Txn) (*T, error) {
	item, err := txn.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	var ret T
	err = item.Value(func(v []byte) error {
		return decode(v, &ret)
	})
	return &ret, err
}

func GetList[T any](prefix []byte, txn *badger.Txn) (map[string]T, error) {
	result := make(map[string]T)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
		item := it.Item()
		var value T
		err := item.Value(func(data []byte) error {
			return decode(data, &value)
		})
		if err != nil {
			return nil, err
		}
		key := string(item.KeyCopy(nil))
		result[key] = value
	}
	return result, nil
}

func decode(data []byte, target interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(target)
}
