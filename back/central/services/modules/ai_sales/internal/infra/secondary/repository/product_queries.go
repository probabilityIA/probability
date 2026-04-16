package repository

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/secamc93/probability/back/migration/shared/models"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

// accentTransformer quita diacriticos (acentos, tildes, cedillas) descomponiendo
// los caracteres Unicode (NFD) y removiendo las marcas combinantes.
// Resultado: "proteinas" == "proteinas", "camion" == "camion", etc.
var accentTransformer = transform.Chain(
	norm.NFD,
	transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}),
	norm.NFC,
)

// removeAccents quita acentos de un string usando descomposicion Unicode
func removeAccents(s string) string {
	result, _, _ := transform.String(accentTransformer, s)
	return result
}

// sqlAccentedChars y sqlPlainChars se usan en la funcion SQL translate()
// para normalizar acentos en columnas de la base de datos.
// Se generan programaticamente para evitar literals UTF-8 densos en el codigo fuente.
var sqlAccentedChars, sqlPlainChars = buildSQLTranslateChars()

func buildSQLTranslateChars() (string, string) {
	// Pares: caracter acentuado -> caracter sin acento
	// Se definen como runas individuales para evitar strings densos de multibyte
	pairs := [][2]rune{
		{'\u00e1', 'a'}, {'\u00e9', 'e'}, {'\u00ed', 'i'}, {'\u00f3', 'o'}, {'\u00fa', 'u'}, // a-grave, etc
		{'\u00e0', 'a'}, {'\u00e8', 'e'}, {'\u00ec', 'i'}, {'\u00f2', 'o'}, {'\u00f9', 'u'},
		{'\u00e2', 'a'}, {'\u00ea', 'e'}, {'\u00ee', 'i'}, {'\u00f4', 'o'}, {'\u00fb', 'u'},
		{'\u00e3', 'a'}, {'\u00f5', 'o'}, {'\u00f1', 'n'}, {'\u00e4', 'a'}, {'\u00eb', 'e'},
		{'\u00ef', 'i'}, {'\u00f6', 'o'}, {'\u00fc', 'u'},
		{'\u00c1', 'A'}, {'\u00c9', 'E'}, {'\u00cd', 'I'}, {'\u00d3', 'O'}, {'\u00da', 'U'},
		{'\u00c0', 'A'}, {'\u00c8', 'E'}, {'\u00cc', 'I'}, {'\u00d2', 'O'}, {'\u00d9', 'U'},
		{'\u00c2', 'A'}, {'\u00ca', 'E'}, {'\u00ce', 'I'}, {'\u00d4', 'O'}, {'\u00db', 'U'},
		{'\u00c3', 'A'}, {'\u00d5', 'O'}, {'\u00d1', 'N'}, {'\u00c7', 'C'},
		{'\u00c4', 'A'}, {'\u00cb', 'E'}, {'\u00cf', 'I'}, {'\u00d6', 'O'}, {'\u00dc', 'U'},
		{'\u00e7', 'c'},
	}

	var accented, plain strings.Builder
	for _, p := range pairs {
		accented.WriteRune(p[0])
		plain.WriteRune(p[1])
	}
	return accented.String(), plain.String()
}

// normalizeCol genera la expresion SQL para normalizar acentos en una columna:
// translate(lower(column), '<accented>', '<plain>')
func normalizeCol(col string) string {
	return fmt.Sprintf("translate(lower(%s), '%s', '%s')", col, sqlAccentedChars, sqlPlainChars)
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

	// Normalizar el patron de busqueda: quitar acentos y lowercase
	normalized := removeAccents(strings.ToLower(query))
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
		TrackInventory:   p.TrackInventory,
		Category:         p.Category,
		Brand:            p.Brand,
		ImageURL:         p.ImageURL,
		IsActive:         p.IsActive,
	}
}
