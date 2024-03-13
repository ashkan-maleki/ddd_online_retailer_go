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

type MigrationInvoker func(db *gorm.DB)

//db.AutoMigrate(&Batches{})

type GormDbDriver string

const (
	InMemory GormDbDriver = "file::memory:?cache=shared"
)

func NewGormRepository[E any](driver GormDbDriver, migrate MigrationInvoker, config *gorm.Config) *GormRepository[E] {
	db, err := gorm.Open(sqlite.Open(string(driver)), config)
	if err != nil {
		panic("failed to connect database")
	}
	migrate(db)
	return &GormRepository[E]{db: db}
}

var _ AbstractRepository[any] = (*GormRepository[any])(nil)

func (repo GormRepository[E]) Insert(ctx context.Context, entity *E) error {
	repo.db.WithContext(ctx).Create(entity)
	return nil
}

func (repo GormRepository[E]) Delete(ctx context.Context, entity *E) error {
	repo.db.WithContext(ctx).Delete(entity)
	return nil
}

func (repo GormRepository[E]) DeleteById(ctx context.Context, id any) error {
	entity, err := repo.Get(ctx, id)
	if err != nil {
		return err
	}
	err = repo.Delete(ctx, &entity)
	if err != nil {
		return err
	}
	return nil
}

func (repo GormRepository[E]) Update(ctx context.Context, entity *E) error {
	repo.db.WithContext(ctx).Save(entity)
	return nil
}

func (repo GormRepository[E]) List(ctx context.Context) ([]E, error) {
	var entityList []E
	repo.db.WithContext(ctx).Find(&entityList)
	return entityList, nil
}

func (repo GormRepository[E]) Get(ctx context.Context, id any) (E, error) {
	var entity E
	repo.db.WithContext(ctx).First(&entity, id)
	return entity, nil
}

func (repo GormRepository[E]) Find(ctx context.Context, specifications ...Specification) ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (repo GormRepository[E]) Count(ctx context.Context, specifications ...Specification) (i int64, err error) {
	//TODO implement me
	panic("implement me")
}
