package infra

import (
    "fmt"
    "os"
    "strings"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)


func runMigrations(config *dbConfig, mode string) error {

    migrationsFilePath := os.Getenv("GOREST_MIGRATIONS_FILE_PATH")

    filePath := fmt.Sprintf("file://%s", migrationsFilePath)
    m, err := migrate.New(filePath, config.connectionUri())
    defer m.Close()
    if err != nil { panic(err) }

    mode = strings.ToLower(mode)
    if mode == "up" {
        err = m.Up()
    } else if mode == "down" {
        err = m.Down()
    } else {
        err = fmt.Errorf("invalid migration mode '%v'", mode)
    }

    // DO NOT consider ErrNoChange as migration failure.
    // Simply ignore it.
    if err == migrate.ErrNoChange { err = nil }
    return err
}

// Apply all migrations under migrations file path
func RunDatabaseMigrations() error {

    config := defaultDbConfig()
    err := runMigrations(config, "up")
    if err != nil {
        pgLogger.Panicf("Migrations failed! Error: %v", err)
        return err
    } else {
        pgLogger.Infof("Migrations applied successfully")
        return nil
    }
}

// Revert all migrations under migrations file path
// NOTE: This function is used only while testing to cleanup database
func RevertDatabaseMigrations() error {

    config := defaultDbConfig()
    err := runMigrations(config, "down")
    if err != nil {
        pgLogger.Panicf("Migrations drop failed! Error: %v", err)
        return err
    } else {
        pgLogger .Infof("Migrations dropped successfully")
        return nil
    }
}
