package office

// ContentType to know which type of file user is uploading in v-office
type ContentType string

const (
	// ContentTypeImage ...
	ContentTypeImage ContentType = "image"
	// ContentTypeArticle ...
	ContentTypeArticle ContentType = "article"
	// ContentTypeCode ...
	ContentTypeCode ContentType = "code"
	// ContentTypeVideo ...
	ContentTypeVideo ContentType = "video"
	// ContentTypeAudio ...
	ContentTypeAudio ContentType = "audio"
	//ContentTypeOther ...
	ContentTypeOther ContentType = "other"
)
