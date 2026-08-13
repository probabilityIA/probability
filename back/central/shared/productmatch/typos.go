package productmatch

import "strings"

type TypoCandidate struct {
	SKU      string
	Name     string
	Quantity int
	HasQty   bool
	Ref      string
	Label    string
}

type TypoSuspect struct {
	ChannelSKU      string
	ChannelName     string
	ChannelRef      string
	ChannelLabel    string
	ChannelQty      int
	HasChannelQty   bool
	ProbabilitySKU  string
	ProbabilityName string
	ProbabilityQty  int
	HasProbQty      bool
	Pattern         string
	FixSide         string
}

const (
	PatternSpacing    = "spacing"
	PatternTransposed = "transposed"
	PatternOneChar    = "one_char"
)

const (
	FixInChannel     = "channel"
	FixInProbability = "probability"
)

const (
	AuthorityUnknown     = ""
	AuthorityProbability = "probability"
	AuthorityChannel     = "channel"
)

type TypoOptions struct {
	Authority string
}

func squeeze(s string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(strings.TrimSpace(s)) {
		if r == ' ' || r == '\t' || r == '_' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func hasBlanks(s string) bool {
	return strings.ContainsAny(strings.TrimSpace(s), " \t_")
}

func splitSKU(squeezed string) (string, string) {
	i := strings.LastIndex(squeezed, "-")
	if i < 0 {
		return squeezed, ""
	}
	return squeezed[:i], squeezed[i+1:]
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

func fixSide(pattern string, channel, prob TypoCandidate, opts TypoOptions) string {
	if pattern == PatternSpacing {
		if hasBlanks(channel.SKU) {
			return FixInChannel
		}
		return FixInProbability
	}
	switch opts.Authority {
	case AuthorityProbability:
		return FixInChannel
	case AuthorityChannel:
		return FixInProbability
	}
	if channel.HasQty && channel.Quantity > 0 && (!prob.HasQty || prob.Quantity <= 0) {
		return FixInProbability
	}
	return FixInChannel
}

type parsedSKU struct {
	key    string
	base   string
	suffix string
}

type bestMatch struct {
	index int
	dup   bool
}

func (b *bestMatch) offer(i int) {
	if b.index >= 0 {
		b.dup = true
		return
	}
	b.index = i
}

func (b bestMatch) usable() bool {
	return b.index >= 0 && !b.dup
}

func DetectTypos(channel, probability []TypoCandidate, opts TypoOptions) []TypoSuspect {
	suspects := make([]TypoSuspect, 0)
	if len(channel) == 0 || len(probability) == 0 {
		return suspects
	}

	parsed := make([]parsedSKU, len(probability))
	for i, p := range probability {
		k := squeeze(p.SKU)
		base, suffix := splitSKU(k)
		parsed[i] = parsedSKU{key: k, base: base, suffix: suffix}
	}

	for _, c := range channel {
		ck := squeeze(c.SKU)
		if ck == "" {
			continue
		}
		cbase, csuffix := splitSKU(ck)

		mismo := bestMatch{index: -1}
		trocado := bestMatch{index: -1}
		unCaracter := bestMatch{index: -1}

		for i, p := range parsed {
			if p.key == "" {
				continue
			}
			switch {
			case p.key == ck:
				if hasBlanks(c.SKU) || hasBlanks(probability[i].SKU) {
					mismo.offer(i)
				}
			case p.suffix != csuffix:
				continue
			case isTransposition(cbase, p.base):
				trocado.offer(i)
			case editDistance(cbase, p.base, 1) == 1 && sameQuantity(c, probability[i]):
				unCaracter.offer(i)
			}
		}

		best, pattern := -1, ""
		switch {
		case mismo.usable():
			best, pattern = mismo.index, PatternSpacing
		case trocado.usable():
			best, pattern = trocado.index, PatternTransposed
		case unCaracter.usable():
			best, pattern = unCaracter.index, PatternOneChar
		}
		if best < 0 {
			continue
		}

		p := probability[best]
		suspects = append(suspects, TypoSuspect{
			ChannelSKU:      c.SKU,
			ChannelName:     c.Name,
			ChannelRef:      c.Ref,
			ChannelLabel:    c.Label,
			ChannelQty:      c.Quantity,
			HasChannelQty:   c.HasQty,
			ProbabilitySKU:  p.SKU,
			ProbabilityName: p.Name,
			ProbabilityQty:  p.Quantity,
			HasProbQty:      p.HasQty,
			Pattern:         pattern,
			FixSide:         fixSide(pattern, c, p, opts),
		})
	}
	return suspects
}
