package db

import (
	"database/sql"
	"fmt"
	tenantMigrate "github.com/cohousing/cohousing-tenant-api/db/migrate"
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

func MigrateTenantDB(tenantDB *sql.DB) error {
	return migrateDB(tenantDB, &migrate.AssetMigrationSource{
		Asset:    tenantMigrate.Asset,
		AssetDir: tenantMigrate.AssetDir,
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
