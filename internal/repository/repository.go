package repository

type Repository interface {
	Close()
	FindByLongUrl(longUrl string) (*MappingModel, error)
	FindByKey(key string) (*MappingModel, error)
	Create(mapping *MappingModel) error
	GetLastId() uint
}

type MappingModel struct {
	ID      string
	Key     string
	LongUrl string
}
