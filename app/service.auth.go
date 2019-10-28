package app

import (
    "errors"
    "fmt"
    "os"
    "time"
    "github.com/dgrijalva/jwt-go"
    "github.com/esquarer/go-rest/domain"
)

const (
    UsernameMinLength int = 3
    UsernameMaxLength int = 16
    PasswordMinLength int = 6
    PasswordMaxLength int = 32

    tokenExpiryDuration = time.Minute * 5
)

var (
    ErrUsernameLength = fmt.Errorf("username length should be between %v to %v characters", UsernameMinLength, UsernameMaxLength)
    ErrPasswordLength = fmt.Errorf("password length should be between %v to %v characters", PasswordMinLength, PasswordMaxLength)
    ErrInvalidToken = errors.New("invalid token")
    ErrIncorrectCredentials = errors.New("invalid user credentials")
)


// 
// 
type userClaims struct {
    jwt.StandardClaims
    Username string `json:"username"`
}


// 
type JWTService struct {
    repo    domain.UserRepository
    secret  []byte
}

func NewJWTService(repo domain.UserRepository) (*JWTService, error) {
    
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return nil, fmt.Errorf("empty value for JWT_SECRET")
    }

    return &JWTService{
        repo: repo,
        secret: []byte(secret),
    }, nil
}

// 
// 
func (service *JWTService) newClaims(user *domain.User) (jwt.Claims, int64) {

    now := time.Now()
    expiresAt := now.Add(tokenExpiryDuration)

    claims := &userClaims{}
    claims.Username = user.Username
    claims.StandardClaims.IssuedAt = now.Unix()
    claims.StandardClaims.ExpiresAt = expiresAt.Unix()

    return claims, claims.StandardClaims.ExpiresAt
}

// 
// 
func (service *JWTService) GenerateToken(user *domain.User) (string, int64, error) {

    claims, expiresAt := service.newClaims(user)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(service.secret)
    return signedToken, expiresAt, err
}

// 
// 
func (service *JWTService) ValidateToken(signedToken string) (*domain.User, error) {

    token, err := jwt.ParseWithClaims(signedToken, &userClaims{}, func(t *jwt.Token) (interface{}, error) {
        return service.secret, nil
    })

    if err != nil || !token.Valid {
        return nil, ErrInvalidToken
    }

    claims, ok := token.Claims.(*userClaims)
    if !ok { return nil, ErrInvalidToken }

    return service.repo.FindByUsername(claims.Username)
}


type userService struct {
    repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
    return &userService{repo}
}

func (service *userService) Register(username, password string) (*domain.User, error) {

    if len(username) < UsernameMinLength || len(username) > UsernameMaxLength {
        return nil, ErrUsernameLength
    }
    if len(password) < PasswordMinLength || len(password) > PasswordMaxLength  {
        return nil, ErrPasswordLength
    }

    user, err := domain.NewUser(username, password)
    if err != nil { return nil, err }

    err = service.repo.Add(user)
    if err != nil { return nil, err }

    return user, nil
}

func (service *userService) Authenticate(username, password string) error {
    user, err := service.repo.FindByUsername(username)
    if err != nil { return ErrIncorrectCredentials }

    if !user.VerifyPassword(password) {
        return ErrIncorrectCredentials
    }

    return nil
}
