package data

import (
	"testing"
)

func TestTagParse(t *testing.T) {
	var src string = `
# Not

## A

### Tag

#tag1
#tag2/subtag
	`
	cache, _, err := MakeNoteCache([]byte(src))
	if err != nil {
		t.Error("Failed to parse source", err)
		return
	}

	expected := []Tag{
		"tag1",
		"tag2/subtag",
	}

	assertTagsMatch(t, expected, cache)

}

func TestTagParseFrontmatter(t *testing.T) {
	src := `---
tags:
  - tag1
  - tag2
---`

	cache, _, err := MakeNoteCache([]byte(src))
	if err != nil {
		t.Error("Failed to parse source", err)
		return
	}

	expected := []Tag{
		"tag1",
		"tag2",
	}

	assertTagsMatch(t, expected, cache)

}

func TestTagParseWeird(t *testing.T) {
	src := `---
tags:
  - tag1
  - -tag2
---

#tag-3
#-tag4

`

	cache, _, err := MakeNoteCache([]byte(src))
	if err != nil {
		t.Error("Failed to parse source", err)
		return
	}

	expected := []Tag{
		"tag1",
		"-tag2",
		"tag-3",
		"-tag4",
	}

	assertTagsMatch(t, expected, cache)

}

func TestTagParseBrokenFrontmatter(t *testing.T) {
	src := `---
tags:
  - tag1
  - tag2-tag3
	- notatag
---

# not a tag
#yes-a-tag
---

`
	expected := []Tag{
		"yes-a-tag",
	}

	cache, _, err := MakeNoteCache([]byte(src))
	if err != nil {
		t.Error("Failed to parse source", err)
		return
	}

	assertTagsMatch(t, expected, cache)

}

func assertTagsMatch(t *testing.T, expected []Tag, cache NoteCache) {
	if len(expected) != cache.Tags.Len() {
		t.Errorf("Tag Mismatch: Expected %v, got %v", expected, cache.Tags)
		return
	}

	failed := false

	for _, tag := range expected {
		if !cache.Tags.Contains(tag) {
			t.Errorf("Expected but could not find tag '%v'", tag)
			failed = true
		}
	}
	if failed {
		t.Logf("Got %v", cache.Tags)
	}
}
