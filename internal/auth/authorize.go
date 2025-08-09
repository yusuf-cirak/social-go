package auth

// A simple policy engine with allow rules and deny-by-default semantics.

type Subject struct {
	UserID int64
	Roles  []string
}

type Resource struct {
	Type    string
	OwnerID int64
	Attr    map[string]any
}

type Rule func(Subject, string, Resource) bool

type PolicyEngine struct {
	rules []Rule
}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{rules: make([]Rule, 0, 8)}
}

func (e *PolicyEngine) Allow(r Rule) *PolicyEngine {
	e.rules = append(e.rules, r)
	return e
}

// Authorize evaluates rules and returns true if any rule allows the action.
func (e *PolicyEngine) Authorize(sub Subject, action string, res Resource) bool {
	for _, rule := range e.rules {
		if rule(sub, action, res) {
			return true
		}
	}
	return false
}

// Common action constants
const (
	ActionPostCreate   = "post:create"
	ActionPostUpdate   = "post:update"
	ActionPostDelete   = "post:delete"
	ActionUserFollow   = "user:follow"
	ActionUserUnfollow = "user:unfollow"
)

// NewDefaultPolicyEngine returns an engine with common default rules.
func NewDefaultPolicyEngine() *PolicyEngine {
	e := NewPolicyEngine()

	// Anyone authenticated can create posts
	e.Allow(func(s Subject, action string, r Resource) bool {
		if action == ActionPostCreate && r.Type == "post" {
			return s.UserID != 0
		}
		return false
	})

	// Only owners can update/delete their posts
	e.Allow(func(s Subject, action string, r Resource) bool {
		if r.Type != "post" {
			return false
		}
		if action == ActionPostUpdate || action == ActionPostDelete {
			return s.UserID != 0 && s.UserID == r.OwnerID
		}
		return false
	})

	// A user can follow/unfollow others, but not themselves
	e.Allow(func(s Subject, action string, r Resource) bool {
		if r.Type != "user" {
			return false
		}
		if action == ActionUserFollow || action == ActionUserUnfollow {
			return s.UserID != 0 && s.UserID != r.OwnerID
		}
		return false
	})

	return e
}
