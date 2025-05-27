package kvs

import (
	"github.com/cockroachdb/pebble"
)

type Pebble struct {
	db *pebble.DB
}

func NewPebble() *Pebble {
	return &Pebble{}
}

func (p *Pebble) Name() string {
	return "pebble"
}

func (p *Pebble) Open(path string) error {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *Pebble) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *Pebble) Set(key string, value *Value) error {
	data, err := value.ToJSON()
	if err != nil {
		return err
	}
	return p.db.Set([]byte(key), data, pebble.NoSync)
}

func (p *Pebble) Get(key string) (*Value, error) {
	data, closer, err := p.db.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	
	result := make([]byte, len(data))
	copy(result, data)
	
	return ValueFromJSON(result)
}