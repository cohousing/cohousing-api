package db

import (
	"fmt"
	"github.com/cohousing/cohousing-api/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	dbCache map[string]*gorm.DB = make(map[string]*gorm.DB)
	confDB  *gorm.DB
)

func InitDB() {
	config.TenantsLoader = loadTenantsFromDB
}

func GetTenantDB(tenant *config.Tenant) *gorm.DB {
	db := dbCache[tenant.Context]
	if db == nil {
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.Loaded.ConfigDB.User, config.Loaded.ConfigDB.Password, config.Loaded.ConfigDB.Host, config.Loaded.ConfigDB.Port, tenant.Context)
		var err error
		db, err = gorm.Open("mysql", connectionString)
		if err != nil {
			panic(err)
		}
		dbCache[tenant.Context] = db
		if err = MigrateTenantDB(db.DB()); err != nil {
			panic(err)
		}
	}
	return db
}

func GetConfDB() *gorm.DB {
	if confDB == nil {
		var err error
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.Loaded.ConfigDB.User, config.Loaded.ConfigDB.Password, config.Loaded.ConfigDB.Host, config.Loaded.ConfigDB.Port, config.Loaded.ConfigDB.Name)
		confDB, err = gorm.Open("mysql", connectionString)
		if err != nil {
			panic(err)
		}
		if err = MigrateConfDB(confDB.DB()); err != nil {
			panic(err)
		}
	}
	return confDB
}

func loadTenantsFromDB() []config.Tenant {
	var tenants []config.Tenant

	GetConfDB().Find(&tenants)

	return tenants
}
