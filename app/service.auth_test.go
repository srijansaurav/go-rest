package app

import (
    "time"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/esquarer/go-rest/domain"
)


func TestJWTServiceGenerateToken(t *testing.T) {

    assert := assert.New(t)
    service := NewJWTService()

    token, err := service.GenerateToken("foobar", "read")
    assert.Nil(err)

    authToken, err := service.ValidateToken(token)
    assert.Nil(err)
    assert.Equal("foobar", authToken.Username)
    assert.Equal(domain.ReadAccess, authToken.Access)
    assert.True(authToken.ExpiresAt > time.Now().Unix())
}
