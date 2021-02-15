package storage

const VersionNotFound = -1

type Storage interface {
	Lookup(key string) (version int64, payload string, err error)
	Write(key string, version int64, payload string) error
	Close() error
}
