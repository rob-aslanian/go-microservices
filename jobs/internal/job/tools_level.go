package job

// Level presents the level of poficiency of each tool 7 technologie
type Level string

const (
	// LevelBeginner to choose the tools knowledge level
	LevelBeginner Level = "beginner"

	// LevelIntermediate to choose the tools knowledge level
	LevelIntermediate Level = "intermediate"

	// LevelAdvanced to choose the tools knowledge level
	LevelAdvanced Level = "advanced"

	// LevelMaster to choose the tools knowledge level
	LevelMaster Level = "master"
)
