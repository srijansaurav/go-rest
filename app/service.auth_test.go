package app

import (
    "errors"
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/esquarer/go-rest/domain"
)

var testPassword = "1q2w3e4r5t"
var jwtSecret = "305f3925544d4bdaa1f7e79b01ced273"

// Dummy implementation of UserRepository
// Returns successful response for all methods
type DummyUserRepository struct {}
func (r *DummyUserRepository) Add(user *domain.User) error { return nil }
func (r *DummyUserRepository) FindByUsername(username string) (*domain.User, error) {
    user, _ := domain.NewUser(username, testPassword)
    return user, nil
}

// Dummy implementation of UserRepository
// Returns error and failed response for all methods
type DummyErrorUserRepository struct {}
func (r *DummyErrorUserRepository) Add(user *domain.User) error { return errors.New("username exists") }
func (r *DummyErrorUserRepository) FindByUsername(username string) (*domain.User, error) { return nil, errors.New("not found") }


var repo domain.UserRepository    = &DummyUserRepository{}
var errRepo domain.UserRepository = &DummyErrorUserRepository{}


func TestNewJWTService(t *testing.T) {

    assert := assert.New(t)
    
    service, err := NewJWTService(repo)
    assert.Nil(service)
    assert.Error(err)

    _ = os.Setenv("JWT_SECRET", jwtSecret)
    service, err = NewJWTService(repo)
    assert.Nil(err)
}


func TestJWTServiceGenerateToken(t *testing.T) {

    assert := assert.New(t)
    _ = os.Setenv("JWT_SECRET", jwtSecret)
    service, _ := NewJWTService(repo)

    user := &domain.User{
        Username: "foobar",
        Password: "1q2w3e4r5t6y7u8i9o0p",
    }

    token, expiry, err := service.GenerateToken(user)

    assert.NotNil(token)
    assert.NotNil(expiry)
    assert.Nil(err)
}


func TestJWTServiceValidateToken(t *testing.T) {

    assert := assert.New(t)
    _ = os.Setenv("JWT_SECRET", jwtSecret)
    service, _ := NewJWTService(repo)

    user := &domain.User{
        Username: "foobar",
        Password: "1q2w3e4r5t6y7u8i9o0p",
    }

    token, expiry, err := service.GenerateToken(user)
    assert.NotNil(token)
    assert.NotNil(expiry)
    assert.Nil(err)

    validatedUser, err := service.ValidateToken(token)
    assert.Nil(err)
    assert.Equal(user.Username, validatedUser.Username)

    // Invalid token sent
    validatedUser, err = service.ValidateToken(token + ".foo")
    assert.Nil(validatedUser)
    assert.Equal(ErrInvalidToken, err)
}


func TestUserServiceRegister(t *testing.T) {

    assert := assert.New(t)
    service := NewUserService(repo)

    var user *domain.User
    var err error

    user, err = service.Register("jo", "0p9o8i7u")
    assert.Nil(user)
    assert.Equal(err, ErrUsernameLength)

    user, err = service.Register("12345678901234567", "0p9o8i7u")
    assert.Nil(user)
    assert.Equal(err, ErrUsernameLength)

    user, err = service.Register("johndoe", "12345")
    assert.Nil(user)
    assert.Equal(err, ErrPasswordLength)

    user, err = service.Register("johndoe", "d81c9a8765d342d699b8bbca2a74b58a1")
    assert.Nil(user)
    assert.Equal(err, ErrPasswordLength)

    user, err = service.Register("johndoe", "d81c9a8765d342d699b8bbca2a74b58a")
    assert.Nil(err)
    assert.Equal("johndoe", user.Username)
    assert.True(len(user.Password) > 0)
}


func TestUserServiceAuthenticate(t *testing.T) {

    assert := assert.New(t)
    service := NewUserService(repo)
    var err error

    // Username and Password is correct
    err = service.Authenticate("foobar", testPassword)
    assert.Nil(err)

    // Password is incorrect
    err = service.Authenticate("foobar", "foobar")
    assert.Equal(ErrIncorrectCredentials, err)

    errRepo := &DummyErrorUserRepository{}
    service = NewUserService(errRepo)

    // Username does not exists
    err = service.Authenticate("foobar", testPassword)
    assert.Equal(ErrIncorrectCredentials, err)
}
