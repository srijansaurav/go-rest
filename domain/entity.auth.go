package domain

import (
    "golang.org/x/crypto/bcrypt"
)


// 
type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// 
func NewUser(username, password string) (*User, error) {

    user := &User{
        Username: username,
    }
    err := user.setPassword(password)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// 
func (u *User) setPassword(rawPassword string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
    if err == nil {
        u.Password = string(hash)
    }
    return err
}

// 
func (u *User) VerifyPassword(inputPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
    return err == nil
}


// 
// 
type UserRepository interface {
    
    // 
    Add(*User) error
    
    // 
    FindByUsername(string) (*User, error)
}


// 
// 
type UserService interface {
    
    // 
    Register(string, string) (*User, error)

    // 
    Authenticate(string, string) (error)
}



type TokenService interface {

    // 
    GenerateToken(*User) (string, int64, error)

    // 
    ValidateToken(string) (*User, error)
}
