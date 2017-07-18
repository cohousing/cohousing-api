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
	object := reflect.New(repository.DomainType).Interface()
	if err := getTenantDB(tenant).First(object, id).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Create(tenant *Tenant, object interface{}) (interface{}, error) {
	if err := getTenantDB(tenant).Create(object).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Update(tenant *Tenant, object interface{}) (interface{}, error) {
	if err := getTenantDB(tenant).Save(object).Error; err == nil {
		return object, nil
	} else {
		return nil, err
	}
}

func (repository *Repository) Delete(tenant *Tenant, id uint64) error {
	item := reflect.New(repository.DomainType).Interface()
	return getTenantDB(tenant).Delete(item, id).Error
}
