package entities

type EntityList interface {
	List(filter string)
}

type EntityDelete interface {
	Delete(id string) error
}
