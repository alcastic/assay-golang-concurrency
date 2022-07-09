package service

import (
	"github.com/alcastic/assay-golang-concurrency/4.0_pattern-at-scale_error-propagation/myerror"
	"github.com/alcastic/assay-golang-concurrency/4.0_pattern-at-scale_error-propagation/systemcommand"
	"os/exec"
)

type ServiceError struct {
	error
}

func Echo(msg string) error {
	echoPath := "/usr/bin/echo"
	isExecutable, err := systemcommand.IsExecutableByTheOwner(echoPath)
	if err != nil {
		return ServiceError{myerror.WrapError(err, "No executable available")}
	}
	if !isExecutable {
		return ServiceError{myerror.WrapError(nil, "Missed owner executable permissions")}
	}
	return exec.Command(echoPath, msg).Run() // Run() can fail with an error (not mapped as an example for an un treated error)
}
