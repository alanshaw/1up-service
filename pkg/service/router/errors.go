package router

import "github.com/alanshaw/ucantone/errors"

const CandidateUnavailableErrorName = "CandidateUnavailable"

var ErrCandidateUnavailable = errors.New(CandidateUnavailableErrorName, "no storage providers available")
