package security

import (
	"net/url"
	"os"
	"strings"
)

const (
	allowedOriginsEnv = "ALLOWED_ORIGINS"
	allowedOriginEnv  = "ALLOWED_ORIGIN"
)

type OriginPolicy struct {
	allowed map[string]struct{}
}

func NewOriginPolicyFromEnv() OriginPolicy {
	rawOrigins := os.Getenv(allowedOriginsEnv)
	if strings.TrimSpace(rawOrigins) == "" {
		rawOrigins = os.Getenv(allowedOriginEnv)
	}

	parts := strings.Split(rawOrigins, ",")
	allowed := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		origin, ok := normalizeOrigin(part)
		if !ok {
			continue
		}
		allowed[origin] = struct{}{}
	}

	return OriginPolicy{
		allowed: allowed,
	}
}

func (p OriginPolicy) IsAllowed(origin string) bool {
	normalized, ok := normalizeOrigin(origin)
	if !ok {
		return false
	}

	_, ok = p.allowed[normalized]
	return ok
}

func (p OriginPolicy) IsAllowedRequest(origin, requestHost string) bool {
	normalizedOrigin, ok := normalizeOrigin(origin)
	if !ok {
		return false
	}

	if p.HasAllowedOrigins() {
		_, ok := p.allowed[normalizedOrigin]
		return ok
	}

	host := strings.ToLower(strings.TrimSpace(requestHost))
	if host == "" {
		return false
	}

	parsedOrigin, err := url.Parse(normalizedOrigin)
	if err != nil {
		return false
	}

	return strings.EqualFold(parsedOrigin.Host, host)
}

func (p OriginPolicy) HasAllowedOrigins() bool {
	return len(p.allowed) > 0
}

func normalizeOrigin(origin string) (string, bool) {
	value := strings.TrimSpace(origin)
	if value == "" {
		return "", false
	}

	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", false
	}

	return strings.ToLower(parsed.Scheme) + "://" + strings.ToLower(parsed.Host), true
}
