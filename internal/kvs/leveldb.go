package kvs

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LevelDB struct {
	db *leveldb.DB
}

func NewLevelDB() *LevelDB {
	return &LevelDB{}
}

func (l *LevelDB) Name() string {
	return "leveldb"
}

func (l *LevelDB) Open(path string) error {
	opts := &opt.Options{
		NoSync: true,
	}
	db, err := leveldb.OpenFile(path, opts)
	if err != nil {
		return err
	}
	l.db = db
	return nil
}

func (l *LevelDB) Close() error {
	if l.db != nil {
		return l.db.Close()
	}
	return nil
}

func (l *LevelDB) Set(key string, value *Value) error {
	data, err := value.ToJSON()
	if err != nil {
		return err
	}
	return l.db.Put([]byte(key), data, nil)
}

func (l *LevelDB) Get(key string) (*Value, error) {
	data, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return ValueFromJSON(data)
}