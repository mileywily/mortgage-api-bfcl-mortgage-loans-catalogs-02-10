package repository

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/out"
)

func ptrString(s string) *string  { return &s }
func ptrInt(i int) *int           { return &i }
func ptrFloat(f float64) *float64 { return &f }

type dummyRepository struct{}

// NewDummyRepository returns a mock repository matching the legacy Java data.
func NewDummyRepository() out.CatalogRepository {
	return &dummyRepository{}
}

func (r *dummyRepository) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	// Dummy response matching EXACTLY the Java legacy documentation
	if req.CatalogName == "SegurosIncendio" || req.CatalogName == "SegurosDesgravamen" || req.CatalogName == "SegurosCesantia" {
		return []domain.InsuranceCatalogItem{
			{
				InsuranceIdentifier:  ptrString("1100438-01-2023-000"),
				Description:          ptrString("INCENDIO - Everest compañía de seguros generales Chile(0.2255300)"),
				Policy:               ptrString("100438-01-2023-000"),
				CompanyName:          ptrString("Everest compañía de seguros generales Chile"),
				CompanyCode:          ptrInt(39),
				InsuranceTypeCode:    ptrInt(0),
				PolicyCorrelative:    ptrInt(28),
				Rate:                 ptrFloat(0.22553),
				Factor:               ptrFloat(0.0002255),
				IndividualPolicyFlag: ptrInt(0),
				PerQuotaValueFlag:    ptrInt(0),
				ExternalPolicyFlag:   ptrInt(0),
			},
		}, nil
	}

	if req.CatalogName == "TiposDocumentos" {
		return []domain.TipoDocumentoCatalogItem{
			{
				Code:        ptrInt(1),
				Description: ptrString("Escritura de Propiedad"),
				GroupID:     ptrInt(2),
			},
		}, nil
	}

	if req.CatalogName == "Vacio" || req.CatalogName == "ErrorRed" {
		// Simulates a network error or empty results from backend
		return []domain.NoResultsCatalogItem{
			{
				CodRespuesta: 3,
				Mensaje:      "Sin resultados.",
				Excepcion:    "Ninguna",
			},
		}, nil
	}

	if req.CatalogName == "Comunas" {
		return []domain.ComunaCatalogItem{
			{
				Code:        ptrString("01101"),
				Description: ptrString("Iquique"),
				RegionID:    ptrInt(1),
			},
		}, nil
	}

	// Standard Catalog exact match
	return []domain.CatalogItem{
		{
			Code:        ptrString("1"),
			Description: ptrString("Vivienda Principal"),
		},
	}, nil
}
