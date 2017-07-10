package main

import "testing"

func TestLoadStaticConfiguration(t *testing.T) {
	LoadStaticConfiguration()

	if config.TenantDomain != "%s.cohousing.nu" {
		t.Error("Configuration file not loaded correctly")
	}
}
