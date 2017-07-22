package db

import (
	"fmt"
	"github.com/cohousing/cohousing-api/config"
	"github.com/cohousing/cohousing-api/domain/admin"
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

func GetTenantDB(tenant *admin.Tenant) *gorm.DB {
	db := dbCache[tenant.Context]
	if db == nil {
		cfg := config.GetConfig()
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.ConfigDB.User, cfg.ConfigDB.Password, cfg.ConfigDB.Host, cfg.ConfigDB.Port, tenant.Context)
		var err error
		db, err = gorm.Open("mysql", connectionString)
		if err != nil {
			panic(err)
		}
		dbCache[tenant.Context] = db
		if err = MigrateTenantDB(db.DB()); err != nil {
			panic(err)
		}
		db.LogMode(true)
	}
	return db
}

func GetConfDB() *gorm.DB {
	if confDB == nil {
		var err error
		cfg := config.GetConfig()
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.ConfigDB.User, cfg.ConfigDB.Password, cfg.ConfigDB.Host, cfg.ConfigDB.Port, cfg.ConfigDB.Name)
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

func loadTenantsFromDB() []admin.Tenant {
	var tenants []admin.Tenant

	GetConfDB().Find(&tenants)

	return tenants
}
