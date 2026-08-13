package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type meliItemsSearchResponse struct {
	Results []string `json:"results"`
	Paging  struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	} `json:"paging"`
}

type meliItemAttribute struct {
	ID        string `json:"id"`
	ValueName string `json:"value_name"`
}

type meliItemVariation struct {
	ID                    int64               `json:"id"`
	AvailableQuantity     int                 `json:"available_quantity"`
	SellerCustomField     string              `json:"seller_custom_field"`
	Attributes            []meliItemAttribute `json:"attributes"`
	AttributeCombinations []struct {
		Name      string `json:"name"`
		ValueName string `json:"value_name"`
	} `json:"attribute_combinations"`
}

func variationLabel(v meliItemVariation) string {
	parts := make([]string, 0, len(v.AttributeCombinations))
	for _, c := range v.AttributeCombinations {
		if strings.TrimSpace(c.ValueName) != "" {
			parts = append(parts, c.ValueName)
		}
	}
	return strings.Join(parts, " / ")
}

func variationName(title string, v meliItemVariation) string {
	label := variationLabel(v)
	if label == "" {
		return title
	}
	return title + " - " + label
}

type meliItemDetail struct {
	ID                string              `json:"id"`
	Title             string              `json:"title"`
	CategoryID        string              `json:"category_id"`
	Price             float64             `json:"price"`
	AvailableQuantity int                 `json:"available_quantity"`
	SellerCustomField string              `json:"seller_custom_field"`
	Attributes        []meliItemAttribute `json:"attributes"`
	Variations        []meliItemVariation `json:"variations"`
	Shipping          struct {
		LogisticType string `json:"logistic_type"`
	} `json:"shipping"`
	Thumbnail string `json:"thumbnail"`
	Pictures  []struct {
		SecureURL string `json:"secure_url"`
		URL       string `json:"url"`
	} `json:"pictures"`
}

func itemImageURL(d meliItemDetail) string {
	for _, picture := range d.Pictures {
		if picture.SecureURL != "" {
			return picture.SecureURL
		}
		if picture.URL != "" {
			return picture.URL
		}
	}
	return d.Thumbnail
}

type meliMultigetItem struct {
	Code int            `json:"code"`
	Body meliItemDetail `json:"body"`
}

func extractSKUFrom(sellerCustomField string, attrs []meliItemAttribute) string {
	if strings.TrimSpace(sellerCustomField) != "" {
		return sellerCustomField
	}
	for _, attr := range attrs {
		if attr.ID == "SELLER_SKU" {
			return attr.ValueName
		}
	}
	return ""
}

func extractSKU(item meliItemDetail) string {
	return extractSKUFrom(item.SellerCustomField, item.Attributes)
}

func atributo(attrs []meliItemAttribute, ids ...string) string {
	for _, wanted := range ids {
		for _, attr := range attrs {
			if attr.ID == wanted && strings.TrimSpace(attr.ValueName) != "" {
				return attr.ValueName
			}
		}
	}
	return ""
}

func extractBarcodeFrom(attrs ...[]meliItemAttribute) string {
	for _, group := range attrs {
		if value := atributo(group, "GTIN", "EAN", "UPC"); value != "" {
			return value
		}
	}
	return ""
}

func extractBarcode(item meliItemDetail) string {
	return extractBarcodeFrom(item.Attributes)
}

func (c *MeliClient) GetProducts(ctx context.Context, accessToken string, sellerID int64) ([]domain.MeliProduct, error) {
	itemIDs := make([]string, 0)
	offset := 0
	limit := 50

	for {
		endpoint := fmt.Sprintf("%s/users/%d/items/search?limit=%d&offset=%d", c.baseURL, sellerID, limit, offset)
		req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
		if err != nil {
			return nil, err
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("meli client: items search failed: %w", err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return nil, domain.ErrTokenExpired
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("meli client: items search status %d: %s", resp.StatusCode, string(body))
		}
		var searchResp meliItemsSearchResponse
		if err := json.Unmarshal(body, &searchResp); err != nil {
			return nil, fmt.Errorf("meli client: parsing items search: %w", err)
		}
		itemIDs = append(itemIDs, searchResp.Results...)
		offset += limit
		if len(searchResp.Results) < limit || offset >= searchResp.Paging.Total {
			break
		}
		if offset >= 1000 {
			break
		}
	}

	products := make([]domain.MeliProduct, 0, len(itemIDs))
	for i := 0; i < len(itemIDs); i += 20 {
		end := i + 20
		if end > len(itemIDs) {
			end = len(itemIDs)
		}
		batch := itemIDs[i:end]
		details, err := c.multigetItems(ctx, accessToken, batch)
		if err != nil {
			return nil, err
		}
		for _, d := range details {
			products = append(products, itemToProducts(d)...)
		}
	}

	c.resolveCategoryNames(ctx, accessToken, products)

	return products, nil
}

func combinacion(v meliItemVariation, nombres ...string) string {
	for _, wanted := range nombres {
		for _, c := range v.AttributeCombinations {
			if strings.EqualFold(c.Name, wanted) && strings.TrimSpace(c.ValueName) != "" {
				return c.ValueName
			}
		}
	}
	return ""
}

func itemToProducts(d meliItemDetail) []domain.MeliProduct {
	brand := atributo(d.Attributes, "BRAND")

	if len(d.Variations) == 0 {
		return []domain.MeliProduct{{
			ID:            d.ID,
			LogisticType:  d.Shipping.LogisticType,
			SKU:           extractSKU(d),
			Barcode:       extractBarcode(d),
			Name:          d.Title,
			Brand:         brand,
			CategoryID:    d.CategoryID,
			Color:         atributo(d.Attributes, "COLOR", "MAIN_COLOR"),
			Size:          atributo(d.Attributes, "SIZE"),
			Price:         d.Price,
			StockQuantity: d.AvailableQuantity,
			ImageURL:      itemImageURL(d),
		}}
	}

	image := itemImageURL(d)
	products := make([]domain.MeliProduct, 0, len(d.Variations))
	for _, v := range d.Variations {
		products = append(products, domain.MeliProduct{
			ID:            d.ID,
			VariationID:   fmt.Sprintf("%d", v.ID),
			ParentName:    d.Title,
			VariantLabel:  variationLabel(v),
			LogisticType:  d.Shipping.LogisticType,
			SKU:           extractSKUFrom(v.SellerCustomField, v.Attributes),
			Barcode:       extractBarcodeFrom(v.Attributes, d.Attributes),
			Name:          variationName(d.Title, v),
			Brand:         brand,
			CategoryID:    d.CategoryID,
			Color:         firstNonBlank(combinacion(v, "Color"), atributo(v.Attributes, "COLOR", "MAIN_COLOR")),
			Size:          firstNonBlank(combinacion(v, "Talle", "Talla", "Size"), atributo(v.Attributes, "SIZE")),
			Price:         d.Price,
			StockQuantity: v.AvailableQuantity,
			ImageURL:      image,
		})
	}
	return products
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (c *MeliClient) multigetItems(ctx context.Context, accessToken string, ids []string) ([]meliItemDetail, error) {
	endpoint := fmt.Sprintf("%s/items?ids=%s&include_attributes=all", c.baseURL, strings.Join(ids, ","))
	req, err := c.newAuthorizedRequest(ctx, http.MethodGet, endpoint, accessToken)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("meli client: multiget failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meli client: multiget status %d: %s", resp.StatusCode, string(body))
	}
	var items []meliMultigetItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("meli client: parsing multiget: %w", err)
	}
	out := make([]meliItemDetail, 0, len(items))
	for _, it := range items {
		if it.Code == http.StatusOK && it.Body.ID != "" {
			out = append(out, it.Body)
		}
	}
	return out, nil
}
