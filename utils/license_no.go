//go:build !license

package utils

type License interface {
	GenMachineID() string
	Activate(code string) bool
}

type NoLicense struct {
}

func (*NoLicense) GenMachineID() string {
	return "******"
}

func (*NoLicense) Activate(code string) bool {
	return true
}
