package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	tenantCache map[string]*Tenant
	config      = struct {
		TenantDomain string

		ConfigDB struct {
			User     string
			Password string
			Host     string
			Port     uint
			Name     string
		}
	}{}
)

type Tenant struct {
	gorm.Model
	Context    string `gorm:"size:100"`
	Name       string
	DbUsername string `gorm:"size:32"`
	DbPassword string `gorm:"size:32"`
	DbHostname string `gorm:"size:50"`
	DbDatabase string `gorm:"size:64"`
	CustomUrl  string
}

func Configure() {
	loadStaticConfiguration()

	dynamicConfigRefresher()
}

func GetTenantByHost(host string) *Tenant {
	if tenantCache == nil {
		refreshTenantCache()
	}
	return tenantCache[host]
}

func loadStaticConfiguration() {
	configor.Load(&config, "config.yml")
}

func dynamicConfigRefresher() {
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				refreshTenantCache()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	refreshTenantCache()
}

func refreshTenantCache() {
	var tenants []Tenant

	getConfDB().Find(&tenants)

	newTenantCache := make(map[string]*Tenant)
	for _, tenant := range tenants {
		newTenantCache[buildTenantDomain(tenant.Context)] = &tenant
		if tenant.CustomUrl != "" {
			newTenantCache[tenant.CustomUrl] = &tenant
		}
	}
	tenantCache = newTenantCache
}

func buildTenantDomain(context string) string {
	return fmt.Sprintf(config.TenantDomain, context)
}
