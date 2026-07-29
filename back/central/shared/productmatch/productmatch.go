package productmatch

import (
	"encoding/json"
	"strings"
)

const (
	FieldSKU        = "sku"
	FieldBarcode    = "barcode"
	FieldExternalID = "external_id"
	FieldVariantID  = "variant_id"
	FieldName       = "name"
)

type Rule struct {
	Probability string `json:"probability"`
	Channel     string `json:"channel"`
}

type Options struct {
	Probability []string `json:"probability"`
	Channel     []string `json:"channel"`
}

var probabilityFields = map[string]bool{
	FieldSKU:        true,
	FieldBarcode:    true,
	FieldExternalID: true,
	FieldName:       true,
}

var channelFields = map[string]bool{
	FieldSKU:        true,
	FieldBarcode:    true,
	FieldExternalID: true,
	FieldVariantID:  true,
	FieldName:       true,
}

func DefaultRules() []Rule {
	return []Rule{{Probability: FieldSKU, Channel: FieldSKU}}
}

func Normalize(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func (r Rule) Valid() bool {
	return probabilityFields[r.Probability] && channelFields[r.Channel]
}

func (r Rule) Key() string {
	return r.Probability + "->" + r.Channel
}

func Sanitize(rules []Rule) []Rule {
	out := make([]Rule, 0, len(rules))
	seen := make(map[string]bool, len(rules))
	for _, r := range rules {
		r.Probability = Normalize(r.Probability)
		r.Channel = Normalize(r.Channel)
		if !r.Valid() || seen[r.Key()] {
			continue
		}
		seen[r.Key()] = true
		out = append(out, r)
	}
	if len(out) == 0 {
		return DefaultRules()
	}
	return out
}

func Parse(raw []byte) []Rule {
	if len(raw) == 0 {
		return nil
	}
	var rules []Rule
	if err := json.Unmarshal(raw, &rules); err != nil || len(rules) == 0 {
		return nil
	}
	sanitized := make([]Rule, 0, len(rules))
	seen := make(map[string]bool, len(rules))
	for _, r := range rules {
		r.Probability = Normalize(r.Probability)
		r.Channel = Normalize(r.Channel)
		if !r.Valid() || seen[r.Key()] {
			continue
		}
		seen[r.Key()] = true
		sanitized = append(sanitized, r)
	}
	if len(sanitized) == 0 {
		return nil
	}
	return sanitized
}

func Resolve(integrationRaw, typeDefaultRaw []byte) []Rule {
	if rules := Parse(integrationRaw); rules != nil {
		return rules
	}
	if rules := Parse(typeDefaultRaw); rules != nil {
		return rules
	}
	return DefaultRules()
}

func Marshal(rules []Rule) ([]byte, error) {
	return json.Marshal(Sanitize(rules))
}

func ParseOptions(raw []byte) Options {
	opts := Options{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &opts)
	}
	if len(opts.Probability) == 0 {
		opts.Probability = []string{FieldSKU, FieldBarcode, FieldExternalID}
	}
	if len(opts.Channel) == 0 {
		opts.Channel = []string{FieldSKU, FieldBarcode, FieldExternalID, FieldVariantID}
	}
	return opts
}

func AllowedByOptions(rules []Rule, opts Options) []Rule {
	prob := make(map[string]bool, len(opts.Probability))
	for _, f := range opts.Probability {
		prob[Normalize(f)] = true
	}
	ch := make(map[string]bool, len(opts.Channel))
	for _, f := range opts.Channel {
		ch[Normalize(f)] = true
	}
	out := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if prob[r.Probability] && ch[r.Channel] {
			out = append(out, r)
		}
	}
	if len(out) == 0 {
		return DefaultRules()
	}
	return out
}
