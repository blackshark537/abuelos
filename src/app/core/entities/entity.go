package entities

type EntityList interface {
	List(filter string)
}

type EntityDelete interface {
	Delete(id string) error
}

type ObjectId interface {
	Hex() string
}

type Entity interface {
	GetDbPort() interface{}
	Count(filter string) (int64, error)
	Save() (interface{}, error)
	Update(id string) error
	Delete(id string) error
	DeleteMany(filters string) error

	GetAll(filters string) ([]interface{}, error)
	FindOne(filters string) error

	List(filter string)
}
