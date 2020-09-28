package shop

// SearchFilter ...
type SearchFilter struct {
	Keyword     string
	Category    []string
	Subcategory []string
	PriceMin    uint32
	PriceMax    uint32
	// Size        []string
	InStock *bool
	IsUsed  *bool
}
