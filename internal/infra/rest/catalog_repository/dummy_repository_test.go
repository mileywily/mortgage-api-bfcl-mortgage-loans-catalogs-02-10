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
	
	// Test Destino
	res, err := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Destino"})
	require.NoError(t, err)
	assert.NotNil(t, res)

	// Test Seguros
	resSeg, errSeg := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "SegurosIncendio"})
	require.NoError(t, errSeg)
	assert.NotNil(t, resSeg)

	// Test Not Found (Inexistente)
	resErr, errErr := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Inexistente"})
	require.ErrorIs(t, errErr, domain.ErrCatalogNotFound)
	assert.Nil(t, resErr)
}
