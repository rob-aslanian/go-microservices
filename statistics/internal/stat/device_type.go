package stat

// Device ...
type Device struct {
	OS        string `bson:"os"`
	OSVersion string `bson:"os_version"`
	Browser   string `bson:"browser"`
	Type      string `bson:"type"`
}
