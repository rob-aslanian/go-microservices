package post

import (
	"testing"
)

func TestFindHashtags(t *testing.T) {
	p := Post{
		Text: `$qwe 123
	  asfd #test sd
	  igoehr 09gh32u4bu
	  b0g8734th jn#
	  jkn# # kfd
	  #test2 #я_very_Красивая`,
	}

	p.FindHashtags()

	if len(p.Hashtags) != 3 {
		t.Errorf("wrong FindHashtags(). Got: %q", p.Hashtags)
	}
}

func TestTrim(t *testing.T) {
	p := Post{
		Text: " a  b\t\tc\n\nd",
	}

	p.Trim()

	if p.Text != "a b c\nd" {
		t.Errorf("wrong TestTrim(). Got: %q", p.Text)
	}
}
