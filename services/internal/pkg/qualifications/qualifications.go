package qualifications

// Qualifications contains skills, languages and tool&technology. it can be used
// in services, request and also in v-office.
type Qualifications struct {
	Skills               []*Skill          `bson:"skills"`
	ToolsAndTechnologies []*ToolTechnology `bson:"tool_technology"`
	Languages            []*Language       `bson:"languages"`
}

