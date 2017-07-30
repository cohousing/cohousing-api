package config

import (
	"testing"
	"time"
)

var (
	tenantsLoaderFuncMock TenantsLoaderFunc = func() []Tenant {
		return []Tenant{
			{
				Context: "tenant1",
				Name:    "Tenant 1",
			},
			{
				Context: "tenant2",
				Name:    "Tenant 2",
			},
			{
				Context:   "tenant3",
				Name:      "Tenant 3",
				CustomUrl: "customurl.example.com",
			},
		}
	}
)

func TestLoadStaticConfiguration(t *testing.T) {
	configFilePath = "../config.yml"

	LoadStaticConfiguration()

	if GetConfig().TenantDomain != "%s.cohousing.nu" {
		t.Error("Configuration file not loaded correctly")
	}
}

func TestRefreshTenantsCache(t *testing.T) {
	TenantsLoader = tenantsLoaderFuncMock
	tenantCache = nil

	refreshTenantCache()

	if len(tenantCache) != 4 {
		t.Error("Expected cache with 4 tenants")
	}
}

func TestGetTenantByHost(t *testing.T) {
	TenantsLoader = tenantsLoaderFuncMock
	tenantCache = nil

	tenant := GetTenantByHost("customurl.example.com")

	if tenant == nil {
		t.Error("Expected a tenant back")
	} else if tenant.Context != "tenant3" {
		t.Error("Expected tenant3 back")
	}
}

func TestDynamicConfigRefresher(t *testing.T) {
	TenantsLoader = tenantsLoaderFuncMock
	tenantCache = nil
	cacheRefresherDuration = 100 * time.Millisecond

	dynamicConfigRefresher()

	time.Sleep(200 * time.Millisecond)

	if len(tenantCache) != 4 {
		t.Error("Expected cache with 4 tenants")
	}

	close(tenantRefresherQuitter)
}
