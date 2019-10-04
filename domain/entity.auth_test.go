package domain

import (
    "encoding/json"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)


func TestNewAccessType(t *testing.T) {

    assert := assert.New(t)
    var access AccessType
    var err error

    access, err = NewAccessType("read")
    assert.Nil(err)
    assert.Equal(ReadAccess, access)

    access, err = NewAccessType("READ")
    assert.Nil(err)
    assert.Equal(ReadAccess, access)

    access, err = NewAccessType("Read")
    assert.Nil(err)
    assert.Equal(ReadAccess, access)

    access, err = NewAccessType("Write")
    assert.Nil(err)
    assert.Equal(WriteAccess, access)

    access, err = NewAccessType("Admin")
    assert.Nil(err)
    assert.Equal(AdminAccess, access)


    access, err = NewAccessType("foo")
    assert.Equal(ErrInvalidAccessType, err)
    assert.Equal(AccessType(""), access)
}

func TestAccessTypeJSONUnmarshal(t *testing.T) {

    assert := assert.New(t)
    var access AccessType
    var err error

    err = json.Unmarshal([]byte(`"read"`), &access)
    assert.Nil(err)
    assert.Equal(ReadAccess, access)

    err = json.Unmarshal([]byte(`"write"`), &access)
    assert.Nil(err)
    assert.Equal(WriteAccess, access)

    err = json.Unmarshal([]byte(`"admin"`), &access)
    assert.Nil(err)
    assert.Equal(AdminAccess, access)

    err = json.Unmarshal([]byte(`"foo"`), &access)
    assert.Error(err)
    assert.Equal(ErrInvalidAccessType, err)

    err = json.Unmarshal([]byte(`admin`), &access)
    assert.Error(err)

    type User struct {
        Access AccessType `json:"access"`
    }

    var user User
    err  = json.Unmarshal([]byte(`{"access": "write"}`), &user)
    assert.Nil(err)
    assert.Equal(WriteAccess, user.Access)
}

func TestNewAuthToken(t *testing.T)  {

    assert := assert.New(t)
    var authToken *AuthToken = NewAuthToken("foobar", ReadAccess)

    assert.Equal("foobar", authToken.Username)
    assert.Equal(ReadAccess, authToken.Access)
    assert.True(authToken.ExpiresAt > authToken.IssuedAt)

}

func TestAuthTokenRefresh(t *testing.T)  {

    assert := assert.New(t)
    var authToken *AuthToken = NewAuthToken("foobar", ReadAccess)
    var err error

    expiry := authToken.ExpiresAt

    time.Sleep(time.Second * 1)
    err = authToken.Refresh()
    assert.Nil(err)
    expiryAfterRefresh := authToken.ExpiresAt

    assert.True(expiryAfterRefresh > expiry)

    // Test refresh when token has already expired

    expiryBeforeRefresh := time.Now().Unix()

    authToken.ExpiresAt = expiryBeforeRefresh
    time.Sleep(time.Second * 1)
    err = authToken.Refresh()
    assert.Equal(ErrTokenExpired, err)
    assert.Equal(expiryBeforeRefresh, authToken.ExpiresAt)
}

func TestAuthTokenIsExpired(t *testing.T)  {

    assert := assert.New(t)
    var authToken *AuthToken = NewAuthToken("foobar", ReadAccess)

    assert.False(authToken.IsExpired())

    authToken.ExpiresAt = time.Now().Unix()
    time.Sleep(time.Second * 1)
    assert.True(authToken.IsExpired())

}

func TestAuthTokenAccess(t *testing.T)  {

    assert := assert.New(t)
    var authToken *AuthToken = NewAuthToken("foobar", ReadAccess)

    assert.False(authToken.HasWriteAccess())
    assert.False(authToken.HasAdminAccess())

    authToken = NewAuthToken("foobar", WriteAccess)
    assert.True(authToken.HasWriteAccess())
    assert.False(authToken.HasAdminAccess())

    authToken = NewAuthToken("foobar", AdminAccess)
    assert.True(authToken.HasWriteAccess())
    assert.True(authToken.HasAdminAccess())
}
