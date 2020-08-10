package portfolio

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Profile contains multiple portfolios
type Profile struct {
	Name             string
	Cash             float64
	CostBasis        float64
	Market           *Market
	Portfolios       []*Portfolio
	TargetAllocation map[string]float64
	MergedPortfolio  *Portfolio
	Status           *ProfileStatus
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

type profileConfig struct {
	Cash       cashConfig        `yaml:"cash"`
	Portfolios []portfolioConfig `yaml:"portfolios"`
}

type cashConfig struct {
	Value            float64 `yaml:"value"`
	TargetAllocation float64 `yaml:"allocation"`
}

// NewProfile returns an empty profile
func NewProfile(name string) *Profile {
	return &Profile{
		Name:             name,
		Market:           NewMarket(),
		Portfolios:       make([]*Portfolio, 0),
		TargetAllocation: make(map[string]float64),
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

	profile.Cash = profileConfig.Cash.Value
	profile.CostBasis = profileConfig.Cash.Value
	profile.TargetAllocation["cash"] = profileConfig.Cash.TargetAllocation
	totalAllocation := profileConfig.Cash.TargetAllocation

	for _, portfolioConfig := range profileConfig.Portfolios {
		portfolio := NewPortfolio()
		portfolio.Load(portfolioConfig)
		profile.Portfolios = append(profile.Portfolios, portfolio)
		profile.TargetAllocation[portfolio.Name] = portfolioConfig.TargetAllocation
		totalAllocation += portfolioConfig.TargetAllocation
		profile.CostBasis += portfolio.CostBasis
	}

	if totalAllocation != 100.0 && totalAllocation != 0.0 {
		return errors.New("Total allocation should be either 0% (ignored) or 100%")
	}

	profile.MergedPortfolio = profile.mergePortfolios()

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

	err := profile.MergedPortfolio.Refresh()
	if err != nil {
		return err
	}

	profile.RefreshStatus()

	return nil
}

// RefreshStatus computes the current status of the entire profile
func (profile *Profile) RefreshStatus() {
	status := ProfileStatus{
		Allocation: make(map[string]float64),
	}

	status.Value = profile.Cash

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

	status.Allocation["cash"] = (profile.Cash / status.Value) * 100

	profile.Status = &status
}

// MergePortfolios merges all portfolios in the profile into a single portfolio
func (profile *Profile) mergePortfolios() *Portfolio {
	portfolio := NewPortfolio()

	portfolio.Name = profile.Name

	for _, port := range profile.Portfolios {
		portfolio.CostBasis += port.CostBasis
		portfolio.Symbols = append(portfolio.Symbols, port.Symbols...)

		for symbol, holding := range port.Holdings {
			portfolio.Holdings[symbol] = holding
		}
	}

	return portfolio
}
