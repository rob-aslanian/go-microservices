package rental

// Status ...
type Status string

const (
	// StatusAny ...
	StatusAny Status = "any"
	// StatusOldBuild ...
	StatusOldBuild Status = "old_build"
	// StatusNewBuilding ...
	StatusNewBuilding Status = "new_building"
	// StatusUnderConstruction ...
	StatusUnderConstruction Status = "under_construction"
	// StatusDeveloped ...
	StatusDeveloped Status = "developed"
	// StatusBuildable ...
	StatusBuildable Status = "buildable"
	// StatusNonBuilding ...
	StatusNonBuilding Status = "non_building"
)
