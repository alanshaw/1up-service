package routing

import "github.com/alanshaw/ucantone/errors"

const CandidateUnavailableErrorName = "CandidateUnavailable"

var ErrCandidateUnavailable = errors.New(CandidateUnavailableErrorName, "no storage providers available")

var ErrNotFound = errors.New("NotFound", "storage provider not found")
