package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cr4z1k/Avito-test-task/internal/repository"
	"github.com/Cr4z1k/Avito-test-task/internal/service"
	"github.com/Cr4z1k/Avito-test-task/internal/transport/rest/handlers"
	"github.com/Cr4z1k/Avito-test-task/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type DelTestSuite struct {
	suite.Suite
	Migrate     *migrate.Migrate
	TestDB      *sqlx.DB
	TestRouter  *gin.Engine
	TestManager auth.Manager
	UserToken   string
	AdmToken    string
	PgConfig    *postgres.Config
}

func (s *DelTestSuite) migrateDB() {
	driver, err := postgres.WithInstance(s.TestDB.DB, &postgres.Config{})
	if err != nil {
		panic(errors.WithStack(err))
	}

	m, err := migrate.NewWithDatabaseInstance("file://../migration/", "test", driver)
	if err != nil {
		panic(errors.WithStack(err))
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.T().Fatalf("Failed to apply migrations: %s", err.Error())
	}

	s.Migrate = m
}

func (s *DelTestSuite) setupDB() {
	pg, err := getConnectionToDB()
	if err != nil {
		panic(errors.WithStack(err))
	}

	s.TestDB = pg

	s.migrateDB()
}

func (s *DelTestSuite) SetupRoutes(h *handlers.Handler) {
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

func (s *DelTestSuite) SetupSuite() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(errors.WithStack(err))
	}

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

func (s *DelTestSuite) TearDownSuite() {
	if err := s.Migrate.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(errors.WithStack(err))
	}

	if err := s.TestDB.Close(); err != nil {
		panic(errors.WithStack(err))
	}
}

func (s *DelTestSuite) TestDeleteBanner() {
	s.T().Run("DeleteBanner admin success", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", basePath+"/banner/1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err.Error())
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d but got %d", http.StatusNoContent, w.Code)
		}
	})

	s.T().Run("DeleteBanner admin fail, bad param", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", basePath+"/banner/bad", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err.Error())
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

	s.T().Run("DeleteBanner fail, no token", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", basePath+"/banner/1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err.Error())
		}

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, w.Code)
		}
	})

	s.T().Run("DeleteBanner user fail, forbidden", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", basePath+"/banner/1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err.Error())
		}

		req.Header.Set("Authorization", "Bearer "+s.UserToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status code %d but got %d", http.StatusForbidden, w.Code)
		}
	})

	s.T().Run("DeleteBanner admin fail, no such banner", func(t *testing.T) {

		req, err := http.NewRequest("DELETE", basePath+"/banner/100", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err.Error())
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d but got %d", http.StatusNotFound, w.Code)
		}
	})
}
