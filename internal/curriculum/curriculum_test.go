package curriculum_test

import (
	"context"
	"testing"

	"github.com/hadden-me/learnoo/internal/curriculum"
	"github.com/hadden-me/learnoo/internal/database/dbtest"
)

func TestCurriculumCatalogueSeed(t *testing.T) {
	db := dbtest.Migrated(t)
	ctx := context.Background()
	store := curriculum.NewStore(db)

	t.Run("Sections", func(t *testing.T) {
		sections, err := store.ListSections(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sections) < 6 {
			t.Fatalf("expected at least 6 sections, got %d", len(sections))
		}

		codes := map[string]bool{}
		for _, s := range sections {
			codes[s.Code] = true
			if s.TitleFR == "" || s.TitleAR == "" {
				t.Errorf("section %s missing bilingual titles", s.Code)
			}
		}
		for _, required := range []string{"math", "sciences", "technique", "informatique", "economie", "lettres"} {
			if !codes[required] {
				t.Errorf("missing official Tunisian section %q", required)
			}
		}
	})

	t.Run("Subjects", func(t *testing.T) {
		subjects, err := store.ListSubjects(ctx, "math")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(subjects) == 0 {
			t.Fatalf("expected subjects for math section, got 0")
		}
		found := false
		for _, s := range subjects {
			if s.Code == "bac-math-math" {
				found = true
				if s.Coefficient != 4 {
					t.Errorf("expected coefficient 4 for Bac Math, got %d", s.Coefficient)
				}
				if s.TitleFR != "Mathématiques" || s.TitleAR != "الرياضيات" {
					t.Errorf("unexpected subject titles: %s / %s", s.TitleFR, s.TitleAR)
				}
			}
		}
		if !found {
			t.Errorf("missing subject bac-math-math")
		}
	})

	t.Run("Chapters", func(t *testing.T) {
		chapters, err := store.ListChapters(ctx, "bac-math-math")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(chapters) != 12 {
			t.Fatalf("expected 12 official Tunisian Bac Math chapters, got %d", len(chapters))
		}

		for i, ch := range chapters {
			expectedOrder := i + 1
			if ch.SortOrder != expectedOrder {
				t.Errorf("chapter %s expected sort_order %d, got %d", ch.Code, expectedOrder, ch.SortOrder)
			}
			if ch.TitleFR == "" || ch.TitleAR == "" {
				t.Errorf("chapter %s missing title", ch.Code)
			}
			if len(ch.DescriptionTN) < 10 {
				t.Errorf("chapter %s missing Tunisian Derja explanation", ch.Code)
			}
			if ch.Domain != "analyse" && ch.Domain != "algebre-geometrie" && ch.Domain != "probabilites-statistiques" {
				t.Errorf("chapter %s has invalid domain %q", ch.Code, ch.Domain)
			}
		}

		// Verify specific key chapter
		complexCh, err := store.GetChapter(ctx, "bac-math-math", "nombres-complexes")
		if err != nil {
			t.Fatalf("get chapter: %v", err)
		}
		if complexCh.TitleFR != "Nombres complexes" || complexCh.TitleAR != "الأعداد المركبة" {
			t.Errorf("unexpected complex numbers chapter titles")
		}
	})

	t.Run("UnknownChapter", func(t *testing.T) {
		_, err := store.GetChapter(ctx, "bac-math-math", "non-existent")
		if err != curriculum.ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}
