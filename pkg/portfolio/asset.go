package portfolio

// AssetClass defines the type for asset classes and subclasses
type AssetClass string

// Asset class definitions
const (
	AssetClassUSStock                      AssetClass = "US Stock"
	AssetClassUSStockLarge                 AssetClass = "US Stock Large"
	AssetClassUSStockLargeValue            AssetClass = "US Stock Large Value"
	AssetClassUSStockLargeGrowth           AssetClass = "US Stock Large Growth"
	AssetClassUSStockLargeTech             AssetClass = "US Stock Large Tech"
	AssetClassUSStockMid                   AssetClass = "US Stock Mid"
	AssetClassUSStockSmall                 AssetClass = "US Stock Small"
	AssetClassInternationalStock           AssetClass = "International Stock"
	AssetClassEmergentMarketStock          AssetClass = "EmergentMarket Stock"
	AssetClassChinaStock                   AssetClass = "China Stock"
	AssetClassUSRealEstate                 AssetClass = "US Real Estate"
	AssetClassInternationalRealEstate      AssetClass = "International Real Estate"
	AssetClassUSBond                       AssetClass = "US Bond"
	AssetClassUSTreasury                   AssetClass = "US Treasury"
	AssetClassUSTreasuryLongTerm           AssetClass = "US Treasury Long Term"
	AssetClassUSTreasuryIntermediateTerm   AssetClass = "US Treasury Intermediate Term"
	AssetClassUSTreasuryShortTerm          AssetClass = "US Treasury Short Term"
	AssetClassUSTreasuryInflationProtected AssetClass = "US Treasury Inflation Protected"
	AssetClassCommodity                    AssetClass = "Commodity"
	AssetClassCrudeOil                     AssetClass = "CrudeOil"
	AssetClassGold                         AssetClass = "Gold"
	AssetClassSilver                       AssetClass = "Silver"
	AssetClassOther                        AssetClass = "Other"
)

// Asset defines a security held in the portfolio
type Asset struct {
	Symbol   string
	Class    AssetClass
	Subclass AssetClass
}

// AssetDB returns a built-in database for assets of interest
var AssetDB = func() map[string]*Asset {
	return map[string]*Asset{
		"VTI":  {"VTI", AssetClassUSStock, AssetClassUSStockLarge},
		"VTV":  {"VTV", AssetClassUSStock, AssetClassUSStockLargeValue},
		"VUG":  {"VUG", AssetClassUSStock, AssetClassUSStockLargeGrowth},
		"QQQ":  {"QQQ", AssetClassUSStock, AssetClassUSStockLargeTech},
		"VO":   {"VO", AssetClassUSStock, AssetClassUSStockMid},
		"VB":   {"VB", AssetClassUSStock, AssetClassUSStockSmall},
		"VXUS": {"VXUS", AssetClassInternationalStock, AssetClassInternationalStock},
		"VWO":  {"VWO", AssetClassInternationalStock, AssetClassEmergentMarketStock},
		"GXC":  {"GXC", AssetClassChinaStock, AssetClassChinaStock},
		"VNQ":  {"VNQ", AssetClassUSRealEstate, AssetClassUSRealEstate},
		"VNQI": {"VNQI", AssetClassInternationalRealEstate, AssetClassInternationalRealEstate},
		"BND":  {"BND", AssetClassUSBond, AssetClassUSBond},
		"GOVT": {"GOVT", AssetClassUSTreasury, AssetClassUSTreasury},
		"VGLT": {"VGLT", AssetClassUSTreasury, AssetClassUSTreasuryLongTerm},
		"SPTI": {"SPTI", AssetClassUSTreasury, AssetClassUSTreasuryIntermediateTerm},
		"SHY":  {"SHY", AssetClassUSTreasury, AssetClassUSTreasuryShortTerm},
		"TIP":  {"TIP", AssetClassUSTreasury, AssetClassUSTreasuryInflationProtected},
		"DBC":  {"DBC", AssetClassCommodity, AssetClassCommodity},
		"USO":  {"USO", AssetClassCommodity, AssetClassCrudeOil},
		"IAU":  {"IAU", AssetClassCommodity, AssetClassGold},
		"SLV":  {"SLV", AssetClassCommodity, AssetClassSilver},
	}
}

// NewAsset returns a new asset
func NewAsset(symbol string) *Asset {
	assets := AssetDB()

	asset := assets[symbol]
	if asset == nil {
		return defaultAsset(symbol)
	}

	return asset
}

func defaultAsset(symbol string) *Asset {
	return &Asset{symbol, AssetClassOther, AssetClassOther}
}
