package delegation

import "github.com/alanshaw/ucantone/ucan"

type Store interface {
	// Query finds delegations matching the given audience, command, and subject.
	// Note: subject MUST not be nil. Matching delegations MAY include powerline
	// delegations (with nil subject) and delegations where command is a matching
	// parent of the passed command.
	Query(aud ucan.Principal, cmd ucan.Command, sub ucan.Subject) ([]ucan.Delegation, error)
}
