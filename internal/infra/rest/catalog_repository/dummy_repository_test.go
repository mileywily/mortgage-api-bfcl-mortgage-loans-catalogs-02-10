package catalog_repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	repository "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/rest/catalog_repository"
)

func TestDummyRepository_GetCatalog(t *testing.T) {
	repo := repository.NewDummyRepository()
	
	// Test Destino (Genérico)
	res, err := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Destino"})
	require.NoError(t, err)
	assert.NotNil(t, res)

	// Test Seguros
	resSeg, errSeg := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "SegurosIncendio"})
	require.NoError(t, errSeg)
	assert.NotNil(t, resSeg)
	
	// Test TiposDocumentos
	resDoc, errDoc := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "TiposDocumentos"})
	require.NoError(t, errDoc)
	assert.NotNil(t, resDoc)
	
	// Test Comunas
	resCom, errCom := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Comunas"})
	require.NoError(t, errCom)
	assert.NotNil(t, resCom)

	// Test Vacio / Error
	resVac, errVac := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Vacio"})
	require.NoError(t, errVac)
	assert.NotNil(t, resVac)

	// Test Fallback (Cualquier otro nombre devuelve el genérico en Dummy)
	resErr, errErr := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Inexistente"})
	require.NoError(t, errErr)
	assert.NotNil(t, resErr)
}
