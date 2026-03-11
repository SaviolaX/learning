package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type URLPair struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

func (u URLPair) Store(filename string) error {

	// convert storage into bytes
	data, err := json.MarshalIndent(u, "", " ")
	if err != nil {
		return fmt.Errorf("marshal storage: %w", err)
	}

	// write bytes into a file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil

}
