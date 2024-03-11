package repository

import (
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormRepository[E any] struct {
	BaseRepository[E]
	db *gorm.DB
}

func NewGormRepository[E any]() *GormRepository[E] {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	// Todo: Add migrations
	//db.AutoMigrate(&Batches{})
	return &GormRepository[E]{db: db}
}

var _ AbstractRepository[any] = (*GormRepository[any])(nil)

func (repo GormRepository[E]) Insert(ctx context.Context, entity *E) error {
	repo.db.Create(entity)
	return nil
}

func (repo GormRepository[E]) Delete(ctx context.Context, entity *E) error {
	//TODO implement me
	panic("implement me")
}

func (repo GormRepository[E]) DeleteById(ctx context.Context, id any) error {
	//TODO implement me
	panic("implement me")
}

func (repo GormRepository[E]) Update(ctx context.Context, entity *E) error {
	//TODO implement me
	panic("implement me")
}

func (repo GormRepository[E]) List(ctx context.Context) ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (repo GormRepository[E]) Get(ctx context.Context, id any) (E, error) {
	//TODO implement me
	panic("implement me")
}

func (repo GormRepository[E]) Find(ctx context.Context, specifications ...Specification) ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (repo GormRepository[E]) Count(ctx context.Context, specifications ...Specification) (i int64, err error) {
	//TODO implement me
	panic("implement me")
}
