package scheduler

import (
	"log"
	"sync"
	"time"

	"snapshotdb/logger"
	"snapshotdb/parser"
	"snapshotdb/store"
)

type Stats struct {
	StartedAt       time.Time
	SnapshotsTotal  int
	SnapshotsErrors int
	LastSnapshotAt  *time.Time
	NextSnapshotAt  *time.Time
	mu              sync.RWMutex
}

func (s *Stats) GetSnapshot() StatsSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return StatsSnapshot{
		StartedAt:       s.StartedAt,
		SnapshotsTotal:  s.SnapshotsTotal,
		SnapshotsErrors: s.SnapshotsErrors,
		LastSnapshotAt:  s.LastSnapshotAt,
		NextSnapshotAt:  s.NextSnapshotAt,
	}
}

type StatsSnapshot struct {
	StartedAt       time.Time
	SnapshotsTotal  int
	SnapshotsErrors int
	LastSnapshotAt  *time.Time
	NextSnapshotAt  *time.Time
}

type Scheduler struct {
	store    *store.Store
	parser   *parser.Client
	logger   *logger.Logger
	freq     time.Duration
	stats    *Stats
	stopCh   chan struct{}
	stoppedCh chan struct{}
}

func New(st *store.Store, p *parser.Client, l *logger.Logger, freqSecs int) *Scheduler {
	return &Scheduler{
		store:    st,
		parser:   p,
		logger:   l,
		freq:     time.Duration(freqSecs) * time.Second,
		stats: &Stats{
			StartedAt: time.Now(),
		},
		stopCh:   make(chan struct{}),
		stoppedCh: make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go s.run()
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
	<-s.stoppedCh
}

func (s *Scheduler) Stats() *Stats {
	return s.stats
}

func (s *Scheduler) run() {
	defer close(s.stoppedCh)

	lastSnapshot, err := s.store.GetLastSnapshotTime()
	if err != nil {
		log.Printf("Error getting last snapshot time: %v", err)
		s.logError(err.Error(), "get_last_snapshot_time")
	}

	var nextFetch time.Time
	if lastSnapshot == nil {
		log.Println("No previous snapshots, fetching immediately")
		nextFetch = time.Now()
	} else {
		nextFetch = lastSnapshot.Add(s.freq)
		if nextFetch.Before(time.Now()) {
			nextFetch = time.Now()
		}
		log.Printf("Last snapshot at %s, next fetch at %s", lastSnapshot.Format(time.RFC3339), nextFetch.Format(time.RFC3339))
	}

	s.stats.mu.Lock()
	s.stats.NextSnapshotAt = &nextFetch
	s.stats.mu.Unlock()

	initialDelay := time.Until(nextFetch)
	if initialDelay > 0 {
		select {
		case <-time.After(initialDelay):
		case <-s.stopCh:
			return
		}
	}

	s.fetchSnapshot()

	ticker := time.NewTicker(s.freq)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.fetchSnapshot()
		case <-s.stopCh:
			return
		}
	}
}

func (s *Scheduler) fetchSnapshot() {
	next := time.Now().Add(s.freq)
	s.stats.mu.Lock()
	s.stats.NextSnapshotAt = &next
	s.stats.mu.Unlock()

	log.Println("Fetching snapshot from parser...")
	snapshot, err := s.parser.Fetch()
	if err != nil {
		log.Printf("Error fetching snapshot: %v", err)
		s.logError(err.Error(), "fetch_snapshot")
		s.stats.mu.Lock()
		s.stats.SnapshotsErrors++
		s.stats.mu.Unlock()
		return
	}

	log.Printf("Received snapshot with %d stories, saving to database...", len(snapshot.Stories))
	if err := s.store.SaveSnapshot(snapshot); err != nil {
		log.Printf("Error saving snapshot: %v", err)
		s.logError(err.Error(), "save_snapshot")
		s.stats.mu.Lock()
		s.stats.SnapshotsErrors++
		s.stats.mu.Unlock()
		return
	}

	now := time.Now()
	s.stats.mu.Lock()
	s.stats.SnapshotsTotal++
	s.stats.LastSnapshotAt = &now
	s.stats.mu.Unlock()

	log.Printf("Snapshot saved successfully (total: %d)", s.stats.SnapshotsTotal)
}

func (s *Scheduler) logError(errMsg, context string) {
	if err := s.logger.LogError(errMsg, context); err != nil {
		log.Printf("Error writing to error log: %v", err)
	}
}
