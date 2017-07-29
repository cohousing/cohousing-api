package db

import (
	"fmt"
	"github.com/cohousing/cohousing-tenant-api/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	dbCache map[string]*gorm.DB = make(map[string]*gorm.DB)
)

func InitDB() {
	config.TenantsLoader = loadTenantsFromDB
}

func GetTenantDB(tenant *config.Tenant) *gorm.DB {
	db := dbCache[tenant.Context]
	if db == nil {
		cfg := config.GetConfig()
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, tenant.Context)
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

func loadTenantsFromDB() []config.Tenant {
	var tenants []config.Tenant

	// TODO How to get tenants now?
	tenants = append(tenants, config.Tenant{
		Context:   "jystrupsavvaerk",
		Name:      "Jystrup Savværk",
		CustomUrl: "intranet.jystrupsavvaerk.dk",
	})

	return tenants
}
