package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewManager("test-secret", "test-issuer", "test-audience")

	userID := int64(123)
	username := "testuser"
	ttl := time.Hour

	token, err := manager.GenerateToken(userID, username, ttl)

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Token should have 3 parts (header.payload.signature)
	assert.Equal(t, 3, len(strings.Split(token, ".")))
}

func TestJWTManager_ParseAndValidate_Success(t *testing.T) {
	manager := NewManager("test-secret", "test-issuer", "test-audience")

	userID := int64(123)
	username := "testuser"

	// Generate a valid token
	token, err := manager.GenerateToken(userID, username, time.Hour)
	require.NoError(t, err)

	// Parse and validate it
	claims, err := manager.ParseAndValidate(token)

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.Equal(t, "test-audience", claims.Audience[0])
}

func TestJWTManager_ParseAndValidate_InvalidToken(t *testing.T) {
	manager := NewManager("test-secret", "test-issuer", "test-audience")

	testCases := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"malformed token", "invalid.token"},
		{"random string", "this-is-not-a-jwt-token"},
		{"wrong format", "header.payload"}, // Missing signature
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := manager.ParseAndValidate(tc.token)
			assert.Error(t, err)
		})
	}
}

func TestJWTManager_ParseAndValidate_ExpiredToken(t *testing.T) {
	manager := NewManager("test-secret", "test-issuer", "test-audience")

	// Generate token that's already expired
	token, err := manager.GenerateToken(123, "testuser", -time.Hour)
	require.NoError(t, err)

	// Should fail validation due to expiration
	_, err = manager.ParseAndValidate(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWTManager_ParseAndValidate_WrongSecret(t *testing.T) {
	manager1 := NewManager("secret1", "test-issuer", "test-audience")
	manager2 := NewManager("secret2", "test-issuer", "test-audience")

	// Generate token with first manager
	token, err := manager1.GenerateToken(123, "testuser", time.Hour)
	require.NoError(t, err)

	// Try to validate with second manager (different secret)
	_, err = manager2.ParseAndValidate(token)
	assert.Error(t, err)
}

func TestJWTManager_ParseAndValidate_WrongIssuer(t *testing.T) {
	manager1 := NewManager("test-secret", "issuer1", "test-audience")
	manager2 := NewManager("test-secret", "issuer2", "test-audience")

	// Generate token with first manager
	token, err := manager1.GenerateToken(123, "testuser", time.Hour)
	require.NoError(t, err)

	// Try to validate with second manager (different issuer)
	_, err = manager2.ParseAndValidate(token)
	assert.Error(t, err)
}

func TestJWTManager_ParseAndValidate_WrongAudience(t *testing.T) {
	manager1 := NewManager("test-secret", "test-issuer", "audience1")
	manager2 := NewManager("test-secret", "test-issuer", "audience2")

	// Generate token with first manager
	token, err := manager1.GenerateToken(123, "testuser", time.Hour)
	require.NoError(t, err)

	// Try to validate with second manager (different audience)
	_, err = manager2.ParseAndValidate(token)
	assert.Error(t, err)
}

func TestJWTManager_CustomTimeProvider(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	manager := NewManager("test-secret", "test-issuer", "test-audience")
	manager.Now = func() time.Time { return fixedTime }

	token, err := manager.GenerateToken(123, "testuser", time.Hour)
	require.NoError(t, err)

	// Also need to set the Now function for validation to use the same time
	manager.Now = func() time.Time { return fixedTime }

	claims, err := manager.ParseAndValidate(token)
	require.NoError(t, err)

	// Check that issued at time matches our fixed time
	assert.Equal(t, fixedTime.Unix(), claims.IssuedAt.Unix())
	assert.Equal(t, fixedTime.Add(time.Hour).Unix(), claims.ExpiresAt.Unix())
}

func TestPolicyEngine_Allow(t *testing.T) {
	engine := NewPolicyEngine()

	// Add a simple rule
	rule := func(s Subject, action string, r Resource) bool {
		return s.UserID == 1 && action == "test:action"
	}

	engine.Allow(rule)

	// Test the rule
	subject := Subject{UserID: 1}
	resource := Resource{}

	assert.True(t, engine.Authorize(subject, "test:action", resource))
	assert.False(t, engine.Authorize(subject, "other:action", resource))

	// Different user
	subject2 := Subject{UserID: 2}
	assert.False(t, engine.Authorize(subject2, "test:action", resource))
}

func TestPolicyEngine_DefaultRules_PostActions(t *testing.T) {
	engine := NewDefaultPolicyEngine()

	user := Subject{UserID: 1}
	ownPost := Resource{Type: "post", OwnerID: 1}
	otherPost := Resource{Type: "post", OwnerID: 2}
	newPost := Resource{Type: "post"}

	// Test post creation - any authenticated user can create
	assert.True(t, engine.Authorize(user, ActionPostCreate, newPost))

	// Test post update/delete - only owner can modify
	assert.True(t, engine.Authorize(user, ActionPostUpdate, ownPost))
	assert.True(t, engine.Authorize(user, ActionPostDelete, ownPost))

	assert.False(t, engine.Authorize(user, ActionPostUpdate, otherPost))
	assert.False(t, engine.Authorize(user, ActionPostDelete, otherPost))

	// Unauthenticated user (UserID = 0) cannot do anything
	unauthUser := Subject{UserID: 0}
	assert.False(t, engine.Authorize(unauthUser, ActionPostCreate, newPost))
	assert.False(t, engine.Authorize(unauthUser, ActionPostUpdate, ownPost))
	assert.False(t, engine.Authorize(unauthUser, ActionPostDelete, ownPost))
}

func TestPolicyEngine_DefaultRules_UserActions(t *testing.T) {
	engine := NewDefaultPolicyEngine()

	user := Subject{UserID: 1}
	otherUser := Resource{Type: "user", OwnerID: 2}
	selfUser := Resource{Type: "user", OwnerID: 1}

	// User can follow/unfollow others
	assert.True(t, engine.Authorize(user, ActionUserFollow, otherUser))
	assert.True(t, engine.Authorize(user, ActionUserUnfollow, otherUser))

	// User cannot follow/unfollow themselves
	assert.False(t, engine.Authorize(user, ActionUserFollow, selfUser))
	assert.False(t, engine.Authorize(user, ActionUserUnfollow, selfUser))

	// Unauthenticated user cannot follow/unfollow
	unauthUser := Subject{UserID: 0}
	assert.False(t, engine.Authorize(unauthUser, ActionUserFollow, otherUser))
	assert.False(t, engine.Authorize(unauthUser, ActionUserUnfollow, otherUser))
}

func TestPolicyEngine_MultipleRules(t *testing.T) {
	engine := NewPolicyEngine()

	// Rule 1: Admins can do anything
	adminRule := func(s Subject, action string, r Resource) bool {
		for _, role := range s.Roles {
			if role == "admin" {
				return true
			}
		}
		return false
	}

	// Rule 2: Users can read their own data
	userRule := func(s Subject, action string, r Resource) bool {
		return action == "read" && s.UserID == r.OwnerID
	}

	engine.Allow(adminRule).Allow(userRule)

	// Test admin (should be allowed to do anything)
	admin := Subject{UserID: 1, Roles: []string{"admin"}}
	resource := Resource{Type: "user", OwnerID: 2}

	assert.True(t, engine.Authorize(admin, "delete", resource))
	assert.True(t, engine.Authorize(admin, "read", resource))

	// Test regular user (can only read own data)
	user := Subject{UserID: 2, Roles: []string{"user"}}
	ownResource := Resource{Type: "user", OwnerID: 2}
	otherResource := Resource{Type: "user", OwnerID: 3}

	assert.True(t, engine.Authorize(user, "read", ownResource))
	assert.False(t, engine.Authorize(user, "read", otherResource))
	assert.False(t, engine.Authorize(user, "delete", ownResource))
}

func TestPolicyEngine_NoRules(t *testing.T) {
	engine := NewPolicyEngine()

	subject := Subject{UserID: 1}
	resource := Resource{Type: "test"}

	// With no rules, everything should be denied
	assert.False(t, engine.Authorize(subject, "any-action", resource))
}

func TestPolicyEngine_ResourceAttributes(t *testing.T) {
	engine := NewPolicyEngine()

	// Rule that checks resource attributes
	rule := func(s Subject, action string, r Resource) bool {
		if action == "read" {
			// Allow reading public resources or own resources
			if public, ok := r.Attr["public"].(bool); ok && public {
				return true
			}
			return s.UserID == r.OwnerID
		}
		return false
	}

	engine.Allow(rule)

	user := Subject{UserID: 1}

	// Public resource - anyone can read
	publicResource := Resource{
		Type:    "post",
		OwnerID: 2,
		Attr:    map[string]any{"public": true},
	}
	assert.True(t, engine.Authorize(user, "read", publicResource))

	// Private resource owned by others - cannot read
	privateResource := Resource{
		Type:    "post",
		OwnerID: 2,
		Attr:    map[string]any{"public": false},
	}
	assert.False(t, engine.Authorize(user, "read", privateResource))

	// Own resource - can read regardless of public flag
	ownResource := Resource{
		Type:    "post",
		OwnerID: 1,
		Attr:    map[string]any{"public": false},
	}
	assert.True(t, engine.Authorize(user, "read", ownResource))
}
