// Package curriculum provides access to the Tunisian educational curriculum
// catalogue (Sections, Subjects, and Chapters) seeded in migration 0003.
package curriculum

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("curriculum entity not found")
)

// Section represents a Tunisian Baccalauréat branch/stream (e.g. Math, Sciences, Technique).
type Section struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	TitleFR   string    `json:"titleFr"`
	TitleAR   string    `json:"titleAr"`
	CreatedAt time.Time `json:"createdAt"`
}

// Subject represents a subject within a Baccalauréat section (e.g. Mathematics).
type Subject struct {
	ID          int64     `json:"id"`
	SectionID   *int64    `json:"sectionId,omitempty"`
	SectionCode string    `json:"sectionCode,omitempty"`
	Code        string    `json:"code"`
	TitleFR     string    `json:"titleFr,omitempty"`
	TitleAR     string    `json:"titleAr,omitempty"`
	Coefficient *int      `json:"coefficient"`
	SortOrder   int       `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Chapter represents a catalogue entry; educational validation is a separate process.
type Chapter struct {
	ID            int64     `json:"id"`
	SubjectID     int64     `json:"subjectId"`
	SubjectCode   string    `json:"subjectCode"`
	Code          string    `json:"code"`
	Domain        string    `json:"domain"`
	TitleFR       string    `json:"titleFr"`
	TitleAR       string    `json:"titleAr"`
	DescriptionTN string    `json:"descriptionTn"`
	SortOrder     int       `json:"sortOrder"`
	CreatedAt     time.Time `json:"createdAt"`
}

// Store reads the curriculum catalogue from PostgreSQL.
type Store struct {
	db *pgxpool.Pool
}

// NewStore creates a new Store bound to the database pool.
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// ListSections returns the registered pilot sections.
func (s *Store) ListSections(ctx context.Context) ([]Section, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, code, title_fr, title_ar, created_at
		  FROM sections
		 ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list sections: %w", err)
	}
	defer rows.Close()

	var sections []Section
	for rows.Next() {
		var sec Section
		if err := rows.Scan(&sec.ID, &sec.Code, &sec.TitleFR, &sec.TitleAR, &sec.CreatedAt); err != nil {
			return nil, err
		}
		sections = append(sections, sec)
	}
	return sections, rows.Err()
}

// ListSubjects returns all subjects, optionally filtered by section code.
func (s *Store) ListSubjects(ctx context.Context, sectionCode string) ([]Subject, error) {
	query := `
		SELECT s.id, s.section_id, COALESCE(sec.code, ''), s.code,
		       COALESCE(s.title_fr, ''), COALESCE(s.title_ar, ''),
		       s.coefficient, COALESCE(s.sort_order, 0), s.created_at
		  FROM subjects s
		  LEFT JOIN sections sec ON sec.id = s.section_id
	`
	var rows pgx.Rows
	var err error
	if sectionCode != "" {
		query += ` WHERE sec.code = $1 ORDER BY s.sort_order ASC, s.id ASC`
		rows, err = s.db.Query(ctx, query, sectionCode)
	} else {
		query += ` ORDER BY s.sort_order ASC, s.id ASC`
		rows, err = s.db.Query(ctx, query)
	}
	if err != nil {
		return nil, fmt.Errorf("list subjects: %w", err)
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var sub Subject
		if err := rows.Scan(
			&sub.ID, &sub.SectionID, &sub.SectionCode, &sub.Code,
			&sub.TitleFR, &sub.TitleAR, &sub.Coefficient, &sub.SortOrder, &sub.CreatedAt,
		); err != nil {
			return nil, err
		}
		subjects = append(subjects, sub)
	}
	return subjects, rows.Err()
}

// ListChapters returns all chapters for a given subject code, ordered by sort_order.
func (s *Store) ListChapters(ctx context.Context, subjectCode string) ([]Chapter, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.id, c.subject_id, s.code, c.code, c.domain,
		       c.title_fr, c.title_ar, c.description_tn, c.sort_order, c.created_at
		  FROM chapters c
		  JOIN subjects s ON s.id = c.subject_id
		 WHERE s.code = $1
		 ORDER BY c.sort_order ASC, c.id ASC`, subjectCode)
	if err != nil {
		return nil, fmt.Errorf("list chapters: %w", err)
	}
	defer rows.Close()

	var chapters []Chapter
	for rows.Next() {
		var ch Chapter
		if err := rows.Scan(
			&ch.ID, &ch.SubjectID, &ch.SubjectCode, &ch.Code, &ch.Domain,
			&ch.TitleFR, &ch.TitleAR, &ch.DescriptionTN, &ch.SortOrder, &ch.CreatedAt,
		); err != nil {
			return nil, err
		}
		chapters = append(chapters, ch)
	}
	return chapters, rows.Err()
}

// GetChapter returns a single chapter by subject code and chapter code.
func (s *Store) GetChapter(ctx context.Context, subjectCode, chapterCode string) (Chapter, error) {
	var ch Chapter
	err := s.db.QueryRow(ctx, `
		SELECT c.id, c.subject_id, s.code, c.code, c.domain,
		       c.title_fr, c.title_ar, c.description_tn, c.sort_order, c.created_at
		  FROM chapters c
		  JOIN subjects s ON s.id = c.subject_id
		 WHERE s.code = $1 AND c.code = $2`, subjectCode, chapterCode).Scan(
		&ch.ID, &ch.SubjectID, &ch.SubjectCode, &ch.Code, &ch.Domain,
		&ch.TitleFR, &ch.TitleAR, &ch.DescriptionTN, &ch.SortOrder, &ch.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Chapter{}, ErrNotFound
	}
	if err != nil {
		return Chapter{}, fmt.Errorf("get chapter: %w", err)
	}
	return ch, nil
}
