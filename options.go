package server

type Option func(*server)

func WithDebug(dbg bool) Option {
	return func(s *server) {
		s.debug = dbg
	}
}
