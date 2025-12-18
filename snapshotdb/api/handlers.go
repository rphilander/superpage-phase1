package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"snapshotdb/scheduler"
	"snapshotdb/store"
)

type Handler struct {
	store     *store.Store
	scheduler *scheduler.Scheduler
}

func NewHandler(st *store.Store, sched *scheduler.Scheduler) *Handler {
	return &Handler{
		store:     st,
		scheduler: sched,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/status", h.handleStatus)
	mux.HandleFunc("/snapshots", h.handleSnapshots)
	mux.HandleFunc("/stories", h.handleStories)
	mux.HandleFunc("/story/", h.handleStory)
	mux.HandleFunc("/doc", h.handleDoc)
}

func (h *Handler) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stats := h.scheduler.Stats().GetSnapshot()
	now := time.Now()

	resp := StatusResponse{
		UptimeSeconds:   int64(now.Sub(stats.StartedAt).Seconds()),
		StartedAt:       stats.StartedAt.Unix(),
		SnapshotsTotal:  stats.SnapshotsTotal,
		SnapshotsErrors: stats.SnapshotsErrors,
	}

	if stats.LastSnapshotAt != nil {
		ts := stats.LastSnapshotAt.Unix()
		resp.LastSnapshotAt = &ts
	}
	if stats.NextSnapshotAt != nil {
		ts := stats.NextSnapshotAt.Unix()
		resp.NextSnapshotAt = &ts
	}

	h.writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleSnapshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	from, to, err := h.parseTimeRange(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	snapshots, err := h.store.GetSnapshotsInRange(from, to)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "database error: "+err.Error())
		return
	}

	dtos := make([]SnapshotDTO, len(snapshots))
	for i, snap := range snapshots {
		stories := make([]StoryDTO, len(snap.Stories))
		for j, s := range snap.Stories {
			stories[j] = StoryDTO{
				StoryID:       s.StoryID,
				Rank:          s.Rank,
				Headline:      s.Headline,
				URL:           s.URL,
				Username:      s.Username,
				Points:        s.Points,
				Comments:      s.Comments,
				DiscussionURL: s.DiscussionURL,
				AgeValue:      s.AgeValue,
				AgeUnit:       s.AgeUnit,
				Page:          s.Page,
			}
		}
		dtos[i] = SnapshotDTO{
			ID:           snap.ID,
			FetchedAt:    snap.FetchedAt.Unix(),
			NumPages:     snap.NumPages,
			TotalStories: snap.TotalStories,
			Stories:      stories,
		}
	}

	h.writeJSON(w, http.StatusOK, SnapshotsResponse{Snapshots: dtos})
}

func (h *Handler) handleStories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	from, to, err := h.parseTimeRange(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	storyIDs, err := h.store.GetStoryIDsInRange(from, to)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "database error: "+err.Error())
		return
	}

	if storyIDs == nil {
		storyIDs = []string{}
	}

	h.writeJSON(w, http.StatusOK, StoriesResponse{
		StoryIDs: storyIDs,
		Count:    len(storyIDs),
	})
}

func (h *Handler) handleStory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	storyID := strings.TrimPrefix(r.URL.Path, "/story/")
	if storyID == "" {
		h.writeError(w, http.StatusBadRequest, "story ID required")
		return
	}

	from, to, err := h.parseTimeRange(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	occurrences, err := h.store.GetStoryInRange(storyID, from, to)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "database error: "+err.Error())
		return
	}

	if len(occurrences) == 0 {
		h.writeError(w, http.StatusNotFound, "story not found in the specified time range")
		return
	}

	dtos := make([]StoryOccurrenceDTO, len(occurrences))
	for i, occ := range occurrences {
		dtos[i] = StoryOccurrenceDTO{
			SnapshotID:    occ.SnapshotID,
			FetchedAt:     occ.FetchedAt.Unix(),
			Rank:          occ.Rank,
			Headline:      occ.Headline,
			URL:           occ.URL,
			Username:      occ.Username,
			Points:        occ.Points,
			Comments:      occ.Comments,
			DiscussionURL: occ.DiscussionURL,
			AgeValue:      occ.AgeValue,
			AgeUnit:       occ.AgeUnit,
			Page:          occ.Page,
		}
	}

	h.writeJSON(w, http.StatusOK, StoryResponse{
		StoryID:     storyID,
		Occurrences: dtos,
	})
}

func (h *Handler) handleDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	doc := DocResponse{
		Name:        "SnapshotDB API",
		Version:     "1.0.0",
		Description: "REST API for querying Hacker News snapshots stored by SnapshotDB",
		Endpoints: []EndpointDoc{
			{
				Method:      "GET",
				Path:        "/status",
				Description: "Returns operational statistics about the SnapshotDB service",
				Response: ResponseDoc{
					ContentType: "application/json",
					Description: "Service status including uptime, snapshot counts, and timing information",
				},
				Example: &EndpointExample{
					Request: "GET /status",
					Response: StatusResponse{
						UptimeSeconds:   3600,
						StartedAt:       1702382400,
						SnapshotsTotal:  60,
						SnapshotsErrors: 2,
						LastSnapshotAt:  ptrInt64(1702386000),
						NextSnapshotAt:  ptrInt64(1702386060),
					},
				},
			},
			{
				Method:      "GET",
				Path:        "/snapshots",
				Description: "Returns all snapshots within a time window, including full story data",
				Parameters: []ParameterDoc{
					{Name: "from", Type: "integer", Required: true, Description: "Start of time window (Unix timestamp)"},
					{Name: "to", Type: "integer", Required: true, Description: "End of time window (Unix timestamp)"},
				},
				Response: ResponseDoc{
					ContentType: "application/json",
					Description: "Array of snapshots with all story data",
				},
				Example: &EndpointExample{
					Request: "GET /snapshots?from=1702382400&to=1702386000",
					Response: SnapshotsResponse{
						Snapshots: []SnapshotDTO{
							{
								ID:           1,
								FetchedAt:    1702382400,
								NumPages:     4,
								TotalStories: 120,
								Stories: []StoryDTO{
									{
										StoryID:       "46243904",
										Rank:          1,
										Headline:      "SQLite JSON at Full Index Speed",
										URL:           "https://example.com/article",
										Username:      "author",
										Points:        206,
										Comments:      74,
										DiscussionURL: "https://news.ycombinator.com/item?id=46243904",
										AgeValue:      5,
										AgeUnit:       "hours",
										Page:          1,
									},
								},
							},
						},
					},
				},
			},
			{
				Method:      "GET",
				Path:        "/stories",
				Description: "Returns deduplicated story IDs within a time window",
				Parameters: []ParameterDoc{
					{Name: "from", Type: "integer", Required: true, Description: "Start of time window (Unix timestamp)"},
					{Name: "to", Type: "integer", Required: true, Description: "End of time window (Unix timestamp)"},
				},
				Response: ResponseDoc{
					ContentType: "application/json",
					Description: "Array of unique story IDs and count",
				},
				Example: &EndpointExample{
					Request: "GET /stories?from=1702382400&to=1702386000",
					Response: StoriesResponse{
						StoryIDs: []string{"46243904", "46174114", "46245923"},
						Count:    3,
					},
				},
			},
			{
				Method:      "GET",
				Path:        "/story/{id}",
				Description: "Returns all data for a specific story across snapshots in a time window",
				Parameters: []ParameterDoc{
					{Name: "id", Type: "string", Required: true, Description: "Hacker News story ID (path parameter)"},
					{Name: "from", Type: "integer", Required: true, Description: "Start of time window (Unix timestamp)"},
					{Name: "to", Type: "integer", Required: true, Description: "End of time window (Unix timestamp)"},
				},
				Response: ResponseDoc{
					ContentType: "application/json",
					Description: "Story data across all snapshots in the time range",
				},
				Example: &EndpointExample{
					Request: "GET /story/46243904?from=1702382400&to=1702386000",
					Response: StoryResponse{
						StoryID: "46243904",
						Occurrences: []StoryOccurrenceDTO{
							{
								SnapshotID:    1,
								FetchedAt:     1702382400,
								Rank:          1,
								Headline:      "SQLite JSON at Full Index Speed",
								URL:           "https://example.com/article",
								Username:      "author",
								Points:        206,
								Comments:      74,
								DiscussionURL: "https://news.ycombinator.com/item?id=46243904",
								AgeValue:      5,
								AgeUnit:       "hours",
								Page:          1,
							},
							{
								SnapshotID:    2,
								FetchedAt:     1702382460,
								Rank:          1,
								Headline:      "SQLite JSON at Full Index Speed",
								URL:           "https://example.com/article",
								Username:      "author",
								Points:        210,
								Comments:      78,
								DiscussionURL: "https://news.ycombinator.com/item?id=46243904",
								AgeValue:      5,
								AgeUnit:       "hours",
								Page:          1,
							},
						},
					},
				},
			},
			{
				Method:      "GET",
				Path:        "/doc",
				Description: "Returns this API documentation",
				Response: ResponseDoc{
					ContentType: "application/json",
					Description: "Full API documentation with examples",
				},
			},
		},
	}

	h.writeJSON(w, http.StatusOK, doc)
}

func (h *Handler) parseTimeRange(r *http.Request) (time.Time, time.Time, error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" {
		return time.Time{}, time.Time{}, &paramError{param: "from", message: "required"}
	}
	if toStr == "" {
		return time.Time{}, time.Time{}, &paramError{param: "to", message: "required"}
	}

	fromUnix, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		return time.Time{}, time.Time{}, &paramError{param: "from", message: "must be a valid Unix timestamp"}
	}

	toUnix, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		return time.Time{}, time.Time{}, &paramError{param: "to", message: "must be a valid Unix timestamp"}
	}

	if fromUnix > toUnix {
		return time.Time{}, time.Time{}, &paramError{param: "from/to", message: "from must be less than or equal to to"}
	}

	return time.Unix(fromUnix, 0), time.Unix(toUnix, 0), nil
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, ErrorResponse{Error: message})
}

type paramError struct {
	param   string
	message string
}

func (e *paramError) Error() string {
	return e.param + ": " + e.message
}

func ptrInt64(v int64) *int64 {
	return &v
}
