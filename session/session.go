package session

type Session struct {
	id string
}

func (s Session) ID() string {
	return s.id
}

func (s *Session) SetID(id string) {
	s.id = id
}
