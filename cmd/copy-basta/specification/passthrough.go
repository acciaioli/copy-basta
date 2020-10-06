package specification

type Passer struct {
	pm *PatternMatcher
}

func NewPasser(passThrough []string) (*Passer, error) {
	pm, err := NewPatternMatcher(passThrough)
	if err != nil {
		return nil, err
	}

	return &Passer{pm: pm}, nil
}

func (i *Passer) Pass(s string) bool {
	return i.pm.Match(s)
}
