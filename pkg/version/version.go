// Package version specifies the version informations
package version

import "fmt"

const (
	major = 2
	minor = 0
	patch = 0
	name  = "G.L.I.F."
)

// Get returns the formatted string containing the version informations
func Get() string {
	return fmt.Sprintf("%s: %d.%d.%d", name, major, minor, patch)
}
