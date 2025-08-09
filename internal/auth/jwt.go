package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Claims represents our JWT claims including standard registered claims and custom fields.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Manager handles JWT generation and validation.
type Manager struct {
	Secret   []byte
	Issuer   string
	Audience string
	// Now provides current time; override in tests.
	Now func() time.Time
}

// NewManager creates a new JWT manager.
func NewManager(secret, issuer, audience string) *Manager {
	return &Manager{
		Secret:   []byte(secret),
		Issuer:   issuer,
		Audience: audience,
		Now:      time.Now,
	}
}

// GenerateToken signs and returns a JWT string for a given user and ttl.
func (m *Manager) GenerateToken(userID int64, username string, ttl time.Duration) (string, error) {
	now := m.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			Subject:   username,
			Audience:  jwt.ClaimStrings{m.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.Secret)
}

// ParseAndValidate parses the token string and validates signature and time-based claims.
func (m *Manager) ParseAndValidate(tokenStr string) (*Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.Secret, nil
	}

	parsedToken, err := jwt.ParseWithClaims(tokenStr, &Claims{}, keyFunc,
		jwt.WithIssuer(m.Issuer),
		jwt.WithAudience(m.Audience),
		jwt.WithLeeway(1*time.Minute), // small clock skew tolerance
	)
	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
