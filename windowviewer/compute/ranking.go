package compute

import (
	"fmt"
	"sort"

	"windowviewer/models"
)

// Criteria constants
const (
	CriteriaBestRank            = "best_rank"
	CriteriaMaxPoints           = "max_points"
	CriteriaMaxComments         = "max_comments"
	CriteriaIncrementalComments = "incremental_comments"
	CriteriaIncrementalPoints   = "incremental_points"
)

// ValidCriteria returns all valid criteria values
func ValidCriteria() []string {
	return []string{
		CriteriaBestRank,
		CriteriaMaxPoints,
		CriteriaMaxComments,
		CriteriaIncrementalComments,
		CriteriaIncrementalPoints,
	}
}

// IsValidCriteria checks if the given criteria is valid
func IsValidCriteria(criteria string) bool {
	for _, c := range ValidCriteria() {
		if c == criteria {
			return true
		}
	}
	return false
}

// storyData holds aggregated data for a story across all snapshots
type storyData struct {
	storyID       string
	headline      string
	url           string
	username      string
	discussionURL string

	// For best_rank (lower is better)
	bestRank int

	// For max_points
	maxPoints int

	// For max_comments
	maxComments int

	// For incremental calculations
	firstPoints   int
	lastPoints    int
	firstComments int
	lastComments  int
	firstFetchedAt int64
	lastFetchedAt  int64
}

// ComputeTopStories computes the top stories based on the given criteria
func ComputeTopStories(snapshots []models.Snapshot, criteria string, limit int) ([]models.RankedStory, error) {
	if !IsValidCriteria(criteria) {
		return nil, fmt.Errorf("invalid criteria: %s", criteria)
	}

	// Aggregate story data across all snapshots
	storyMap := make(map[string]*storyData)

	// Sort snapshots by fetched_at to ensure correct ordering for incremental calculations
	sortedSnapshots := make([]models.Snapshot, len(snapshots))
	copy(sortedSnapshots, snapshots)
	sort.Slice(sortedSnapshots, func(i, j int) bool {
		return sortedSnapshots[i].FetchedAt < sortedSnapshots[j].FetchedAt
	})

	for _, snapshot := range sortedSnapshots {
		for _, story := range snapshot.Stories {
			sd, exists := storyMap[story.StoryID]
			if !exists {
				sd = &storyData{
					storyID:        story.StoryID,
					headline:       story.Headline,
					url:            story.URL,
					username:       story.Username,
					discussionURL:  story.DiscussionURL,
					bestRank:       story.Rank,
					maxPoints:      story.Points,
					maxComments:    story.Comments,
					firstPoints:    story.Points,
					lastPoints:     story.Points,
					firstComments:  story.Comments,
					lastComments:   story.Comments,
					firstFetchedAt: snapshot.FetchedAt,
					lastFetchedAt:  snapshot.FetchedAt,
				}
				storyMap[story.StoryID] = sd
			} else {
				// Update best rank (lower is better)
				if story.Rank < sd.bestRank {
					sd.bestRank = story.Rank
				}

				// Update max values
				if story.Points > sd.maxPoints {
					sd.maxPoints = story.Points
				}
				if story.Comments > sd.maxComments {
					sd.maxComments = story.Comments
				}

				// Update latest values (snapshots are sorted by time)
				if snapshot.FetchedAt > sd.lastFetchedAt {
					sd.lastPoints = story.Points
					sd.lastComments = story.Comments
					sd.lastFetchedAt = snapshot.FetchedAt
				}

				// Keep headline and other metadata from latest snapshot
				sd.headline = story.Headline
				sd.url = story.URL
				sd.username = story.Username
				sd.discussionURL = story.DiscussionURL
			}
		}
	}

	// Convert to slice and compute metric
	stories := make([]struct {
		data   *storyData
		metric int
	}, 0, len(storyMap))

	for _, sd := range storyMap {
		var metric int
		switch criteria {
		case CriteriaBestRank:
			// For best_rank, lower rank number is better, but we want to sort descending by "goodness"
			// So we negate the rank: rank 1 becomes -1, rank 10 becomes -10
			// When we sort descending by metric, -1 > -10, so rank 1 comes first
			metric = -sd.bestRank
		case CriteriaMaxPoints:
			metric = sd.maxPoints
		case CriteriaMaxComments:
			metric = sd.maxComments
		case CriteriaIncrementalComments:
			metric = sd.lastComments - sd.firstComments
		case CriteriaIncrementalPoints:
			metric = sd.lastPoints - sd.firstPoints
		}

		stories = append(stories, struct {
			data   *storyData
			metric int
		}{data: sd, metric: metric})
	}

	// Sort by metric descending
	sort.Slice(stories, func(i, j int) bool {
		return stories[i].metric > stories[j].metric
	})

	// Take top N
	if limit > len(stories) {
		limit = len(stories)
	}

	result := make([]models.RankedStory, limit)
	for i := 0; i < limit; i++ {
		sd := stories[i].data

		result[i] = models.RankedStory{
			StoryID:             sd.storyID,
			Headline:            sd.headline,
			URL:                 sd.url,
			Username:            sd.username,
			DiscussionURL:       sd.discussionURL,
			BestRank:            sd.bestRank,
			MaxPoints:           sd.maxPoints,
			MaxComments:         sd.maxComments,
			IncrementalPoints:   sd.lastPoints - sd.firstPoints,
			IncrementalComments: sd.lastComments - sd.firstComments,
		}
	}

	return result, nil
}
