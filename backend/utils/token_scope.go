package utils

import (
	"fmt"
	"strings"
)

var allowedTokenScopes = map[string]struct{}{
	"*":          {},
	"all":        {},
	"excel":      {},
	"dict":       {},
	"settlement": {},
	"meta":       {},
}

// NormalizeTokenScopes validates an optional comma-separated scope list.
// An empty list means the token can call every admin API. "all" is stored as "*".
func NormalizeTokenScopes(raw string) (string, error) {
	raw = strings.ReplaceAll(raw, "，", ",")
	parts := strings.Split(raw, ",")
	seen := map[string]struct{}{}
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		scope := strings.ToLower(strings.TrimSpace(part))
		if scope == "" {
			continue
		}
		if _, ok := allowedTokenScopes[scope]; !ok {
			return "", fmt.Errorf("unknown scope %q; allowed values: *, all, excel, dict, settlement, meta", scope)
		}
		if scope == "all" {
			scope = "*"
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	if _, ok := seen["*"]; ok {
		return "*", nil
	}
	return strings.Join(out, ","), nil
}

// TokenScopeAllows reports whether a stored scope list permits path.
// Empty scopes and "*" allow every admin path.
func TokenScopeAllows(scopes, path string) bool {
	scopes = strings.TrimSpace(scopes)
	if scopes == "" || scopes == "*" {
		return true
	}
	for _, part := range strings.Split(scopes, ",") {
		switch strings.TrimSpace(part) {
		case "", "*":
			return true
		case "excel":
			if strings.Contains(path, "/excel") {
				return true
			}
		case "dict":
			if strings.Contains(path, "/dict-manage") {
				return true
			}
		case "settlement":
			if strings.Contains(path, "/settlement-form") {
				return true
			}
		case "meta":
			if strings.HasSuffix(path, "/menus") || strings.Contains(path, "/options/") {
				return true
			}
		}
	}
	return false
}
