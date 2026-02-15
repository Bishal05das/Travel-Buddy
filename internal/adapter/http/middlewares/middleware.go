package middleeware

import "github.com/bishal05das/travelbuddy/config"

type Middleware struct {
	cnf *config.Config
}

func NeewMiddleware(cnf *config.Config) *Middleware {
	return &Middleware{
		cnf: cnf,
	}
}