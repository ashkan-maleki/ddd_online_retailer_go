package repository

type Specification interface {
	GetQuery() string
	GetValues() []any
}
