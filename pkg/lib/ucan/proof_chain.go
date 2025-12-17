package ucanlib

import (
	"github.com/alanshaw/1up-service/pkg/store/delegation"
	"github.com/alanshaw/ucantone/ucan"
)

func ConstructProofChain(store delegation.Store, aud ucan.Principal, cmd ucan.Command, sub ucan.Principal) ([]ucan.Delegation, []ucan.Link, error) {
	proofs := []ucan.Delegation{}
	links := []ucan.Link{}

	matches, err := store.Query(aud, cmd, sub)
	if err != nil {
		return nil, nil, err
	}

	for _, d := range matches {
		if d.Subject() != nil && d.Subject().DID() == d.Issuer().DID() {
			proofs = append(proofs, d)
			links = append(links, d.Link())
			break
		}
		// if subject is nil, or subject != issuer, we need more proof
		ps, ls, err := ConstructProofChain(store, d.Issuer(), d.Command(), sub)
		if err != nil {
			return nil, nil, err
		}
		if len(ps) == 0 {
			continue // try a different path
		}
		proofs = append(proofs, d)
		proofs = append(proofs, ps...)
		links = append(links, d.Link())
		links = append(links, ls...)
		break
	}

	return proofs, links, nil
}
