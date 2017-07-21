package db

import (
	"database/sql"
	"fmt"
	"github.com/cohousing/cohousing-api/db/conf"
	"github.com/cohousing/cohousing-api/db/tenant"
	"github.com/rubenv/sql-migrate"
	"os"
)

type StdOutMigrateLogger struct {
}

func (l StdOutMigrateLogger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stdout, format, v)
}

func (l StdOutMigrateLogger) Verbose() bool {
	return true
}

func MigrateConfDB(confDB *sql.DB) error {
	return migrateDB(confDB, &migrate.AssetMigrationSource{
		Asset:    conf.Asset,
		AssetDir: conf.AssetDir,
		Dir:      "db/conf",
	})
}

func MigrateTenantDB(tenantDB *sql.DB) error {
	return migrateDB(tenantDB, &migrate.AssetMigrationSource{
		Asset:    tenant.Asset,
		AssetDir: tenant.AssetDir,
		Dir:      "db/tenant",
	})
}

func migrateDB(db *sql.DB, migrationSource migrate.MigrationSource) error {
	n, err := migrate.Exec(db, "mysql", migrationSource, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Applied %d migrations!\n", n)
	return nil
}
