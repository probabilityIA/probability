package productmatch

import "fmt"

type DataFieldSQL struct {
	Key     string
	Label   string
	Note    string
	Column  string
	Cast    string
	Channel string
	Own     string
	Extra   string
}

func (f DataFieldSQL) FillCondition() string {
	return fmt.Sprintf("%s <> '' AND %s = '' %s", f.Channel, f.Own, f.Extra)
}

func (f DataFieldSQL) ReplaceCondition() string {
	return fmt.Sprintf("%s <> '' AND %s <> '' AND %s <> %s %s", f.Channel, f.Own, f.Channel, f.Own, f.Extra)
}

func DataFieldsSQL() []DataFieldSQL {
	return []DataFieldSQL{
		{
			Key: "description", Label: "Descripcion", Column: "description",
			Channel: `COALESCE(m.channel_snapshot->>'description','')`, Own: `COALESCE(p.description,'')`,
		},
		{
			Key: "category", Label: "Categoria", Column: "category",
			Channel: `COALESCE(m.channel_snapshot->>'category','')`, Own: `COALESCE(p.category,'')`,
		},
		{
			Key: "brand", Label: "Marca", Column: "brand",
			Channel: `COALESCE(m.channel_snapshot->>'brand','')`, Own: `COALESCE(p.brand,'')`,
		},
		{
			Key: "image_url", Label: "Imagen", Column: "image_url",
			Note:    "se guarda la URL del canal, no se copia el archivo",
			Channel: `COALESCE(m.channel_snapshot->>'image_url','')`, Own: `COALESCE(p.image_url,'')`,
		},
		{
			Key: "barcode", Label: "Codigo de barras", Column: "barcode",
			Note:    "se ignora cuando el canal repite el SKU",
			Channel: `COALESCE(m.channel_snapshot->>'barcode','')`, Own: `COALESCE(p.barcode,'')`,
			Extra: `AND UPPER(TRIM(COALESCE(m.channel_snapshot->>'barcode',''))) <> UPPER(TRIM(COALESCE(p.sku,'')))`,
		},
		{
			Key: "variant_label", Label: "Variante", Column: "variant_label",
			Note:    "color y talla del canal",
			Channel: `TRIM(BOTH ' /' FROM CONCAT_WS(' / ', NULLIF(m.channel_snapshot->>'color',''), NULLIF(m.channel_snapshot->>'size','')))`,
			Own:     `COALESCE(p.variant_label,'')`,
		},
		{
			Key: "weight", Label: "Peso", Column: "weight", Cast: "::numeric",
			Channel: `CASE WHEN m.channel_snapshot->>'weight' ~ '^[0-9]+(\.[0-9]+){0,1}$' THEN ROUND((m.channel_snapshot->>'weight')::numeric,3)::text ELSE '' END`,
			Own:     `COALESCE(ROUND(p.weight::numeric,3)::text,'')`,
		},
		{
			Key: "name", Label: "Nombre", Column: "name",
			Channel: `COALESCE(m.channel_snapshot->>'name','')`, Own: `COALESCE(p.name,'')`,
		},
	}
}

func DataFieldByKey(key string) (DataFieldSQL, bool) {
	for _, f := range DataFieldsSQL() {
		if f.Key == key {
			return f, true
		}
	}
	return DataFieldSQL{}, false
}
