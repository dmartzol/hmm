package storage

type Storage interface {
}

type MacroStorage struct {
	db Storage
}

func NewStorage(db Storage) *MacroStorage {
	return &MacroStorage{
		db: db,
	}
}
