package productmatch

import "strings"

// TypoCandidate es un SKU que quedo sin pareja de alguno de los dos lados.
type TypoCandidate struct {
	SKU      string
	Name     string
	Quantity int
	HasQty   bool
	Ref      string
	Label    string
}

// TypoSuspect es un par que probablemente sea el mismo producto escrito distinto.
type TypoSuspect struct {
	ChannelSKU     string
	ChannelName    string
	ChannelRef     string
	ChannelLabel   string
	ChannelQty     int
	HasChannelQty  bool
	ProbabilitySKU string
	ProbabilityQty int
	HasProbQty     bool
	Reason         string
}

const (
	// ReasonTransposed son dos caracteres contiguos intercambiados (SDH <-> SHD).
	// Es muy especifico: casi nunca ocurre por casualidad.
	ReasonTransposed = "transposed"
	// ReasonOneCharAndQty es un caracter de diferencia, pero solo se acepta si la
	// cantidad tambien coincide. Un caracter por si solo no alcanza: con SKUs
	// cortos casi todo queda a distancia 1 de algo.
	ReasonOneCharAndQty = "one_char_and_qty"
)

func typoKey(s string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(strings.TrimSpace(s)) {
		if r == ' ' || r == '\t' || r == '_' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func isTransposition(a, b string) bool {
	if len(a) != len(b) || a == b {
		return false
	}
	diff := make([]int, 0, 3)
	for i := range a {
		if a[i] != b[i] {
			diff = append(diff, i)
			if len(diff) > 2 {
				return false
			}
		}
	}
	if len(diff) != 2 || diff[1] != diff[0]+1 {
		return false
	}
	return a[diff[0]] == b[diff[1]] && a[diff[1]] == b[diff[0]]
}

// editDistance corta apenas supera el tope, no hace falta el valor exacto.
func editDistance(a, b string, cap int) int {
	if a == b {
		return 0
	}
	if len(a)-len(b) > cap || len(b)-len(a) > cap {
		return cap + 1
	}
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		best := cur[0]
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min3(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
			if cur[j] < best {
				best = cur[j]
			}
		}
		if best > cap {
			return cap + 1
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

func min3(a, b, c int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}

// sameQuantity solo confirma si hay algo que confirmar: dos ceros no son
// evidencia, son el valor por defecto de todo lo que nadie ha tocado.
func sameQuantity(a, b TypoCandidate) bool {
	if !a.HasQty || !b.HasQty {
		return false
	}
	if a.Quantity <= 0 && b.Quantity <= 0 {
		return false
	}
	d := a.Quantity - b.Quantity
	if d < 0 {
		d = -d
	}
	return d <= 1
}

// DetectTypos busca pares que probablemente sean el mismo producto escrito
// distinto en cada sistema. Solo devuelve los que tienen dos senales, porque una
// sola produce una sugerencia para casi todo y el listado deja de servir.
// Un SKU del canal se reporta una sola vez, con su mejor candidato.
func DetectTypos(channel, probability []TypoCandidate) []TypoSuspect {
	suspects := make([]TypoSuspect, 0)
	if len(channel) == 0 || len(probability) == 0 {
		return suspects
	}

	probKeys := make([]string, len(probability))
	for i, p := range probability {
		probKeys[i] = typoKey(p.SKU)
	}

	for _, c := range channel {
		ck := typoKey(c.SKU)
		if ck == "" {
			continue
		}
		best := -1
		bestReason := ""
		for i, p := range probability {
			pk := probKeys[i]
			if pk == "" || pk == ck {
				continue
			}
			switch {
			case isTransposition(ck, pk):
				best, bestReason = i, ReasonTransposed
			case bestReason != ReasonTransposed && editDistance(ck, pk, 1) == 1 && sameQuantity(c, p):
				best, bestReason = i, ReasonOneCharAndQty
			}
			if bestReason == ReasonTransposed {
				break
			}
		}
		if best < 0 {
			continue
		}
		p := probability[best]
		suspects = append(suspects, TypoSuspect{
			ChannelSKU:     c.SKU,
			ChannelName:    c.Name,
			ChannelRef:     c.Ref,
			ChannelLabel:   c.Label,
			ChannelQty:     c.Quantity,
			HasChannelQty:  c.HasQty,
			ProbabilitySKU: p.SKU,
			ProbabilityQty: p.Quantity,
			HasProbQty:     p.HasQty,
			Reason:         bestReason,
		})
	}
	return suspects
}
