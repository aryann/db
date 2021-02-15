package storage

type Storage interface {
	Insert(key string, payload string) error
	Close() error
}
