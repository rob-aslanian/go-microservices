package location

// LocationType shows us if the service is remote or on-site
type LocationType string

const (
	// LocationTypeRemoTeOnly LocationType = "RemoTeOnly"
	LocationTypeRemoTeOnly LocationType = "RemoteOnly"

	// LocationTypeOnSiteWork OnSiteWork = "OnSiteWork"
	LocationTypeOnSiteWork LocationType = "OnSiteWork"
)
