package e2e_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	app "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/usecases"
	handler "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/http/gin/catalog/get"
	repository "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/rest/catalog_repository"
)

func setupFullApp() *gin.Engine {
	gin.SetMode(gin.TestMode)
	dummyRepository := repository.NewDummyRepository()
	usecase := usecases.NewCatalogUsecase(domain.CatalogBackendDummy, dummyRepository, dummyRepository, dummyRepository)
	controller := app.NewController(usecase, app.NewRequestMapper(), app.NewResponseMapper())

	router := gin.New()
	router.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", handler.NewHandler(controller, handler.NewDecoder()))
	return router
}

func TestE2EFullCatalogFlow_DummyBackend(t *testing.T) {
	app := setupFullApp()
	request := httptest.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/Destino", nil)
	request.Header.Set("X-Channel", "APP")
	request.Header.Set("X-Commerce", "TOTTUS")
	request.Header.Set("X-Transaction-ID", "999888")
	response := httptest.NewRecorder()

	app.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	var items []map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &items))
	require.NotEmpty(t, items)
	assert.Equal(t, "1", items[0]["codigo_adm"])
	assert.Equal(t, "Vivienda Principal", items[0]["descripcion"])
	assert.Empty(t, response.Header().Get("X-Transaction-ID"))
}

func TestE2EFullCatalogFlow_RealBackend(t *testing.T) {
	// 1. Iniciar mock del servidor backend
	backendMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"codigo_adm": "REAL1", "descripcion": "Item desde backend real"}]`))
	}))
	defer backendMock.Close()

	// 2. Setup de repositorios apuntando al mock "real"
	gin.SetMode(gin.TestMode)
	realRepo := repository.NewRestRepository(backendMock.URL, "/", true, "user", "pass")

	usecase := usecases.NewCatalogUsecase(domain.CatalogBackendFinnflow, realRepo, realRepo, realRepo)
	controller := app.NewController(usecase, app.NewRequestMapper(), app.NewResponseMapper())

	router := gin.New()
	router.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", handler.NewHandler(controller, handler.NewDecoder()))

	// 3. Ejecutar request contra la app
	request := httptest.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/DestinoReal", nil)
	request.Header.Set("X-Channel", "APP")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	// 4. Validar respuestas
	require.Equal(t, http.StatusOK, response.Code)
	var items []map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &items))
	require.NotEmpty(t, items)
	assert.Equal(t, "REAL1", items[0]["codigo_adm"])
	assert.Equal(t, "Item desde backend real", items[0]["descripcion"])
}
