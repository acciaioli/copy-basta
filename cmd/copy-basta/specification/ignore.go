package specification

type Ignore []string

func (i Ignore) validate() error {
	return nil
}

type Ignorer struct {
	pm *PatternMatcher
}

func NewIgnorer(ignore Ignore) (*Ignorer, error) {
	pm, err := NewPatternMatcher(ignore)
	if err != nil {
		return nil, err
	}

	return &Ignorer{pm: pm}, nil
}

func (i *Ignorer) Ignore(s string) bool {
	return i.pm.Match(s)
}
