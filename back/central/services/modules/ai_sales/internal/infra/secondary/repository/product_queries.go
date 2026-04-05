package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/migration/shared/models"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

const (
	// Caracteres acentuados y sus equivalentes sin acento para búsqueda accent-insensitive
	accentedChars = "áéíóúàèìòùâêîôûãõñäëïöüÁÉÍÓÚÀÈÌÒÙÂÊÎÔÛÃÕÑÄËÏÖÜçÇ"
	plainChars    = "aeiouaeiouaeiouaonaeiouAEIOUAEIOUAEIOUAONAEIOUcC"
)

// accentReplacer normaliza acentos en Go para que el patrón de búsqueda coincida
// con la normalización SQL (translate) aplicada a las columnas
var accentReplacer = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
	"à", "a", "è", "e", "ì", "i", "ò", "o", "ù", "u",
	"â", "a", "ê", "e", "î", "i", "ô", "o", "û", "u",
	"ã", "a", "õ", "o", "ñ", "n", "ç", "c",
	"ä", "a", "ë", "e", "ï", "i", "ö", "o", "ü", "u",
	"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U",
	"À", "A", "È", "E", "Ì", "I", "Ò", "O", "Ù", "U",
	"Â", "A", "Ê", "E", "Î", "I", "Ô", "O", "Û", "U",
	"Ã", "A", "Õ", "O", "Ñ", "N", "Ç", "C",
	"ä", "a", "ë", "e", "ï", "i", "ö", "o", "ü", "u",
)

// normalizeCol genera la expresión SQL para normalizar acentos en una columna:
// translate(lower(column), 'áéí...', 'aei...')
func normalizeCol(col string) string {
	return fmt.Sprintf("translate(lower(%s), '%s', '%s')", col, accentedChars, plainChars)
}

// spanishStem aplica stemming básico para español: quita sufijos de plural
// para que "proteinas" también encuentre "proteina", "camisetas" encuentre "camiseta", etc.
func spanishStem(word string) string {
	// Separar en palabras y hacer stem de cada una
	words := strings.Fields(word)
	for i, w := range words {
		switch {
		case strings.HasSuffix(w, "es") && len(w) > 4:
			words[i] = strings.TrimSuffix(w, "es")
		case strings.HasSuffix(w, "s") && len(w) > 3:
			words[i] = strings.TrimSuffix(w, "s")
		}
	}
	return strings.Join(words, " ")
}

// SearchProducts busca productos por query accent-insensitive en name, title, description,
// short_description, category y brand.
// Tabla consultada: products (gestionada por modulo products)
// Replicado localmente para evitar compartir repositorios entre modulos
func (r *repository) SearchProducts(ctx context.Context, businessID uint, query string, limit int) ([]domain.ProductSearchResult, error) {
	if limit <= 0 || limit > 20 {
		limit = 5
	}

	// Normalizar el patrón de búsqueda: quitar acentos y lowercase
	normalized := accentReplacer.Replace(strings.ToLower(query))
	// Stemming básico español: quitar plural para ampliar la búsqueda
	stemmed := spanishStem(normalized)

	searchPattern := fmt.Sprintf("%%%s%%", normalized)
	stemmedPattern := fmt.Sprintf("%%%s%%", stemmed)

	// Buscar en 6 campos con normalización de acentos en ambos lados
	// Usa OR entre el patrón original y el stemmed para cubrir singular/plural
	col := func(name string) string { return normalizeCol(name) }
	searchCondition := fmt.Sprintf(
		"%s LIKE ? OR %s LIKE ? OR %s LIKE ? OR %s LIKE ? OR %s LIKE ? OR %s LIKE ?"+
			" OR %s LIKE ? OR %s LIKE ? OR %s LIKE ? OR %s LIKE ? OR %s LIKE ? OR %s LIKE ?",
		col("name"), col("title"), col("description"),
		col("short_description"), col("category"), col("brand"),
		col("name"), col("title"), col("description"),
		col("short_description"), col("category"), col("brand"),
	)

	var products []models.Product
	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Where("is_active = ?", true).
		Where("deleted_at IS NULL").
		Where(searchCondition,
			searchPattern, searchPattern, searchPattern,
			searchPattern, searchPattern, searchPattern,
			stemmedPattern, stemmedPattern, stemmedPattern,
			stemmedPattern, stemmedPattern, stemmedPattern,
		).
		Limit(limit).
		Find(&products).Error

	if err != nil {
		return nil, fmt.Errorf("error searching products: %w", err)
	}

	results := make([]domain.ProductSearchResult, 0, len(products))
	for _, p := range products {
		results = append(results, mapProductToDomain(&p))
	}

	return results, nil
}

// GetProductBySKU obtiene un producto por su SKU dentro de un negocio
func (r *repository) GetProductBySKU(ctx context.Context, businessID uint, sku string) (*domain.ProductSearchResult, error) {
	var product models.Product

	err := r.db.Conn(ctx).
		Where("business_id = ? AND sku = ? AND is_active = ? AND deleted_at IS NULL", businessID, sku, true).
		First(&product).Error

	if err != nil {
		return nil, &domain.ErrProductNotFound{SKU: sku}
	}

	result := mapProductToDomain(&product)
	return &result, nil
}

func mapProductToDomain(p *models.Product) domain.ProductSearchResult {
	return domain.ProductSearchResult{
		ID:               p.ID,
		SKU:              p.SKU,
		Name:             p.Name,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		Price:            p.Price,
		Currency:         p.Currency,
		StockQuantity:    p.StockQuantity,
		Category:         p.Category,
		Brand:            p.Brand,
		ImageURL:         p.ImageURL,
		IsActive:         p.IsActive,
	}
}
