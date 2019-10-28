package infra

import (
    "fmt"
    "os"
    "sync"
    "database/sql"
    log "github.com/sirupsen/logrus"
    _ "github.com/lib/pq"
)


var pgLogger = log.WithFields(log.Fields{"ctx": "postgres"})


type dbConfig struct {
    host                string
    port                string
    dbName              string
    user                string
    password            string
}

func (config *dbConfig) connectionUri() string {
    user := config.user
    if config.password != "" {
        user = fmt.Sprintf("%s:%s", config.user, config.password)
    }
    hostAddr := fmt.Sprintf("%s:%s", config.host, config.port)
    return fmt.Sprintf("postgres://%s@%s/%s?sslmode=disable", user, hostAddr, config.dbName)
}

func defaultDbConfig() *dbConfig {
    return &dbConfig{
        user: os.Getenv("GOREST_DEFAULT_DB_USER"),
        password: os.Getenv("GOREST_DEFAULT_DB_PASSWORD"),
        host: os.Getenv("GOREST_DEFAULT_DB_HOST"),
        port: os.Getenv("GOREST_DEFAULT_DB_PORT"),
        dbName: os.Getenv("GOREST_DEFAULT_DB_NAME"),
    }
}


/**
 * --- Connection Pool Singleton ---
 * 
 * Database connection pool is managed by creating a singleton of
 * `connectionPool` struct which maintains a lazy connection pool
 * on `connectionPool.db`. 
 */

type connectionPool struct {
    db  *sql.DB  // Database connection pool
}

func newConnectionPool(config *dbConfig) *connectionPool {

    pgConnLogger := pgLogger.WithFields(log.Fields{
        "database": config.dbName,
        "host": config.host,
    })

    db, err := sql.Open("postgres", config.connectionUri())
    if err != nil {
        pgConnLogger.Panicf("Unable to establish connection to database. " +
                            "Error: %v", err)
    }

    err = db.Ping()  // send a ping to ensure connection is established
    if err != nil {
        pgConnLogger.Panicf("Unable to establish connection to database. " +
                            "Error: %v", err)
    }

    pgConnLogger.Infof("Connection established to database.")
    return &connectionPool{db}
}


var pgDefaultConnPool *connectionPool
var once *sync.Once = nil

// Returns a connection pool object from the defualt connection
func DefaultDb() *sql.DB {

    if once == nil {
        once = new(sync.Once)
    }

    once.Do(func() {
        config := defaultDbConfig()
        pgDefaultConnPool = newConnectionPool(config)
    })
    return pgDefaultConnPool.db
}

// Closes the default database connection pool
func CloseDefaultDb() {
    if pgDefaultConnPool != nil {
        err := pgDefaultConnPool.db.Close()
        if err != nil {
            pgLogger.Warnf("Error closing default connection. Error: %v", err)
        }
        pgLogger.Infof("Default connection closed!")
        once = nil
    }
}
