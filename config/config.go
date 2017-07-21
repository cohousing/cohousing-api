package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	tenantCache            map[string]*Tenant
	cacheRefresherDuration time.Duration = 1 * time.Minute
	tenantRefresherQuitter chan struct{}

	Loaded = struct {
		TenantDomain string

		ConfigDB struct {
			User     string
			Password string
			Host     string
			Port     uint
			Name     string
		}
	}{}

	TenantsLoader TenantsLoaderFunc
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
	return tenantCache[host]
}

func LoadStaticConfiguration() {
	configor.Load(&Loaded, "config.yml")
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
	tenants := TenantsLoader()

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
	return fmt.Sprintf(Loaded.TenantDomain, context)
}