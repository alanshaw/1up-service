package router

import "github.com/alanshaw/ucantone/errors"

const CandidateUnavailableErrorName = "CandidateUnavailable"

func NewCandidateUnavailableError(msg string) error {
	return errors.New(CandidateUnavailableErrorName, msg)
}
