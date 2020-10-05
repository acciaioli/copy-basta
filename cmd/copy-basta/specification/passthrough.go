package specification

type PassThrough []string

func (p PassThrough) validate() error {
	return nil
}

type Passer struct {
	pm *PatternMatcher
}

func NewPasser(passThrough PassThrough) (*Passer, error) {
	pm, err := NewPatternMatcher(passThrough)
	if err != nil {
		return nil, err
	}

	return &Passer{pm: pm}, nil
}

func (i *Passer) Pass(s string) bool {
	return i.pm.Match(s)
}
