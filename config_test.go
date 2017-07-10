package main

import (
	"testing"
	"time"
)

var (
	origTenantsLoaderFunc TenantsLoaderFunc
	tenantsLoaderFuncMock TenantsLoaderFunc = func() []Tenant {
		return []Tenant{
			Tenant{
				Context: "tenant1",
				Name:    "Tenant 1",
			},
			Tenant{
				Context: "tenant2",
				Name:    "Tenant 2",
			},
			Tenant{
				Context:   "tenant3",
				Name:      "Tenant 3",
				CustomUrl: "customurl.example.com",
			},
		}
	}
)

func TestLoadStaticConfiguration(t *testing.T) {
	LoadStaticConfiguration()

	if config.TenantDomain != "%s.cohousing.nu" {
		t.Error("Configuration file not loaded correctly")
	}
}

func TestRefreshTenantsCache(t *testing.T) {
	tenantsLoaderFunc = tenantsLoaderFuncMock
	tenantCache = nil

	refreshTenantCache()

	if len(tenantCache) != 4 {
		t.Error("Expected cache with 4 tenants")
	}
}

func TestGetTenantByHost(t *testing.T) {
	tenantsLoaderFunc = tenantsLoaderFuncMock
	tenantCache = nil

	tenant := GetTenantByHost("customurl.example.com")

	if tenant == nil {
		t.Error("Expected a tenant back")
	} else if tenant.Context != "tenant3" {
		t.Error("Expected tenant3 back")
	}
}

func TestDynamicConfigRefresher(t *testing.T) {
	tenantsLoaderFunc = tenantsLoaderFuncMock
	tenantCache = nil
	cacheRefresherDuration = 1 * time.Second

	dynamicConfigRefresher()

	time.Sleep(1500 * time.Millisecond)

	if len(tenantCache) != 4 {
		t.Error("Expected cache with 4 tenants")
	}

	close(tenantRefresherQuitter)
}
