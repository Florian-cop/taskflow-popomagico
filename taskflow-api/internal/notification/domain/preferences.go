package domain

type Preferences struct {
	UserID  string
	Enabled map[string]bool
}

func DefaultPreferences(userID string) *Preferences {
	return &Preferences{
		UserID: userID,
		Enabled: map[string]bool{
			"email":  true,
			"in_app": true,
		},
	}
}

func (p *Preferences) IsEnabled(channel string) bool {
	if p == nil || p.Enabled == nil {
		return true
	}
	v, ok := p.Enabled[channel]
	return !ok || v
}

func (p *Preferences) Set(channel string, enabled bool) {
	if p.Enabled == nil {
		p.Enabled = map[string]bool{}
	}
	p.Enabled[channel] = enabled
}
