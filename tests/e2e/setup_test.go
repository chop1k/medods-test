package e2e

import (
	"database/sql"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	stdhttp "net/http"

	"github.com/chop1k/medods-test/internal/app"
	"github.com/chop1k/medods-test/internal/app/config"
	"github.com/chop1k/medods-test/internal/transport/http"
	"github.com/stretchr/testify/require"
)

var (
	testDB     *sql.DB
	testURL    string
	testClient *stdhttp.Client
	testSeed   int64
)

func TruncateDB(t *testing.T) {
	_, err := testDB.Exec("truncate table \"app\".\"templates\", \"app\".\"tasks\", \"app\".\"tags\", \"app\".\"templates_tags\" cascade")

	require.Nil(t, err, "cannot truncate test database: ", err)
}

func TestMain(m *testing.M) {
	cfg := config.Config{
		DB: &config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "app_db",
			SSLMode:  "disable",
		},
	}

	db, err := app.OpenDatabase(cfg.DB)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	srv := httptest.NewServer(
		http.NewRouter(db),
	)
	defer srv.Close()

	testDB = db
	testURL = srv.URL
	testClient = &stdhttp.Client{
		Timeout: time.Second * 10,
	}
	testSeed = time.Now().UnixNano()

	os.Exit(
		m.Run(),
	)
}
