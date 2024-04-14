package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Cr4z1k/Avito-test-task/internal/core"
	"github.com/Cr4z1k/Avito-test-task/internal/repository"
	"github.com/Cr4z1k/Avito-test-task/internal/service"
	"github.com/Cr4z1k/Avito-test-task/internal/transport/rest/handlers"
	"github.com/Cr4z1k/Avito-test-task/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	host     = "localhost:8000" // 192.168.56.1
	basePath = "http://" + host
)

type GetTestSuite struct {
	suite.Suite
	Migrate     *migrate.Migrate
	TestDB      *sqlx.DB
	TestRouter  *gin.Engine
	TestManager auth.Manager
	UserToken   string
	AdmToken    string
	PgConfig    *postgres.Config
}

func getConnectionString() string {
	type Config struct {
		Postgres struct {
			Host   string `yaml:"host"`
			Port   string `yaml:"port"`
			Dbname string `yaml:"dbname"`
			User   string `yaml:"user"`
		} `yaml:"postgres"`
	}

	configData, err := os.ReadFile("../internal/config/conf.yaml")
	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Postgres.User, os.Getenv("DB_PASS"), "test", config.Postgres.Host, config.Postgres.Port)

	return connectionString
}

func getConnectionToDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", getConnectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *GetTestSuite) migrateDB() {
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

func (s *GetTestSuite) setupDB() {
	pg, err := getConnectionToDB()
	if err != nil {
		panic(errors.WithStack(err))
	}

	s.TestDB = pg

	s.migrateDB()
}

func (s *GetTestSuite) SetupRoutes(h *handlers.Handler) {
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

func (s *GetTestSuite) SetupSuite() {
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

func (s *GetTestSuite) TearDownSuite() {
	if err := s.Migrate.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(errors.WithStack(err))
	}

	if err := s.TestDB.Close(); err != nil {
		panic(errors.WithStack(err))
	}
}

func (s *GetTestSuite) TestGetBanner() {
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

func (s *GetTestSuite) TestGetBannerWithFilters() {
	s.T().Run("GetBannerWithFilters by tag adm success", func(t *testing.T) {

		req, err := http.NewRequest("GET", basePath+"/banner?tag_id=4", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var response []core.BannerWithFilters

		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}

		if len(response) != 3 {
			t.Fatalf("Expected 3 JSON objects in result, got %d instead", len(response))
		}
	})

	s.T().Run("GetBannerWithFilters by feature adm success", func(t *testing.T) {

		req, err := http.NewRequest("GET", basePath+"/banner?feature_id=1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var response []core.BannerWithFilters

		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}

		if len(response) != 2 {
			t.Fatalf("Expected 2 JSON objects in result, got %d instead", len(response))
		}
	})

	s.T().Run("GetBannerWithFilters by feature adm success", func(t *testing.T) {

		req, err := http.NewRequest("GET", basePath+"/banner?feature_id=1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var response []core.BannerWithFilters

		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}

		if len(response) != 2 {
			t.Fatalf("Expected 2 JSON objects in result, got %d instead", len(response))
		}
	})

	s.T().Run("GetBannerWithFilters with limit adm success", func(t *testing.T) {

		req, err := http.NewRequest("GET", basePath+"/banner?limit=1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var response []core.BannerWithFilters

		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}

		if len(response) != 1 {
			t.Fatalf("Expected 1 JSON object in result, got %d instead", len(response))
		}
	})

	s.T().Run("GetBannerWithFilters with offset adm success", func(t *testing.T) {

		req, err := http.NewRequest("GET", basePath+"/banner?offset=2", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.AdmToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d but got %d", http.StatusOK, w.Code)
		}

		var response []core.BannerWithFilters

		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err)
		}

		if len(response) != 2 {
			t.Fatalf("Expected 2 JSON objects in result, got %d instead", len(response))
		}
	})

	s.T().Run("GetBannerWithFilters adm fail, bad param", func(t *testing.T) {

		req, err := http.NewRequest("GET", basePath+"/banner?tag_id=gdgf", nil)
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

		if err := json.Unmarshal(w.Body.Bytes(), &expectedJSON); err != nil {
			t.Fatalf("Failed to unmarshal response body: %s", err.Error())
		}
	})

	s.T().Run("GetBannerWithFilters fail, unauthorized", func(t *testing.T) {

		req, err := http.NewRequest("GET", basePath+"/banner?offset=2", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, w.Code)
		}
	})

	s.T().Run("GetBannerWithFilters user fail, forbidden", func(t *testing.T) {

		req, err := http.NewRequest("GET", basePath+"/banner?offset=2", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %s", err)
		}

		req.Header.Set("Authorization", "Bearer "+s.UserToken)

		w := httptest.NewRecorder()
		s.TestRouter.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status code %d but got %d", http.StatusForbidden, w.Code)
		}
	})
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestHttpSuite(t *testing.T) {
	suite.Run(t, new(GetTestSuite))
	suite.Run(t, new(AddTestSuite))
	suite.Run(t, new(UpdTestSuite))
	suite.Run(t, new(DelTestSuite))
}
