package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"
)

// PersonalAccessTokenPrefix marks a bearer credential as a personal access token
// so auth can tell it apart from a session JWT.
const PersonalAccessTokenPrefix = "ssc_pat_"

// PersonalAccessToken is a long-lived API credential. Only the SHA-256 hex digest
// is stored. The plaintext is returned once by the create API and never again.
type PersonalAccessToken struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint       `gorm:"index;not null" json:"userId"`
	Name       string     `gorm:"type:varchar(100);not null" json:"name"`
	TokenHash  string     `gorm:"type:char(64);uniqueIndex;not null" json:"-"`
	Prefix     string     `gorm:"type:varchar(24);not null" json:"prefix"`
	Scopes     string     `gorm:"type:varchar(500);not null;default:''" json:"scopes"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	ExpiresAt  *time.Time `json:"expiresAt"`
	RevokedAt  *time.Time `json:"revokedAt"`
}

func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
}

func IsPersonalAccessToken(token string) bool {
	return len(token) > len(PersonalAccessTokenPrefix) && strings.HasPrefix(token, PersonalAccessTokenPrefix)
}

func HashPersonalAccessToken(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}

// NewPersonalAccessTokenSecret returns the one-time plaintext, its SHA-256 hex
// hash, and the short prefix stored for display.
func NewPersonalAccessTokenSecret() (plaintext string, hash string, prefix string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", "", err
	}
	plaintext = PersonalAccessTokenPrefix + base64.RawURLEncoding.EncodeToString(buf)
	hash = HashPersonalAccessToken(plaintext)
	prefix = plaintext
	if len(prefix) > 12 {
		prefix = plaintext[:12]
	}
	return plaintext, hash, prefix, nil
}
