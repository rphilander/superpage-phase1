package store

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

type Snapshot struct {
	ID           int64
	FetchedAt    time.Time
	NumPages     int
	TotalStories int
	Stories      []Story
}

type Story struct {
	ID            int64
	SnapshotID    int64
	StoryID       string
	Rank          int
	Headline      string
	URL           string
	Username      string
	Points        int
	Comments      int
	DiscussionURL string
	AgeValue      int
	AgeUnit       string
	Page          int
}

func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.createTables(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) createTables() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			fetched_at DATETIME NOT NULL,
			num_pages INTEGER NOT NULL,
			total_stories INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS stories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			snapshot_id INTEGER NOT NULL,
			story_id TEXT NOT NULL,
			rank INTEGER NOT NULL,
			headline TEXT NOT NULL,
			url TEXT,
			username TEXT,
			points INTEGER NOT NULL,
			comments INTEGER NOT NULL,
			discussion_url TEXT,
			age_value INTEGER NOT NULL,
			age_unit TEXT NOT NULL,
			page INTEGER NOT NULL,
			FOREIGN KEY (snapshot_id) REFERENCES snapshots(id)
		);

		CREATE INDEX IF NOT EXISTS idx_snapshots_fetched_at ON snapshots(fetched_at);
		CREATE INDEX IF NOT EXISTS idx_stories_snapshot_id ON stories(snapshot_id);
		CREATE INDEX IF NOT EXISTS idx_stories_story_id ON stories(story_id);
	`)
	return err
}

func (s *Store) SaveSnapshot(snapshot *Snapshot) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO snapshots (fetched_at, num_pages, total_stories) VALUES (?, ?, ?)",
		snapshot.FetchedAt, snapshot.NumPages, snapshot.TotalStories,
	)
	if err != nil {
		return err
	}

	snapshotID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO stories (snapshot_id, story_id, rank, headline, url, username, points, comments, discussion_url, age_value, age_unit, page)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, story := range snapshot.Stories {
		_, err := stmt.Exec(
			snapshotID, story.StoryID, story.Rank, story.Headline, story.URL,
			story.Username, story.Points, story.Comments, story.DiscussionURL,
			story.AgeValue, story.AgeUnit, story.Page,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) GetLastSnapshotTime() (*time.Time, error) {
	var fetchedAt time.Time
	err := s.db.QueryRow("SELECT fetched_at FROM snapshots ORDER BY fetched_at DESC LIMIT 1").Scan(&fetchedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fetchedAt, nil
}

func (s *Store) GetSnapshotCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM snapshots").Scan(&count)
	return count, err
}

func (s *Store) GetSnapshotsInRange(from, to time.Time) ([]Snapshot, error) {
	rows, err := s.db.Query(
		"SELECT id, fetched_at, num_pages, total_stories FROM snapshots WHERE fetched_at >= ? AND fetched_at <= ? ORDER BY fetched_at",
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []Snapshot
	for rows.Next() {
		var snap Snapshot
		if err := rows.Scan(&snap.ID, &snap.FetchedAt, &snap.NumPages, &snap.TotalStories); err != nil {
			return nil, err
		}

		stories, err := s.getStoriesForSnapshot(snap.ID)
		if err != nil {
			return nil, err
		}
		snap.Stories = stories
		snapshots = append(snapshots, snap)
	}

	return snapshots, rows.Err()
}

func (s *Store) getStoriesForSnapshot(snapshotID int64) ([]Story, error) {
	rows, err := s.db.Query(`
		SELECT id, snapshot_id, story_id, rank, headline, url, username, points, comments, discussion_url, age_value, age_unit, page
		FROM stories WHERE snapshot_id = ? ORDER BY rank
	`, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []Story
	for rows.Next() {
		var story Story
		if err := rows.Scan(
			&story.ID, &story.SnapshotID, &story.StoryID, &story.Rank, &story.Headline,
			&story.URL, &story.Username, &story.Points, &story.Comments, &story.DiscussionURL,
			&story.AgeValue, &story.AgeUnit, &story.Page,
		); err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}

	return stories, rows.Err()
}

func (s *Store) GetStoryIDsInRange(from, to time.Time) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT DISTINCT stories.story_id
		FROM stories
		JOIN snapshots ON stories.snapshot_id = snapshots.id
		WHERE snapshots.fetched_at >= ? AND snapshots.fetched_at <= ?
		ORDER BY stories.story_id
	`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storyIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		storyIDs = append(storyIDs, id)
	}

	return storyIDs, rows.Err()
}

func (s *Store) GetStoryInRange(storyID string, from, to time.Time) ([]StoryOccurrence, error) {
	rows, err := s.db.Query(`
		SELECT stories.id, stories.snapshot_id, snapshots.fetched_at, stories.rank, stories.headline,
			stories.url, stories.username, stories.points, stories.comments, stories.discussion_url,
			stories.age_value, stories.age_unit, stories.page
		FROM stories
		JOIN snapshots ON stories.snapshot_id = snapshots.id
		WHERE stories.story_id = ? AND snapshots.fetched_at >= ? AND snapshots.fetched_at <= ?
		ORDER BY snapshots.fetched_at
	`, storyID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occurrences []StoryOccurrence
	for rows.Next() {
		var occ StoryOccurrence
		if err := rows.Scan(
			&occ.ID, &occ.SnapshotID, &occ.FetchedAt, &occ.Rank, &occ.Headline,
			&occ.URL, &occ.Username, &occ.Points, &occ.Comments, &occ.DiscussionURL,
			&occ.AgeValue, &occ.AgeUnit, &occ.Page,
		); err != nil {
			return nil, err
		}
		occurrences = append(occurrences, occ)
	}

	return occurrences, rows.Err()
}

type StoryOccurrence struct {
	ID            int64
	SnapshotID    int64
	FetchedAt     time.Time
	Rank          int
	Headline      string
	URL           string
	Username      string
	Points        int
	Comments      int
	DiscussionURL string
	AgeValue      int
	AgeUnit       string
	Page          int
}

func (s *Store) Close() error {
	return s.db.Close()
}
