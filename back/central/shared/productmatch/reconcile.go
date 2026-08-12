package productmatch

type Item struct {
	SKU        string
	Barcode    string
	ExternalID string
	VariantID  string
	Name       string
}

func (i Item) Values() Values {
	return Values{
		FieldSKU:        i.SKU,
		FieldBarcode:    i.Barcode,
		FieldExternalID: i.ExternalID,
		FieldVariantID:  i.VariantID,
		FieldName:       i.Name,
	}
}

func toValues(items []Item) []Values {
	out := make([]Values, len(items))
	for i, it := range items {
		out[i] = it.Values()
	}
	return out
}

type ExternalRefs struct {
	ProductID    string
	VariantID    string
	SKU          string
	Barcode      string
	LogisticType string
}

type Pair struct {
	ProbabilityIndex int
	ChannelIndex     int
	Rule             Rule
}

type Outcome struct {
	Rules                  []Rule
	Pairs                  []Pair
	OnlyInProbability      []int
	OnlyInChannel          []int
	ProbabilityUnmatchable int
	ChannelUnmatchable     int

	ProbabilityNoKey []int
	ChannelNoKey     []int
}

func Reconcile(rules []Rule, probability, channel []Item) Outcome {
	rules = Sanitize(rules)
	probValues := toValues(probability)
	chanValues := toValues(channel)

	index := IndexProbability(rules, probValues)
	channelProbe := IndexChannel(rules, chanValues)

	out := Outcome{
		Rules:             rules,
		Pairs:             []Pair{},
		OnlyInProbability: []int{},
		OnlyInChannel:     []int{},
		ProbabilityNoKey:  []int{},
		ChannelNoKey:      []int{},
	}

	matchedProb := make(map[int]bool, len(probability))
	for ci, cv := range chanValues {
		if !index.HasAnyKey(cv) {
			out.ChannelUnmatchable++
			out.ChannelNoKey = append(out.ChannelNoKey, ci)
			continue
		}
		pi, rule, ok := index.Find(cv)
		if !ok {
			out.OnlyInChannel = append(out.OnlyInChannel, ci)
			continue
		}
		matchedProb[pi] = true
		out.Pairs = append(out.Pairs, Pair{ProbabilityIndex: pi, ChannelIndex: ci, Rule: rule})
	}

	for pi, pv := range probValues {
		if matchedProb[pi] {
			continue
		}
		if !channelProbe.HasAnyKey(pv) {
			out.ProbabilityUnmatchable++
			out.ProbabilityNoKey = append(out.ProbabilityNoKey, pi)
			continue
		}
		out.OnlyInProbability = append(out.OnlyInProbability, pi)
	}

	return out
}
