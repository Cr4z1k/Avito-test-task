package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Cr4z1k/Avito-test-task/internal/repository"
	"github.com/Cr4z1k/Avito-test-task/internal/service"
	"github.com/Cr4z1k/Avito-test-task/internal/transport/rest/handlers"
	"github.com/Cr4z1k/Avito-test-task/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	host     = "192.168.56.1:8000"
	basePath = "http://" + host
)

type HttpTestSuite struct {
	suite.Suite
	Migrate     *migrate.Migrate
	TestDB      *sqlx.DB
	TestRouter  *gin.Engine
	TestManager auth.Manager
	UserToken   string
	AdmToken    string
	PgConfig    *postgres.Config
}

func getConnectionToDB() (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		"postgres", "postgres", "Avito-DB", "pg_database", "5432")

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *HttpTestSuite) migrateDB() {
	driver, err := postgres.WithInstance(s.TestDB.DB, &postgres.Config{})
	if err != nil {
		panic(errors.WithStack(err))
	}

	m, err := migrate.NewWithDatabaseInstance("file://../migration/", "Avito-DB", driver)
	if err != nil {
		panic(errors.WithStack(err))
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(errors.WithStack(err))
	}

	s.Migrate = m
}

func (s *HttpTestSuite) setupDB() {
	pg, err := getConnectionToDB()
	if err != nil {
		panic(errors.WithStack(err))
	}

	s.TestDB = pg

	s.migrateDB()
}

func (s *HttpTestSuite) SetupRoutes(h *handlers.Handler) {
	r := s.TestRouter

	mwToken := r.Group("", h.CheckTokenIsAdmin)
	{
		mwAdm := mwToken.Group("", h.IdentifyAdmin)
		{
			banner := mwAdm.Group("/banner")
			{
				banner.GET("", h.GetBannerWithFilter)
				banner.POST("", h.CreateBanner)
				banner.PATCH("/:id", h.UpdateBanner)
				banner.DELETE("/:id", h.DeleteBanner)
			}
		}

		mwToken.GET("/user_banner", h.GetBanner)
	}
}

func (s *HttpTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	s.TestRouter = gin.New()

	s.setupDB()

	tknMngr, err := auth.NewTokenManager("123")
	if err != nil {
		panic(errors.WithStack(err))
	}

	s.TestManager = *tknMngr

	s.AdmToken, err = s.TestManager.NewToken(true)
	if err != nil {
		panic(errors.WithStack(err))
	}

	s.UserToken, err = s.TestManager.NewToken(false)
	if err != nil {
		panic(errors.WithStack(err))
	}

	repo := repository.NewRepository(s.TestDB)
	service := service.NewService(repo, tknMngr)
	handler := handlers.NewHandler(service)

	s.SetupRoutes(handler)
}

func (s *HttpTestSuite) TearDownSuite() {
	if err := s.Migrate.Down(); err != nil {
		panic(errors.WithStack(err))
	}

	if err := s.TestDB.Close(); err != nil {
		panic(errors.WithStack(err))
	}
}

func (s *HttpTestSuite) TestGetBanner() {
	s.T().Run("GetBanner user success", func(t *testing.T) {
		req, err := http.NewRequest("GET", basePath+"/user_banner?tag_id=1&feature_id=1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.UserToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var expectedJSON struct {
			Title string `json:"title"`
			Text  string `json:"text"`
			Url   string `json:"url"`
		}

		err = json.Unmarshal(w.Body.Bytes(), &expectedJSON)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}
	})

	s.T().Run("GetBanner admin success", func(t *testing.T) {
		req, err := http.NewRequest("GET", basePath+"/user_banner?tag_id=1&feature_id=1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var expectedJSON struct {
			Title string `json:"title"`
			Text  string `json:"text"`
			Url   string `json:"url"`
		}

		err = json.Unmarshal(w.Body.Bytes(), &expectedJSON)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}
	})

	s.T().Run("GetBanner admin success, disabled banner", func(t *testing.T) {
		req, err := http.NewRequest("GET", basePath+"/user_banner?tag_id=3&feature_id=1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var expectedJSON struct {
			Title string `json:"title"`
			Text  string `json:"text"`
			Url   string `json:"url"`
		}

		err = json.Unmarshal(w.Body.Bytes(), &expectedJSON)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}
	})

	s.T().Run("GetBanner user fail, disabled banner", func(t *testing.T) {
		req, err := http.NewRequest("GET", basePath+"/user_banner?tag_id=3&feature_id=1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.UserToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, w.Code)
		}
	})

	s.T().Run("GetBanner fail, no JWT token provided", func(t *testing.T) {
		req, err := http.NewRequest("GET", basePath+"/user_banner?tag_id=3&feature_id=1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, w.Code)
		}
	})

	s.T().Run("GetBanner fail, bad data", func(t *testing.T) {
		req, err := http.NewRequest("GET", basePath+"/user_banner?tag_id=gdfg&feature_id=gdsg", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

		var expectedJSON struct {
			Error string `json:"error"`
		}

		err = json.Unmarshal(w.Body.Bytes(), &expectedJSON)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}
	})
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestHttpSuite(t *testing.T) {
	suite.Run(t, new(HttpTestSuite))
}
