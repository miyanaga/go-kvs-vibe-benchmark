package kvs

import "encoding/json"

type Value struct {
	Single *int `json:"single,omitempty"`
	Double *int `json:"double,omitempty"`
}

func (v *Value) ToJSON() ([]byte, error) {
	return json.Marshal(v)
}

func ValueFromJSON(data []byte) (*Value, error) {
	var v Value
	err := json.Unmarshal(data, &v)
	return &v, err
}

type KVS interface {
	Open(path string) error
	Close() error
	Set(key string, value *Value) error
	Get(key string) (*Value, error)
	Name() string
}