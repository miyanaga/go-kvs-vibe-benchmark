package kvs

import (
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

type BBolt struct {
	db *bolt.DB
}

var bucketName = []byte("kvs")

func NewBBolt() *BBolt {
	return &BBolt{}
}

func (b *BBolt) Name() string {
	return "bbolt"
}

func (b *BBolt) Open(path string) error {
	dbPath := filepath.Join(path, "bbolt.db")
	opts := bolt.DefaultOptions
	opts.NoSync = true
	db, err := bolt.Open(dbPath, 0600, opts)
	if err != nil {
		return err
	}
	
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		db.Close()
		return err
	}
	
	b.db = db
	return nil
}

func (b *BBolt) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

func (b *BBolt) Set(key string, value *Value) error {
	data, err := value.ToJSON()
	if err != nil {
		return err
	}
	
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		return bucket.Put([]byte(key), data)
	})
}

func (b *BBolt) Get(key string) (*Value, error) {
	var data []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		data = bucket.Get([]byte(key))
		if data == nil {
			return bolt.ErrInvalid
		}
		tmp := make([]byte, len(data))
		copy(tmp, data)
		data = tmp
		return nil
	})
	if err != nil {
		return nil, err
	}
	
	return ValueFromJSON(data)
}