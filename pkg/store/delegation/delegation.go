package delegation

import ucanlib "github.com/alanshaw/libracha/ucan"

type Store interface {
	ucanlib.DelegationFinder
}
