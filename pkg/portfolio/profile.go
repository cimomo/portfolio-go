package portfolio

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Profile contains multiple portfolios
type Profile struct {
	Name       string
	Portfolios []*Portfolio
}

type profileConfig []portfolioConfig

// NewProfile returns an empty profile
func NewProfile(name string) *Profile {
	return &Profile{
		Name:       name,
		Portfolios: make([]*Portfolio, 0),
	}
}

// Load loads a profile from the given file
func (profile *Profile) Load(name string) error {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	profileConfig := profileConfig{}

	err = yaml.Unmarshal(file, &profileConfig)
	if err != nil {
		return err
	}

	portfolioConfig := profileConfig[0]

	portfolio := NewPortfolio()

	portfolio.Load(portfolioConfig)

	profile.Portfolios = append(profile.Portfolios, portfolio)

	return nil
}
