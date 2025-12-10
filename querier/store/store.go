package store

import (
	"sync"

	"querier/models"
)

// Store is a thread-safe in-memory store for stories
type Store struct {
	mu        sync.RWMutex
	stories   []models.Story
	fetchedAt string
}

// New creates a new empty store
func New() *Store {
	return &Store{
		stories: nil,
	}
}

// Update replaces the stored data with new stories
func (s *Store) Update(stories []models.Story, fetchedAt string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stories = stories
	s.fetchedAt = fetchedAt
}

// Get returns a copy of all stories and the fetched_at timestamp
func (s *Store) Get() ([]models.Story, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.stories == nil {
		return nil, s.fetchedAt
	}

	// Return a copy to prevent modification
	copied := make([]models.Story, len(s.stories))
	copy(copied, s.stories)
	return copied, s.fetchedAt
}

// IsEmpty returns true if no data has been loaded
func (s *Store) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stories == nil
}

// Count returns the number of stories
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.stories)
}
