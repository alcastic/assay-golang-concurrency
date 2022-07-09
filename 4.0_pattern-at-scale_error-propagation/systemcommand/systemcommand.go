package systemcommand

import (
	"github.com/alcastic/assay-golang-concurrency/4.0_pattern-at-scale_error-propagation/myerror"
	"os"
)

type LowLevelError struct {
	error
}

func IsExecutableByTheOwner(executable string) (bool, error) {
	fi, err := os.Stat(executable)
	if err != nil {
		return false, LowLevelError{myerror.WrapError(err, "executable not found")}
	}
	return fi.Mode().Perm()&0100 == 0100, nil
}
