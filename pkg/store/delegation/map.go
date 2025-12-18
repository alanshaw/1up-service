package delegation

import (
	"context"
	"iter"

	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/command"
)

const nullSubject = "null"

type Audience = did.DID
type Subject = string

type MapDelegationStore struct {
	data map[Audience]map[ucan.Command]map[Subject][]ucan.Delegation
}

var _ Store = (*MapDelegationStore)(nil)

func NewMapDelegationStore(delegations []ucan.Delegation) *MapDelegationStore {
	data := map[Audience]map[ucan.Command]map[Subject][]ucan.Delegation{}
	for _, d := range delegations {
		aud := d.Audience().DID()
		cmd := d.Command()
		var sub Subject
		if d.Subject() == nil {
			sub = nullSubject // powerline delegation
		} else {
			sub = d.Subject().DID().String()
		}
		if _, ok := data[aud]; !ok {
			data[aud] = map[ucan.Command]map[Subject][]ucan.Delegation{}
		}
		if _, ok := data[aud][cmd]; !ok {
			data[aud][cmd] = map[Subject][]ucan.Delegation{}
		}
		if _, ok := data[aud][cmd][sub]; !ok {
			data[aud][cmd][sub] = []ucan.Delegation{}
		}
		data[aud][cmd][sub] = append(data[aud][cmd][sub], d)
	}
	return &MapDelegationStore{data}
}

func (m *MapDelegationStore) Query(ctx context.Context, aud ucan.Principal, cmd ucan.Command, sub ucan.Subject) iter.Seq2[ucan.Delegation, error] {
	return func(yield func(ucan.Delegation, error) bool) {
		cmdDelegations, ok := m.data[aud.DID()]
		if !ok {
			return
		}

		segs := cmd.Segments()
		for i := len(segs) - 1; i >= 0; i-- {
			cmd := command.Top().Join(segs[0 : i+1]...)
			subDelegations, ok := cmdDelegations[cmd]
			if !ok {
				return
			}
			dlgs, ok := subDelegations[sub.DID().String()]
			if !ok {
				dlgs, ok = subDelegations[nullSubject]
				if !ok {
					return
				}
			} else {
				powerlineDlgs, ok := subDelegations[nullSubject]
				if ok {
					dlgs = append(dlgs, powerlineDlgs...)
				}
			}
			for _, d := range dlgs {
				if !yield(d, nil) {
					return
				}
			}
		}
	}
}
