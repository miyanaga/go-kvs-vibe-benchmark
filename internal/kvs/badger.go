package kvs

import (
	"github.com/dgraph-io/badger/v4"
)

type Badger struct {
	db *badger.DB
}

func NewBadger() *Badger {
	return &Badger{}
}

func (b *Badger) Name() string {
	return "badger"
}

func (b *Badger) Open(path string) error {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	b.db = db
	return nil
}

func (b *Badger) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

func (b *Badger) Set(key string, value *Value) error {
	data, err := value.ToJSON()
	if err != nil {
		return err
	}
	
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

func (b *Badger) Get(key string) (*Value, error) {
	var data []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			data = append([]byte{}, val...)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	
	return ValueFromJSON(data)
}