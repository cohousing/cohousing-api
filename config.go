package main

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

var (
	tenantCache            map[string]*Tenant
	cacheRefresherDuration time.Duration = 1 * time.Minute
	tenantRefresherQuitter chan struct{}

	config = struct {
		TenantDomain string

		ConfigDB struct {
			User     string
			Password string
			Host     string
			Port     uint
			Name     string
		}
	}{}

	tenantsLoaderFunc TenantsLoaderFunc = loadTenantsFromDB
)

type TenantsLoaderFunc func() []Tenant

type Tenant struct {
	gorm.Model
	Context   string `gorm:"size:100"`
	Name      string
	CustomUrl string
}

func Configure() {
	LoadStaticConfiguration()

	dynamicConfigRefresher()
}

func GetTenantByHost(host string) *Tenant {
	if tenantCache == nil {
		refreshTenantCache()
	}
	fmt.Fprintf(os.Stdout, "Looking up tenant for host: %s\n", host)
	return tenantCache[host]
}

func LoadStaticConfiguration() {
	configor.Load(&config, "config.yml")
}

func dynamicConfigRefresher() {
	ticker := time.NewTicker(cacheRefresherDuration)
	tenantRefresherQuitter = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				refreshTenantCache()
			case <-tenantRefresherQuitter:
				ticker.Stop()
				return
			}
		}
	}()
}

func refreshTenantCache() {
	tenants := tenantsLoaderFunc()

	newTenantCache := make(map[string]*Tenant)
	for _, tenant := range tenants {
		newTenantCache[buildTenantDomain(tenant.Context)] = &tenant
		if tenant.CustomUrl != "" {
			newTenantCache[tenant.CustomUrl] = &tenant
		}
	}
	tenantCache = newTenantCache
	fmt.Fprintf(os.Stdout, "Tenant cache refreshed %v\n", tenantCache)
}

func buildTenantDomain(context string) string {
	return fmt.Sprintf(config.TenantDomain, context)
}

func loadTenantsFromDB() []Tenant {
	var tenants []Tenant

	getConfDB().Find(&tenants)

	return tenants
}
