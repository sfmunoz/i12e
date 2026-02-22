package config

import (
	"fmt"
	"net/netip"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Docs refs:
//   https://pkg.go.dev/github.com/go-playground/validator/v10
//   ~/go/pkg/mod/github.com/go-playground/validator/v10@v10.30.1/baked_in.go
//   ~/go/pkg/mod/github.com/go-playground/validator/v10@v10.30.1/_examples/simple/main.go
// Home:
//   https://github.com/go-playground/validator

func validSemverV(fl validator.FieldLevel) bool {
	s1 := fl.Field().String()
	s2 := strings.TrimPrefix(s1, "v")
	return validator.New().Var(s2, "semver") == nil
}

func validButaneMode(fl validator.FieldLevel) bool {
	return slices.Contains(ValidModes(), fl.Field().String())
}

func validButaneOutput(fl validator.FieldLevel) bool {
	return slices.Contains(ValidOutputs(), fl.Field().String())
}

func validMeshNetwork(fl validator.FieldLevel) bool {
	p, ok := fl.Field().Addr().Interface().(*netip.Prefix)
	if !ok {
		return false
	}
	meshAddr := p.Addr()
	if !meshAddr.Is4() {
		fmt.Printf("config: 'mesh.network_address=%s' is not IPv4\n", p)
		return false
	}
	if !meshAddr.IsPrivate() {
		fmt.Printf("config: 'mesh.network_address=%s' is not private\n", p)
		return false
	}
	b := p.Bits() // from /12 (20 bits for host) to /29 (3 bits for host)
	if b < 12 {
		fmt.Printf("config: wrong 'mesh.network_address=%s' (bits=%d, min=12)\n", p, b)
		return false
	}
	if b > 29 {
		fmt.Printf("config: wrong 'mesh.network_address=%s' (bits=%d, max=29)\n", p, b)
		return false
	}
	return true
}
