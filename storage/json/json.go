package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func (d documents) find(newDocument document) (index int, found bool) {
	start := 0
	limit := len(d)
	for start < limit {
		mid := start + (limit-start)/2
		if d[mid].Key == newDocument.Key {
			return mid, true
		}
		if d[mid].Key < newDocument.Key {
			start = mid + 1
		} else {
			limit = mid - 1
		}
	}

	if start < len(d) {
		return start, d[start].Key == newDocument.Key
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

func (j *JSONStorage) Insert(key string, payload string) error {
	if _, err := j.file.Seek(0, 0); err != nil {
		return err
	}
	data, err := ioutil.ReadAll(j.file)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}
	var current documents
	if err := json.Unmarshal(data, &current); err != nil {
		return fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	newDocument := document{
		Key:     key,
		Version: 0,
		Payload: payload,
	}

	insertAt, found := current.find(newDocument)
	if found {
		return fmt.Errorf("document with key '%s' already exists", key)
	}

	newDocuments := append(append(current[:insertAt], newDocument), current[insertAt:]...)
	marshaledDocuments, err := json.Marshal(newDocuments)
	if err != nil {
		return fmt.Errorf("could not marshal documents: %v", err)
	}
	if _, err := j.file.WriteAt(marshaledDocuments, 0); err != nil {
		return err
	}

	fmt.Println(newDocuments)
	return nil
}

func (j *JSONStorage) Close() error {
	return j.file.Close()
}
