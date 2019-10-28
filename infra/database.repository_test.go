package infra

import (
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/esquarer/go-rest/domain"
)


func TestMain(m *testing.M) {

    testDbName := os.Getenv("GOREST_TEST_DB_NAME")
    os.Setenv("GOREST_DEFAULT_DB_NAME", testDbName)

    testDbUser := os.Getenv("GOREST_TEST_DB_USER")
    os.Setenv("GOREST_DEFAULT_DB_USER", testDbUser)

    testDbPassword := os.Getenv("GOREST_TEST_DB_PASSWORD")
    os.Setenv("GOREST_DEFAULT_DB_PASSWORD", testDbPassword)

    // Before applying migrations, remove any dirty data in
    // database if exists
    RevertDatabaseMigrations()
    RunDatabaseMigrations()
    code := m.Run()
    // Delete all tables form the database
    RevertDatabaseMigrations()
    os.Exit(code)
}

func ClearUsers() {
    db := DefaultDb()
    _, err := db.Exec("delete from auth_user")
    if err != nil { panic(err) }
}

// Utility function to check total number of user
// objects in the database
func totalUserCount() int {
    db, count := DefaultDb(), 0
    _ = db.QueryRow("select count(*) from auth_user").Scan(&count)
    return count
}


func TestUserRepository(t *testing.T) {

    assert := assert.New(t)
    repo := NewUserRepository()
    ClearUsers()

    user, _ := domain.NewUser("johndoe", "0p9o8i7u6y")

    // Create a new user with username "johndoe"
    err := repo.Add(user)
    assert.Nil(err)
    // Ensure that 1 user exists in database
    assert.Equal(1, totalUserCount())

    // Searching for user with username "johndoe" should return
    // the appropriate user
    searchedUser, err := repo.FindByUsername(user.Username)
    assert.Nil(err)
    assert.Equal(user.Username, searchedUser.Username)
    assert.Equal(user.Password, searchedUser.Password)
}

func TestErrUserAlreadyExists(t *testing.T) {

    assert := assert.New(t)
    repo := NewUserRepository()
    ClearUsers()

    // Create user with username "johndoe"
    user, _ := domain.NewUser("johndoe", "0p9o8i7u6y")
    err := repo.Add(user)
    assert.Nil(err)
    // Creating another user with username "johndoe" should
    // return error ErrUserAlreadyExists
    err = repo.Add(user)
    assert.Equal(ErrUserAlreadyExists, err)
    assert.Equal(1, totalUserCount())
}

func TestErrUserNotFound(t *testing.T) {

    assert := assert.New(t)
    repo := NewUserRepository()
    ClearUsers()

    // Empty database. Create user with username "johndoe"
    user, _ := domain.NewUser("johndoe", "0p9o8i7u6y")
    err := repo.Add(user)
    assert.Nil(err)
    // Only one user (johndoe) exists in database. Searching
    // for "foobar" should return error ErrUserNotFound
    _, err = repo.FindByUsername("foobar")
    assert.Equal(ErrUserNotFound, err)
}
