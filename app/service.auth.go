package app

import (
    "os"
    "encoding/json"
    "github.com/dgrijalva/jwt-go"
    "github.com/esquarer/go-rest/domain"
)


type AuthTokenService interface {

    GenerateToken(string, string) (string, error)
    ValidateToken(string) (*domain.AuthToken, error)
    RefreshToken(string) (string, error)
}


type JWTService struct {}

func NewJWTService() *JWTService {
    return &JWTService{}
}

func (service *JWTService) getSecretKey() []byte {
    key := os.Getenv("JWT_SECRET_KEY")
    if key == "" { key = "foobar" }
    return []byte(key)
}

func (service *JWTService) claim(authToken *domain.AuthToken) jwt.Claims {
    var claim jwt.MapClaims
    j, _ := json.Marshal(authToken)
    _ = json.Unmarshal(j, &claim)
    return claim
}

func (service *JWTService) generateSignedToken(claim jwt.Claims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    return token.SignedString(service.getSecretKey())
}

func (service *JWTService) GenerateToken(username, access string) (string, error) {

    accessType, err := domain.NewAccessType(access)
    if err != nil { return "", err }

    authToken := domain.NewAuthToken(username, accessType)
    claim := service.claim(authToken)
    return service.generateSignedToken(claim)
}

func (service *JWTService) ValidateToken(token string) (*domain.AuthToken, error) {

    var claim jwt.MapClaims
    var authToken *domain.AuthToken

    _, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
        return service.getSecretKey(), nil
    })

    if err != nil { return nil, err }

    j, _ := json.Marshal(claim)
    _ = json.Unmarshal(j, &authToken)
    return authToken, nil
}

func (service *JWTService) RefreshToken(token string) (string, error) {

    var refreshedToken string

    authToken, err := service.ValidateToken(token)
    if err != nil { return refreshedToken, err }

    err = authToken.Refresh()
    if err != nil { return refreshedToken, err }

    claim := service.claim(authToken)
    return service.generateSignedToken(claim)
}
