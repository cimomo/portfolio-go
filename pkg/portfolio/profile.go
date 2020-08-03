package portfolio

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Profile contains multiple portfolios
type Profile struct {
	Name       string
	CostBasis  float64
	Market     *Market
	Portfolios []*Portfolio
	Status     *ProfileStatus
}

// ProfileStatus defines the real-time status of the entire portfolio
type ProfileStatus struct {
	Value                      float64
	RegularMarketChange        float64
	RegularMarketChangePercent float64
	Unrealized                 float64
	UnrealizedPercent          float64
	Allocation                 map[string]float64
}

type profileConfig []portfolioConfig

// NewProfile returns an empty profile
func NewProfile(name string) *Profile {
	return &Profile{
		Name:       name,
		Market:     NewMarket(),
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

	for _, portfolioConfig := range profileConfig {
		portfolio := NewPortfolio()
		portfolio.Load(portfolioConfig)
		profile.Portfolios = append(profile.Portfolios, portfolio)

		profile.CostBasis += portfolio.CostBasis
	}

	return nil
}

// Refresh computes the current status of the entire profile and its portfolios
func (profile *Profile) Refresh() error {
	for _, portfolio := range profile.Portfolios {
		err := portfolio.Refresh()
		if err != nil {
			return err
		}
	}

	profile.RefreshStatus()

	return nil
}

// RefreshStatus computes the current status of the entire profile
func (profile *Profile) RefreshStatus() {
	status := ProfileStatus{
		Allocation: make(map[string]float64),
	}

	for _, portfolio := range profile.Portfolios {
		status.Value += portfolio.Status.Value
		status.RegularMarketChange += portfolio.Status.RegularMarketChange
		status.Unrealized += portfolio.Status.Unrealized
	}

	previousValue := status.Value - status.RegularMarketChange

	status.RegularMarketChangePercent = 0
	if previousValue != 0 {
		status.RegularMarketChangePercent = (status.RegularMarketChange / previousValue) * 100
	}

	status.UnrealizedPercent = 0
	if profile.CostBasis != 0 {
		status.UnrealizedPercent = (status.Unrealized / profile.CostBasis) * 100
	}

	for _, portfolio := range profile.Portfolios {
		status.Allocation[portfolio.Name] = 0
		if status.Value != 0 {
			status.Allocation[portfolio.Name] = (portfolio.Status.Value / status.Value) * 100
		}
	}

	profile.Status = &status
}
