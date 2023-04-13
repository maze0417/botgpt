package orm

import (
	"gorm.io/gorm"
)

type BaseOrmI interface {
	GetDBInstance() *gorm.DB
}

type BaseOrm struct {
	db *gorm.DB
}

func NewBaseOrm(db *gorm.DB) *BaseOrm {
	return &BaseOrm{db: db}
}

func (repo *BaseOrm) GetDBInstance() *gorm.DB {
	return repo.db
}
