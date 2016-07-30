package ironclad

type SortOrder string

const (
	PriceHighToLow SortOrder = "price"
	PriceLowToHigh SortOrder = "-price"
	AgeOldToNew    SortOrder = "age"
	AgeNewToOld    SortOrder = "-age"
)

func ParseSortOrder(f string) SortOrder {
	switch SortOrder(f) {
	case PriceHighToLow, PriceLowToHigh, AgeOldToNew, AgeNewToOld:
		return SortOrder(f)
	default:
		return AgeNewToOld
	}
}

type byOrder struct {
	L []Listing
	F SortOrder
}

func (o byOrder) Less(i, j int) bool {
	a, b := o.L[i], o.L[j]

	switch o.F {
	case AgeNewToOld:
		return a.LastUpdated.After(b.LastUpdated)
	case AgeOldToNew:
		return a.LastUpdated.Before(b.LastUpdated)
	case PriceLowToHigh:
		return a.Price < b.Price
	case PriceHighToLow:
		return a.Price > b.Price
	}
	return false
}
func (o byOrder) Swap(i, j int) { o.L[i], o.L[j] = o.L[j], o.L[i] }
func (o byOrder) Len() int      { return len(o.L) }
