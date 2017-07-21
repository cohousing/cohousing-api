package db

import (
	"github.com/cohousing/cohousing-api/config"
	"reflect"
)

type Repository struct {
	DomainType reflect.Type
}

func CreateRepository(domainType reflect.Type) *Repository {
	return &Repository{
		DomainType: domainType,
	}
}

func (repository *Repository) GetList(tenant *config.Tenant) interface{} {
	list := reflect.New(reflect.SliceOf(repository.DomainType)).Interface()
	GetTenantDB(tenant).Find(list)
	return list
}

func (repository *Repository) GetById(tenant *config.Tenant, id uint64) (interface{}, error) {
	object := reflect.New(repository.DomainType).Interface()
	if err := GetTenantDB(tenant).First(object, id).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Create(tenant *config.Tenant, object interface{}) (interface{}, error) {
	if err := GetTenantDB(tenant).Create(object).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Update(tenant *config.Tenant, object interface{}) (interface{}, error) {
	if err := GetTenantDB(tenant).Save(object).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Delete(tenant *config.Tenant, id uint64) error {
	item := reflect.New(repository.DomainType).Interface()
	return GetTenantDB(tenant).Delete(item, id).Error
}
