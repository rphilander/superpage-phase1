package query

import (
	"strings"

	"github.com/sahilm/fuzzy"

	"querier/models"
)

const defaultFuzzyThreshold = 50

// Filter applies filters to stories and returns matching stories
func Filter(stories []models.Story, filters *models.Filters) []models.Story {
	if filters == nil {
		return stories
	}

	result := make([]models.Story, 0, len(stories))
	for _, story := range stories {
		if matchesAllFilters(story, filters) {
			result = append(result, story)
		}
	}
	return result
}

func matchesAllFilters(story models.Story, filters *models.Filters) bool {
	if filters.Headline != nil && !matchesFuzzy(story.Headline, filters.Headline) {
		return false
	}
	if filters.Username != nil && !matchesFuzzy(story.Username, filters.Username) {
		return false
	}
	if filters.URL != nil && !matchesFuzzy(story.URL, filters.URL) {
		return false
	}
	if filters.Points != nil && !matchesRange(story.Points, filters.Points) {
		return false
	}
	if filters.Comments != nil && !matchesRange(story.Comments, filters.Comments) {
		return false
	}
	if filters.Rank != nil && !matchesRange(story.Rank, filters.Rank) {
		return false
	}
	if filters.Page != nil && !matchesRange(story.Page, filters.Page) {
		return false
	}
	if filters.Age != nil && !matchesAge(story, filters.Age) {
		return false
	}
	return true
}

func matchesFuzzy(text string, filter *models.FuzzyFilter) bool {
	if filter.Match == "" {
		return true
	}

	threshold := defaultFuzzyThreshold
	if filter.Threshold != nil {
		threshold = *filter.Threshold
	}

	// Use fuzzy matching
	matches := fuzzy.Find(strings.ToLower(filter.Match), []string{strings.ToLower(text)})
	if len(matches) == 0 {
		return false
	}

	// Calculate score as percentage (0-100)
	// The fuzzy library returns a score where higher is better
	// We normalize based on the length of the matched text
	match := matches[0]
	maxScore := len(filter.Match) * 2 // Approximate max score
	score := (match.Score * 100) / maxScore
	if score > 100 {
		score = 100
	}

	return score >= threshold
}

func matchesRange(value int, r *models.RangeInt) bool {
	if r.Min != nil && value < *r.Min {
		return false
	}
	if r.Max != nil && value > *r.Max {
		return false
	}
	return true
}

func matchesAge(story models.Story, ageRange *models.AgeRange) bool {
	// Convert story age to minutes for comparison
	storyAgeMinutes := toMinutes(story.AgeValue, story.AgeUnit)

	// Convert filter range to minutes
	filterMin := 0
	filterMax := int(^uint(0) >> 1) // Max int

	if ageRange.Min != nil {
		filterMin = toMinutes(*ageRange.Min, ageRange.Unit)
	}
	if ageRange.Max != nil {
		filterMax = toMinutes(*ageRange.Max, ageRange.Unit)
	}

	return storyAgeMinutes >= filterMin && storyAgeMinutes <= filterMax
}

func toMinutes(value int, unit string) int {
	switch unit {
	case "minutes":
		return value
	case "hours":
		return value * 60
	case "days":
		return value * 60 * 24
	default:
		return value
	}
}
