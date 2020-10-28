package specification

type Ignorer struct {
	pm *PatternMatcher
}

func NewIgnorer(patterns []string) (*Ignorer, error) {
	pm, err := NewPatternMatcher(patterns)
	if err != nil {
		return nil, err
	}

	return &Ignorer{pm: pm}, nil
}

func (i *Ignorer) Ignore(s string) bool {
	return i.pm.Match(s)
}
