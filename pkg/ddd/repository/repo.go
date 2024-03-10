package repository

import (
	"context"
	"reflect"
)

type Repository[E any] interface {
	GetEntityType() reflect.Type
}

type BaseRepository[E any] struct {
}

func (repo BaseRepository[E]) GetEntityType() reflect.Type {
	return reflect.TypeOf(repo)
}

type Inserter[E any] interface {
	Repository[E]
	Insert(ctx context.Context, entity *E) error
}

type Deleter[E any] interface {
	Repository[E]
	Delete(ctx context.Context, entity *E) error
	DeleteById(ctx context.Context, id any) error
}

type Updater[E any] interface {
	Repository[E]
	Update(ctx context.Context, entity *E) error
}

type Lister[E any] interface {
	Repository[E]
	List(ctx context.Context) ([]E, error)
}

type Getter[E any] interface {
	Repository[E]
	Get(ctx context.Context, id any) (E, error)
}

type CrudRepository[E any] interface {
	Repository[E]
	Inserter[E]
	Deleter[E]
	Updater[E]
	Lister[E]
	Getter[E]
}

type Specifier[E any] interface {
	Repository[E]
	Find(ctx context.Context, specifications ...Specification) ([]E, error)
}

type Counter[E any] interface {
	Repository[E]
	Count(ctx context.Context, specifications ...Specification) (i int64, err error)
}

type GenericRepository[E any] interface {
	Insert(ctx context.Context, entity *E) error
	Delete(ctx context.Context, entity *E) error
	DeleteById(ctx context.Context, id any) error
	Update(ctx context.Context, entity *E) error
	Get(ctx context.Context, id any) (E, error)
	List(ctx context.Context) ([]E, error)
}

type ComplexRepository[E any] interface {
	Find(ctx context.Context, specifications ...Specification) ([]E, error)
	FindWithLimit(ctx context.Context, limit int, offset int, specifications ...Specification) ([]E, error)
	Count(ctx context.Context, specifications ...Specification) (i int64, err error)
}
