package domain

import (
    "testing"
    "github.com/stretchr/testify/assert"
)


func TestNewUser(t *testing.T) {

    assert := assert.New(t)
    var user *User
    var err error

    user, err = NewUser("foobar", "mypass123")
    assert.Nil(err)
    
    // Username should be equal to the given value
    assert.Equal("foobar", user.Username)
    // Password should not be an empty string or equal
    // to raw password
    assert.NotEqual("", user.Password)
    assert.NotEqual("mypass123", user.Password)

}

func TestUserPasswordVerification(t *testing.T) {

    assert := assert.New(t)
    var user *User
    var err error

    password := "$tr0ngP@$$w0rD"
    user, err = NewUser("foobar", password)

    assert.Nil(err)
    assert.Equal("foobar", user.Username)
    assert.NotEqual("", user.Password)

    // Incorrect password
    assert.False(user.VerifyPassword("lorem"))
    assert.False(user.VerifyPassword(""))

    // Correct password
    assert.True(user.VerifyPassword(password))
}
