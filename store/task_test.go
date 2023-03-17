package store

import (
	"context"
	"testing"

	"github.com/Luftalian/shokai_golang_webapp/clock"
	"github.com/Luftalian/shokai_golang_webapp/entity"
	"github.com/Luftalian/shokai_golang_webapp/testutil"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()
	// Test will fail if it mixes up other tests that make entity.Task
	// So, putting on transaction, make it table status only this test case.
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// If this test accomplished, rollback
	t.Cleanup(func() {
		_ = tx.Rollback()
	})
	if err != nil {
		t.Fatal(err)
	}
	wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err:= sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("got differs: (-got +want)\n%s", d)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()

	// Make it clean
	if _, err := con.ExecContext(ctx, "DELETE FROM tasks"); err != nil {
		t.Logf("failed to initialize tasks: %v", err)
	}
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{Title: "want task 1", Status: entity.TaskStatusTodo, Created: c.Now(), Modified: c.Now()},
		{Title: "want task 2", Status: entity.TaskStatusTodo, Created: c.Now(), Modified: c.Now()},
		{Title: "want task 3", Status: entity.TaskStatusDone, Created: c.Now(), Modified: c.Now()},
	}
	result, err := con.ExecContext(ctx, `INSERT INTO tasks (title, status, created, modified) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?);`, wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified, wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified, wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title: "ok task",
		Status: entity.TaskStatusTodo,
		Created: c.Now(),
		Modified: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func () {
		_ = db.Close()
	})
	mock.ExpectExec(`INSERT INTO tasks \(title, status, created, modified\) VALUES \(\?, \?, \?, \?\)`).WithArgs(okTask.Title, okTask.Status, okTask.Created, okTask.Modified).WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}