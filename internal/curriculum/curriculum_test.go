package curriculum_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hadden-me/learnoo/internal/curriculum"
	"github.com/hadden-me/learnoo/internal/database/dbtest"
)

func TestPilotSeedHasNoUnverifiedCurriculum(t *testing.T) {
	db := dbtest.Migrated(t)
	ctx := context.Background()
	store := curriculum.NewStore(db)
	sections, err := store.ListSections(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Code != "math" {
		t.Fatalf("unexpected sections: %+v", sections)
	}
	subjects, err := store.ListSubjects(ctx, "math")
	if err != nil {
		t.Fatal(err)
	}
	if len(subjects) != 1 || subjects[0].Code != "bac-math-math" || subjects[0].Coefficient != nil {
		t.Fatalf("unexpected subjects: %+v", subjects)
	}
	chapters, err := store.ListChapters(ctx, "bac-math-math")
	if err != nil || len(chapters) != 0 {
		t.Fatalf("unverified chapters seeded: %+v, %v", chapters, err)
	}
	other, err := store.ListSubjects(ctx, "sciences")
	if err != nil || len(other) != 0 {
		t.Fatalf("non-pilot section seeded: %+v, %v", other, err)
	}
}

func TestCatalogueQueriesAndConstraints(t *testing.T) {
	db := dbtest.Migrated(t)
	ctx := context.Background()
	store := curriculum.NewStore(db)
	// Synthetic query fixtures only: no educational curriculum claims.
	var subjectID int64
	if err := db.QueryRow(ctx, "SELECT id FROM subjects WHERE code='bac-math-math'").Scan(&subjectID); err != nil {
		t.Fatal(err)
	}
	insert := `INSERT INTO chapters (subject_id,code,domain,title_fr,title_ar,sort_order) VALUES ($1,$2,'test','Fixture non validée','اختبار غير معتمد',$3)`
	for _, tc := range []struct {
		code  string
		order int
	}{{"test-second", 2}, {"test-first", 1}} {
		if _, err := db.Exec(ctx, insert, subjectID, tc.code, tc.order); err != nil {
			t.Fatal(err)
		}
	}
	chapters, err := store.ListChapters(ctx, "bac-math-math")
	if err != nil {
		t.Fatal(err)
	}
	if len(chapters) != 2 || chapters[0].Code != "test-first" || chapters[1].Code != "test-second" {
		t.Fatalf("wrong order: %+v", chapters)
	}
	ch, err := store.GetChapter(ctx, "bac-math-math", "test-first")
	if err != nil || ch.ID != chapters[0].ID || ch.DescriptionTN != "" {
		t.Fatalf("unexpected chapter: %+v, %v", ch, err)
	}
	if _, err := store.GetChapter(ctx, "unknown", "test-first"); !errors.Is(err, curriculum.ErrNotFound) {
		t.Fatalf("subject filter: %v", err)
	}
	if _, err := store.GetChapter(ctx, "bac-math-math", "missing"); !errors.Is(err, curriculum.ErrNotFound) {
		t.Fatalf("missing chapter: %v", err)
	}
	if _, err := db.Exec(ctx, insert, subjectID, "test-first", 3); err == nil {
		t.Fatal("duplicate chapter accepted")
	}
	if _, err := db.Exec(ctx, insert, -1, "orphan", 3); err == nil {
		t.Fatal("orphan chapter accepted")
	}
	if _, err := db.Exec(ctx, "DELETE FROM subjects WHERE id=$1", subjectID); err == nil {
		t.Fatal("referenced subject deleted")
	}
	if _, err := db.Exec(ctx, "DELETE FROM sections WHERE code='math'"); err == nil {
		t.Fatal("referenced section deleted")
	}
	// Existing reviewer scopes need not yet belong to a catalogue section.
	if _, err := db.Exec(ctx, "INSERT INTO subjects (code) VALUES ('test-legacy')"); err != nil {
		t.Fatal(err)
	}
	subjects, err := store.ListSubjects(ctx, "")
	if err != nil || len(subjects) != 2 {
		t.Fatalf("legacy scope query: %+v, %v", subjects, err)
	}
}
