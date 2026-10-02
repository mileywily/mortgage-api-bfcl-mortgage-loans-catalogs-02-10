package repository_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	repository "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/rest/catalog_repository"
)

func TestRestRepository_GetCatalog_Finnflow_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		user, pass, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "user", user)
		assert.Equal(t, "pass", pass)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"IdentificadorSeguro": "123", "Descripcion": "Test Seguros"}]`))
	}))
	defer ts.Close()

	repo := repository.NewRestRepository(ts.URL, "/", true, "user", "pass")
	req := domain.GetCatalogRequest{CatalogName: "Seguros"}

	res, err := repo.GetCatalog(context.Background(), req)
	require.NoError(t, err)

	items, ok := res.([]domain.InsuranceCatalogItem)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.NotNil(t, items[0].InsuranceIdentifier)
	assert.Equal(t, "123", *items[0].InsuranceIdentifier)
	assert.NotNil(t, items[0].Description)
	assert.Equal(t, "Test Seguros", *items[0].Description)
}

func TestRestRepository_GetCatalog_Legacy_TiposDocumentos(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"codigo_adm": 1, "descripcion": "Doc 1", "grupo_id": 2}]`))
	}))
	defer ts.Close()

	repo := repository.NewRestRepository(ts.URL, "/", false, "", "")
	req := domain.GetCatalogRequest{
		CatalogName:   "TiposDocumentos",
		Channel:       "APP",
		Commerce:      "TOTTUS",
		TransactionID: "123",
	}

	res, err := repo.GetCatalog(context.Background(), req)
	require.NoError(t, err)

	items, ok := res.([]domain.TipoDocumentoCatalogItem)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.NotNil(t, items[0].Code)
	assert.Equal(t, 1, *items[0].Code)
	assert.Equal(t, "Doc 1", *items[0].Description)
}

func TestRestRepository_GetCatalog_Legacy_Comunas(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"codigo_adm": 10, "descripcion": "Santiago", "region_id": 13}]`))
	}))
	defer ts.Close()

	repo := repository.NewRestRepository(ts.URL, "/", false, "", "")
	req := domain.GetCatalogRequest{CatalogName: "Comunas"}

	res, err := repo.GetCatalog(context.Background(), req)
	require.NoError(t, err)

	items, ok := res.([]domain.ComunaCatalogItem)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.NotNil(t, items[0].Code)
	assert.Equal(t, "10", *items[0].Code)
	assert.Equal(t, "Santiago", *items[0].Description)
}

func TestRestRepository_GetCatalog_Legacy_Comunas_MissingCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"descripcion": "Santiago", "region_id": 13}]`))
	}))
	defer ts.Close()

	repo := repository.NewRestRepository(ts.URL, "/", false, "", "")
	req := domain.GetCatalogRequest{CatalogName: "Comunas"}

	res, err := repo.GetCatalog(context.Background(), req)
	require.NoError(t, err)

	items, ok := res.([]domain.ComunaCatalogItem)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.Nil(t, items[0].Code)
	assert.NotNil(t, items[0].Description)
}

func TestRestRepository_GetCatalog_Generic_NoResults(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"codRespuesta": 3}]`))
	}))
	defer ts.Close()

	repo := repository.NewRestRepository(ts.URL, "/", false, "", "")
	req := domain.GetCatalogRequest{CatalogName: "Otros"}

	res, err := repo.GetCatalog(context.Background(), req)
	require.NoError(t, err)

	items, ok := res.([]domain.NoResultsCatalogItem)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.Equal(t, 3, items[0].CodRespuesta)
}

func TestRestRepository_GetCatalog_EmptyResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer ts.Close()

	repo := repository.NewRestRepository(ts.URL, "/", false, "", "")
	req := domain.GetCatalogRequest{CatalogName: "Otros"}

	res, err := repo.GetCatalog(context.Background(), req)
	require.NoError(t, err)

	items, ok := res.([]domain.NoResultsCatalogItem)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.Equal(t, 3, items[0].CodRespuesta)
}

func TestRestRepository_GetCatalog_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	repo := repository.NewRestRepositoryWithTimeout(ts.URL, "/", false, "", "", 1*time.Millisecond)
	req := domain.GetCatalogRequest{CatalogName: "Otros"}

	res, err := repo.GetCatalog(context.Background(), req)
	require.NoError(t, err)

	items, ok := res.([]domain.NoResultsCatalogItem)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.Equal(t, 3, items[0].CodRespuesta)
}

func TestRestRepository_GetCatalog_Errors(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		expectedErr error
	}{
		{"Unauthorized", http.StatusUnauthorized, domain.ErrAuthentication},
		{"NotFound", http.StatusNotFound, domain.ErrCatalogNotFound},
		{"PaymentRequired", http.StatusPaymentRequired, domain.ErrCatalogNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer ts.Close()

			repo := repository.NewRestRepository(ts.URL, "/", false, "", "")
			req := domain.GetCatalogRequest{CatalogName: "Otros"}

			_, err := repo.GetCatalog(context.Background(), req)
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
