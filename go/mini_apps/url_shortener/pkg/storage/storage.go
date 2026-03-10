package storage

type Storage struct {
	LongUrl  string
	ShortUrl string
}

func (s Storage) Write() {}

func (s Storage) Read() {}
