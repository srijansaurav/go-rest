package infra

import (
    "errors"
    "database/sql"
    "github.com/lib/pq"
    log "github.com/sirupsen/logrus"
    "github.com/esquarer/go-rest/domain"
)


var (
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrUserNotFound = errors.New("user not found")
)

// Concerte implementation of `domain.UserRepository`
// for PostgreSQL database.
type PgUserRepository struct {
    db  *sql.DB
}

func NewUserRepository() domain.UserRepository {
    db := DefaultDb()
    return &PgUserRepository{db}
}

func (repo *PgUserRepository) Add(user *domain.User) error {

    var id uint64
    query := `insert into auth_user (username, password) values ($1, $2) returning id;`
    err := repo.db.QueryRow(query, user.Username, user.Password).Scan(&id)

    if err == nil {
        pgLogger.WithFields(log.Fields{
            "id": id, "username": user.Username,
        }).Info("User created")
    }

    // Check for unique contraint violation
    if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
        return ErrUserAlreadyExists
    }

    return err
}

func (repo *PgUserRepository) FindByUsername(username string) (*domain.User, error) {

    query := `select username, password from auth_user where username = $1;`
    
    user := &domain.User{}
    err := repo.db.QueryRow(query, username).Scan(&user.Username, &user.Password)

    if err == sql.ErrNoRows { err = ErrUserNotFound }
    return user, err
}
