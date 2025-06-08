package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// run migrations for test db
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("truncating hables err: %v", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)
	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "push day",
				Description:     "upper body day",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "bench press",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FltPtr(135.5),
						Notes:        "warm up",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:           "full body",
				Description:     "full body workout",
				DurationMinutes: 90,
				CaloriesBurned:  100,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPtr(60),
						Notes:        "form",
						OrderIndex:   1,
					},
					{
						ExerciseName: "squats",
						Sets:         4,
						Reps:         IntPtr(12),
						Weight:       FltPtr(185.0),
						Notes:        "whatever",
						OrderIndex:   2,
					},
				},
			},
		},
	}
}

func IntPtr(i int) *int {
	return &i
}

func FltPtr(i float64) *float64 {
	return &i
}
