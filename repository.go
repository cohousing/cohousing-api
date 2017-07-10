package main

import (
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

func (repository *Repository) GetList(tenant *Tenant) interface{} {
	list := reflect.New(reflect.SliceOf(repository.DomainType)).Interface()
	getTenantDB(tenant).Find(list)
	return list
}

func (repository *Repository) GetById(tenant *Tenant, id uint64) (interface{}, error) {
	item := reflect.New(repository.DomainType).Interface()
	if err := getTenantDB(tenant).First(item, id).Error; err == nil {
		return item, nil
	} else {
		return nil, err
	}
}
