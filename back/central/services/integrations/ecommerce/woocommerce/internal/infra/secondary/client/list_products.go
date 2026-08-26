package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

type wooFlexString string

func (f *wooFlexString) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "null" {
		*f = ""
		return nil
	}
	if len(trimmed) > 0 && trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = wooFlexString(s)
		return nil
	}
	*f = wooFlexString(trimmed)
	return nil
}

type wooImageResponse struct {
	Src string `json:"src"`
}

type wooTermResponse struct {
	Name string `json:"name"`
}

type wooDimensionsResponse struct {
	Length wooFlexString `json:"length"`
	Width  wooFlexString `json:"width"`
	Height wooFlexString `json:"height"`
}

type wooProductResponse struct {
	ID               int64                 `json:"id"`
	Name             string                `json:"name"`
	SKU              string                `json:"sku"`
	GlobalUniqueID   string                `json:"global_unique_id"`
	Type             string                `json:"type"`
	Price            wooFlexString         `json:"price"`
	RegularPrice     wooFlexString         `json:"regular_price"`
	Description      string                `json:"description"`
	ShortDescription string                `json:"short_description"`
	Categories       []wooTermResponse     `json:"categories"`
	Brands           []wooTermResponse     `json:"brands"`
	Weight           wooFlexString         `json:"weight"`
	Dimensions       wooDimensionsResponse `json:"dimensions"`
	StockQuantity    *int                  `json:"stock_quantity"`
	Images           []wooImageResponse    `json:"images"`
}

type wooAttributeResponse struct {
	Name   string `json:"name"`
	Option string `json:"option"`
}

type wooVariationResponse struct {
	ID             int64                  `json:"id"`
	Attributes     []wooAttributeResponse `json:"attributes"`
	SKU            string                 `json:"sku"`
	GlobalUniqueID string                 `json:"global_unique_id"`
	Price          wooFlexString          `json:"price"`
	RegularPrice   wooFlexString          `json:"regular_price"`
	Description    string                 `json:"description"`
	Weight         wooFlexString          `json:"weight"`
	Dimensions     wooDimensionsResponse  `json:"dimensions"`
	StockQuantity  *int                   `json:"stock_quantity"`
	Image          *wooImageResponse      `json:"image"`
}

func firstImage(images []wooImageResponse) string {
	if len(images) == 0 {
		return ""
	}
	return images[0].Src
}

func atributosDeVariacion(attrs []wooAttributeResponse) map[string]string {
	out := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		nombre := strings.ToLower(strings.TrimSpace(attr.Name))
		valor := strings.TrimSpace(attr.Option)
		if nombre == "" || valor == "" {
			continue
		}
		out[nombre] = valor
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func etiquetaDeVariacion(attrs []wooAttributeResponse) string {
	partes := make([]string, 0, len(attrs))
	for _, attr := range attrs {
		if valor := strings.TrimSpace(attr.Option); valor != "" {
			partes = append(partes, valor)
		}
	}
	return strings.Join(partes, " / ")
}

func firstTerm(terms []wooTermResponse) string {
	if len(terms) == 0 {
		return ""
	}
	return terms[0].Name
}

func numeroWoo(value string) *float64 {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil || parsed <= 0 {
		return nil
	}
	return &parsed
}

func firstText(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (c *WooCommerceClient) GetProducts(ctx context.Context, storeURL, consumerKey, consumerSecret string) ([]domain.WooProduct, error) {
	storeURL = strings.TrimRight(storeURL, "/")

	products := make([]domain.WooProduct, 0)
	page := 1
	for {
		endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products?per_page=100&page=%d&status=any", storeURL, page)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("woocommerce client: creating request: %w", err)
		}
		req.SetBasicAuth(consumerKey, consumerSecret)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("woocommerce client: request failed: %w", err)
		}

		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			resp.Body.Close()
			return nil, domain.ErrInvalidCredentials
		}
		if resp.StatusCode != http.StatusOK {
			raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
			return nil, fmt.Errorf("woocommerce client: unexpected status %d listing products: %s", resp.StatusCode, string(raw))
		}

		var list []wooProductResponse
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("woocommerce client: decoding products response: %w", err)
		}
		resp.Body.Close()

		if len(list) == 0 {
			break
		}

		for _, p := range list {
			if p.Type == "variable" {
				variations, verr := c.getProductVariations(ctx, storeURL, consumerKey, consumerSecret, p.ID)
				if verr != nil {
					return nil, verr
				}
				parentID := strconv.FormatInt(p.ID, 10)
				for _, v := range variations {
					vprice := 0.0
					if v.Price != "" {
						vprice, _ = strconv.ParseFloat(string(v.Price), 64)
					}
					vstock := 0
					if v.StockQuantity != nil {
						vstock = *v.StockQuantity
					}
					image := firstImage(p.Images)
					if v.Image != nil && v.Image.Src != "" {
						image = v.Image.Src
					}
					products = append(products, domain.WooProduct{
						ID:                strconv.FormatInt(v.ID, 10),
						ParentID:          parentID,
						ParentName:        p.Name,
						VariantLabel:      etiquetaDeVariacion(v.Attributes),
						VariantAttributes: atributosDeVariacion(v.Attributes),
						SKU:               v.SKU,
						Barcode:           v.GlobalUniqueID,
						Name:              p.Name,
						Description:       firstText(v.Description, p.Description, p.ShortDescription),
						Category:          firstTerm(p.Categories),
						Brand:             firstTerm(p.Brands),
						Weight:            numeroWoo(firstText(string(v.Weight), string(p.Weight))),
						Length:            numeroWoo(firstText(string(v.Dimensions.Length), string(p.Dimensions.Length))),
						Width:             numeroWoo(firstText(string(v.Dimensions.Width), string(p.Dimensions.Width))),
						Height:            numeroWoo(firstText(string(v.Dimensions.Height), string(p.Dimensions.Height))),
						Price:             vprice,
						StockQuantity:     vstock,
						ImageURL:          image,
					})
				}
				continue
			}

			price := 0.0
			if p.Price != "" {
				price, _ = strconv.ParseFloat(string(p.Price), 64)
			}
			stock := 0
			if p.StockQuantity != nil {
				stock = *p.StockQuantity
			}
			products = append(products, domain.WooProduct{
				ID:            strconv.FormatInt(p.ID, 10),
				SKU:           p.SKU,
				Barcode:       p.GlobalUniqueID,
				Name:          p.Name,
				Description:   firstText(p.Description, p.ShortDescription),
				Category:      firstTerm(p.Categories),
				Brand:         firstTerm(p.Brands),
				Weight:        numeroWoo(string(p.Weight)),
				Length:        numeroWoo(string(p.Dimensions.Length)),
				Width:         numeroWoo(string(p.Dimensions.Width)),
				Height:        numeroWoo(string(p.Dimensions.Height)),
				Price:         price,
				StockQuantity: stock,
				ImageURL:      firstImage(p.Images),
			})
		}

		if len(list) < 100 {
			break
		}
		page++
	}

	return products, nil
}

func (c *WooCommerceClient) getProductVariations(ctx context.Context, storeURL, consumerKey, consumerSecret string, parentID int64) ([]wooVariationResponse, error) {
	all := make([]wooVariationResponse, 0)
	page := 1
	for {
		endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products/%d/variations?per_page=100&page=%d", storeURL, parentID, page)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("woocommerce client: creating variations request: %w", err)
		}
		req.SetBasicAuth(consumerKey, consumerSecret)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("woocommerce client: variations request failed: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
			resp.Body.Close()
			return nil, fmt.Errorf("woocommerce client: unexpected status %d listing variations of %d: %s", resp.StatusCode, parentID, string(raw))
		}

		var list []wooVariationResponse
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("woocommerce client: decoding variations response: %w", err)
		}
		resp.Body.Close()

		if len(list) == 0 {
			break
		}
		all = append(all, list...)
		if len(list) < 100 {
			break
		}
		page++
	}
	return all, nil
}
