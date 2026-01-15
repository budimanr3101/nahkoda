package semantic

import (
	"errors"
	"regexp"
)

var (
	// RFC 1123 DNS Subdomain validation: lowercase alphanumeric, '-', can't start/end with '-'
	// We allow slightly more permissive for some specialized resources, but definitely no spaces or shell characters.
	// We also allow uppercase because Kubernetes resource names (like nodes) might sometimes be case specific?
	// Actually K8s names are lowercase usually. But finding by name is strict.
	// To prevent command injection, we primarily need to block ";", "|", "$", "`", etc.
	// Let's stick to safe chars: a-z, A-Z, 0-9, -, _, .
	validResourceName = regexp.MustCompile(`^[a-zA-Z0-9-._]+$`)
)

func ValidateResourceName(name string) error {
	if name == "" {
		return nil // Missing target handled by resolver logic
	}
	if !validResourceName.MatchString(name) {
		return errors.New("nama resource tidak valid (hanya boleh huruf, angka, -, _, .)")
	}
	return nil
}
