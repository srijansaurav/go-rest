package domain

import (
    "errors"
    "strings"
    "strconv"
    "time"
)


var ErrInvalidAccessType = errors.New("invalid value for AccessType")
var ErrTokenExpired = errors.New("AuthToken expired")


type AccessType string

const (
    ReadAccess AccessType  = "READ"
    WriteAccess AccessType = "WRITE"
    AdminAccess AccessType = "ADMIN"
)

func NewAccessType(access string) (AccessType, error) {

    access = strings.ToUpper(access)
    accessType := AccessType(access)

    switch accessType {
    case ReadAccess, WriteAccess, AdminAccess:
        break
    default:
        return AccessType(""), ErrInvalidAccessType
    }

    return accessType, nil
}

func (access *AccessType) UnmarshalJSON(data []byte) error {

    unquoteData, err  := strconv.Unquote(string(data))
    if err != nil { return err }
    ac, err := NewAccessType(unquoteData)
    if err != nil { return err }
    *access = ac
    return nil
}


const tokenExpiryDuration = time.Minute * 5

type AuthToken struct {
    IssuedAt    int64 `json:"iat"`
    ExpiresAt   int64 `json:"exp"`
    Username    string `json:"preferred_username"`
    Access      AccessType `json:"access"`
}

func NewAuthToken(username string, access AccessType) *AuthToken {
    now := time.Now()
    expiry := now.Add(tokenExpiryDuration)
    return &AuthToken{
        IssuedAt: now.Unix(),
        ExpiresAt: expiry.Unix(),
        Username: username,
        Access: access,
    }
}

func (token *AuthToken) Refresh() error {
    if token.IsExpired() { return ErrTokenExpired }
    expiry := time.Now().Add(tokenExpiryDuration)
    token.ExpiresAt = expiry.Unix()
    return nil
}

func (token *AuthToken) IsExpired() bool {
    return token.ExpiresAt < time.Now().Unix()
}

func (token *AuthToken) HasWriteAccess() bool {
    return token.Access == WriteAccess || token.Access == AdminAccess
}

func (token *AuthToken) HasAdminAccess() bool {
    return token.Access == AdminAccess
}
