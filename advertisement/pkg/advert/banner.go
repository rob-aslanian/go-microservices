package advert

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Banner represent banner advert
type Banner struct {
	Advert `bson:",inline"`
	// Places         []Place   `bson:"places"`
	// DestinationURL string    `bson:"destination_url"`
	IsResponsive bool      `bson:"is_responsive"`
	Contents     []Content `bson:"contents"`
}

// Validate ...
func (b Banner) Validate() error {
	err := b.Advert.Validate()
	if err != nil {
		return err
	}

	// if len(b.Places) == 0 {
	// 	return errors.New("wrong_places")
	// }

	// if b.DestinationURL == "" {
	// 	return errors.New("empty_destination_url")
	// }

	if len(b.Contents) == 0 {
		return errors.New("empty_content")
	}

	for i, c := range b.Contents {
		if i == 0 {
			if c.DestinationURL == "" {
				return errors.New("empty_destination_url")
			}
		}

		// if b.IsResponsive {
		if c.Content == "" {
			return errors.New("empty_content")
		}
		if utf8.RuneCountInString(c.Content) > 100 {
			return errors.New("content_is_too_long")
		}
		if c.Title == "" {
			return errors.New("empty_title")
		}
		if utf8.RuneCountInString(c.Title) > 40 {
			return errors.New("title_is_too_long")
		}
		// }
	}

	return nil
}

// Trim ...
func (b *Banner) Trim() {
	b.Advert.Trim()
	// b.DestinationURL = strings.TrimSpace(b.DestinationURL)

	// if b.IsResponsive {
	for i := range b.Contents {
		b.Contents[i].Content = strings.TrimSpace(b.Contents[i].Content)
		b.Contents[i].Title = strings.TrimSpace(b.Contents[i].Title)
		b.Contents[i].DestinationURL = strings.TrimSpace(b.Contents[i].DestinationURL)
	}
	// }
}
