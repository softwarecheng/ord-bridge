package proto3

import (
	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Set(key []byte, data protoreflect.ProtoMessage, db *badger.DB) error {
	dataBytes, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return db.Update(func(txn *badger.Txn) error {
		err = txn.Set(key, dataBytes)
		if err != nil {
			return err
		}
		return nil
	})
}

func BatchSet(key []byte, data protoreflect.ProtoMessage, wb *badger.WriteBatch) error {
	dataBytes, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	err = wb.Set([]byte(key), dataBytes)
	if err != nil {
		return err
	}
	return nil
}

func Get(key []byte, protoMsg protoreflect.ProtoMessage, db *badger.DB) error {
	var err error
	return db.View(func(txn *badger.Txn) error {
		err = TxnGet(key, protoMsg, txn)
		if err != nil {
			return err
		}
		return nil
	})
}

func TxnGet(key []byte, protoMsg protoreflect.ProtoMessage, txn *badger.Txn) error {
	item, err := txn.Get([]byte(key))
	if err != nil {
		return err
	}
	return item.Value(func(v []byte) error {
		return proto.Unmarshal(v, protoMsg)
	})
}

func TxnGetWithType[T protoreflect.ProtoMessage](key []byte, txn *badger.Txn) (*T, error) {
	var ret T
	item, err := txn.Get([]byte(key))
	if err != nil {
		return &ret, err
	}
	err = item.Value(func(v []byte) error {
		return proto.Unmarshal(v, ret)
	})
	if err != nil {
		return nil, err
	}
	return &ret, err
}

func TxnGetList[T protoreflect.ProtoMessage](prefix []byte, txn *badger.Txn) (map[string]*T, error) {
	result := make(map[string]*T)
	itr := txn.NewIterator(badger.DefaultIteratorOptions)
	defer itr.Close()

	for itr.Seek([]byte(prefix)); itr.ValidForPrefix([]byte(prefix)); itr.Next() {
		item := itr.Item()
		var value T
		err := item.Value(func(v []byte) error {
			return proto.Unmarshal(v, value)
		})
		if err != nil {
			return nil, err
		}
		key := string(item.KeyCopy(nil))
		result[key] = &value
	}
	return result, nil
}
