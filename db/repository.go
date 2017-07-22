package db

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"reflect"
)

type DBFactory func(c *gin.Context) *gorm.DB

type Repository struct {
	DomainType reflect.Type
	DBFactory  DBFactory
}

func CreateRepository(domainType reflect.Type, dbFactory DBFactory) *Repository {
	return &Repository{
		DomainType: domainType,
		DBFactory:  dbFactory,
	}
}

func (repository *Repository) GetList(c *gin.Context, lookupObject interface{}, start, limit int) (interface{}, int) {
	list := reflect.New(reflect.SliceOf(repository.DomainType)).Interface()
	var count int
	repository.DBFactory(c).Model(reflect.New(repository.DomainType).Interface()).Where(lookupObject).Count(&count)
	if count > 0 {
		repository.DBFactory(c).Where(lookupObject).Offset(start).Limit(limit).Find(list)
	}
	return list, count
}

func (repository *Repository) GetById(c *gin.Context, id uint64) (interface{}, error) {
	object := reflect.New(repository.DomainType).Interface()
	if err := repository.DBFactory(c).First(object, id).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Create(c *gin.Context, object interface{}) (interface{}, error) {
	if err := repository.DBFactory(c).Create(object).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Update(c *gin.Context, object interface{}) (interface{}, error) {
	if err := repository.DBFactory(c).Save(object).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Delete(c *gin.Context, id uint64) error {
	item := reflect.New(repository.DomainType).Interface()
	return repository.DBFactory(c).Delete(item, id).Error
}
