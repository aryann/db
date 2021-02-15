package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/aryann/db/storage"
)

type JSONStorage struct {
	file *os.File
}

type document struct {
	Key     string `json:"key"`
	Version int64  `json:"version"`
	Payload string `json:"payload"`
}

type documents []document

func (d documents) find(key string) (index int, found bool) {
	start := 0
	limit := len(d)
	for start < limit {
		mid := start + (limit-start)/2
		if d[mid].Key == key {
			return mid, true
		}
		if d[mid].Key < key {
			start = mid + 1
		} else {
			limit = mid - 1
		}
	}

	if start < len(d) {
		return start, d[start].Key == key
	} else {
		return start, false
	}
}

func NewJSONStorage(dir string) (*JSONStorage, error) {
	filepath := path.Join(dir, "data.json")
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	if err := maybeInitFile(file); err != nil {
		return nil, err
	}

	return &JSONStorage{
		file: file,
	}, err
}

func maybeInitFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		data, err := json.Marshal([]document{})
		if err != nil {
			return err
		}
		if _, err := file.Write(data); err != nil {
			return err
		}
		if err := file.Sync(); err != nil {
			return err
		}
	}
	return nil
}

func (j *JSONStorage) currentDocuments() (documents, error) {
	if _, err := j.file.Seek(0, 0); err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(j.file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}
	var current documents
	if err := json.Unmarshal(data, &current); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON: %v", err)
	}
	return current, nil
}

func (j *JSONStorage) Lookup(key string) (version int64, payload string, err error) {
	current, err := j.currentDocuments()
	if err != nil {
		return 0, "", err
	}
	i, found := current.find(key)
	if !found {
		return storage.VersionNotFound, "", nil
	}
	return current[i].Version, current[i].Payload, nil
}

func (j *JSONStorage) Write(key string, version int64, payload string) error {
	current, err := j.currentDocuments()
	if err != nil {
		return err
	}

	newDocument := document{
		Key:     key,
		Version: version,
		Payload: payload,
	}

	i, found := current.find(key)
	if found {
		current[i] = newDocument
	} else {
		current = append(append(current[:i], newDocument), current[i:]...)
	}

	marshaledDocuments, err := json.Marshal(current)
	if err != nil {
		return fmt.Errorf("could not marshal documents: %v", err)
	}
	if _, err := j.file.WriteAt(marshaledDocuments, 0); err != nil {
		return err
	}

	fmt.Println(current)
	return nil
}

func (j *JSONStorage) Close() error {
	return j.file.Close()
}
