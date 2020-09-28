package careercenter

import "time"

// CareerCenter ...
type CareerCenter struct {
	IsOpened            bool      `bson:"is_opened"`
	Title               string    `bson:"title"`
	Description         string    `bson:"description"`
	CVButtonEnabled     bool      `bson:"cb_button_enabled"`
	CustomButtonEnabled bool      `bson:"custom_button_enabled"`
	CustomButtonTitle   string    `bson:"custom_button_title"`
	CustomButtonURL     string    `bson:"custom_button_url"`
	CreatedAt           time.Time `bson:"created_at"`
}
