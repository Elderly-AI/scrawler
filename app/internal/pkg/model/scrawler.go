package model

type TagDB struct {
	ID         uint64 `db:"tag_id"`
	ExternalID uint64 `db:"external_id"`
	Title      string `db:"tag_title"`
	Total      uint64 `db:"total"`
}

type LessonLogDB struct {
	ID           uint64  `db:"lesson_log_id"`
	ExternalID   uint64  `db:"external_id"`
	LessonsCount float64 `db:"lessons_count"`
}

type LessonLog struct {
	ID           uint64
	ExternalID   uint64
	LessonsCount float64
	Tags         []Tag
}

type ExternalTags struct {
	Tags []ExternalTag `json:"tags"`
}

type ExternalTag struct {
	Title string `json:"nameRu"`
	ID    uint64 `json:"id"`
}

type Tag struct {
	ID         uint64
	ExternalID uint64
	Title      string
}

type ExternalLessonResponse struct {
	Result ExternalLessonResult `json:"result"`
}

type ExternalLessonResult struct {
	Entities []ExternalEntity `json:"entities"`
	Count    uint64           `json:"Count"`
}

type ExternalEntity struct {
	ID            uint64                      `json:"id"`
	MentorProfile ExternalEntityMentorProfile `json:"mentorProfile"`
}

type ExternalEntityMentorProfile struct {
	Statistics ExternalMentorProfileStatistics `json:"statistics"`
	Tags       []ExternalTag                   `json:"tags"`
}

type ExternalMentorProfileStatistics struct {
	Sessions ExternalMentorProfileStatisticsSessions `json:"sessions"`
}

type ExternalMentorProfileStatisticsSessions struct {
	Count float64 `json:"count"`
}
