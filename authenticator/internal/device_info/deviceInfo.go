package deviceinfo

// DeviceInfo ...
type DeviceInfo struct {
	Browser Browser `bson:"browser"`
	OS      OS      `bson:"os"`
	Type    string  `bson:"type"`
}

// Browser ...
type Browser struct {
	Name    string  `bson:"name"`
	Version Version `bson:"version"`
}

// OS ...
type OS struct {
	Name     string  `bson:"name"`
	Platform string  `bson:"platform"`
	Version  Version `bson:"version"`
}

// Version ...
type Version struct {
	Major int `bson:"major"`
	Minor int `bson:"minor"`
	Patch int `bson:"patch"`
}
