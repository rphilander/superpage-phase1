package query

import (
	"sort"
	"strings"

	"querier/models"
)

// Sort sorts stories according to the sort specifications
func Sort(stories []models.Story, specs []models.SortSpec) []models.Story {
	if len(specs) == 0 {
		return stories
	}

	// Use stable sort to preserve order for equal elements
	sort.SliceStable(stories, func(i, j int) bool {
		for _, spec := range specs {
			cmp := compareField(stories[i], stories[j], spec.Field)
			if cmp == 0 {
				continue
			}
			if spec.Direction == "desc" {
				return cmp > 0
			}
			return cmp < 0
		}
		return false
	})

	return stories
}

func compareField(a, b models.Story, field string) int {
	switch field {
	case "id":
		return strings.Compare(a.ID, b.ID)
	case "headline":
		return strings.Compare(a.Headline, b.Headline)
	case "url":
		return strings.Compare(a.URL, b.URL)
	case "discussion_url":
		return strings.Compare(a.DiscussionURL, b.DiscussionURL)
	case "username":
		return strings.Compare(a.Username, b.Username)
	case "points":
		return compareInt(a.Points, b.Points)
	case "comments":
		return compareInt(a.Comments, b.Comments)
	case "rank":
		return compareInt(a.Rank, b.Rank)
	case "page":
		return compareInt(a.Page, b.Page)
	case "age_value":
		return compareInt(a.AgeValue, b.AgeValue)
	case "age":
		// Compare normalized age in minutes
		aMinutes := toMinutes(a.AgeValue, a.AgeUnit)
		bMinutes := toMinutes(b.AgeValue, b.AgeUnit)
		return compareInt(aMinutes, bMinutes)
	default:
		return 0
	}
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
