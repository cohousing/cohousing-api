package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"time"
)

type Tenant struct {
	Context   string
	Name      string
	TenantUrl string
	CustomUrl string
}

var (
	tenantCache            map[string]*Tenant
	cacheRefresherDuration time.Duration = 1 * time.Minute
	tenantRefresherQuitter chan struct{}
	configFilePath         string = "config.yml"

	loaded = Config{}

	TenantsLoader TenantsLoaderFunc
)

type Config struct {
	TenantDomain string

	DB struct {
		User     string
		Password string
		Host     string
		Port     uint
		Name     string
	}
}

func GetConfig() *Config {
	return &loaded
}

type TenantsLoaderFunc func() []Tenant

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
	configor.Load(&loaded, configFilePath)
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
		tenantUrl := buildTenantDomain(tenant.Context)
		(&tenant).TenantUrl = tenantUrl
		newTenantCache[tenantUrl] = &tenant
		if tenant.CustomUrl != "" {
			newTenantCache[tenant.CustomUrl] = &tenant
		}
	}
	tenantCache = newTenantCache
}

func buildTenantDomain(context string) string {
	return fmt.Sprintf(GetConfig().TenantDomain, context)
}
