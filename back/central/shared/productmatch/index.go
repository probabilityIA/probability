package productmatch

type Values map[string]string

func (v Values) Get(field string) string {
	if v == nil {
		return ""
	}
	return Normalize(v[field])
}

type Index struct {
	rules       []Rule
	keyFields   []string
	probeFields []string
	maps        []map[string]int
}

func newIndex(rules []Rule, items []Values, keyOf func(Rule) string, probeOf func(Rule) string) *Index {
	rules = Sanitize(rules)
	idx := &Index{
		rules:       rules,
		keyFields:   make([]string, len(rules)),
		probeFields: make([]string, len(rules)),
		maps:        make([]map[string]int, len(rules)),
	}
	for i, r := range rules {
		idx.keyFields[i] = keyOf(r)
		idx.probeFields[i] = probeOf(r)
		idx.maps[i] = make(map[string]int, len(items))
	}
	for pos, item := range items {
		for i := range rules {
			key := item.Get(idx.keyFields[i])
			if key == "" {
				continue
			}
			if _, exists := idx.maps[i][key]; exists {
				idx.maps[i][key] = -1
				continue
			}
			idx.maps[i][key] = pos
		}
	}
	return idx
}

func IndexProbability(rules []Rule, items []Values) *Index {
	return newIndex(rules, items,
		func(r Rule) string { return r.Probability },
		func(r Rule) string { return r.Channel },
	)
}

func IndexChannel(rules []Rule, items []Values) *Index {
	return newIndex(rules, items,
		func(r Rule) string { return r.Channel },
		func(r Rule) string { return r.Probability },
	)
}

func (i *Index) Rules() []Rule {
	if i == nil {
		return DefaultRules()
	}
	return i.rules
}

func (i *Index) Find(probe Values) (int, Rule, bool) {
	if i == nil {
		return -1, Rule{}, false
	}
	for n := range i.rules {
		value := probe.Get(i.probeFields[n])
		if value == "" {
			continue
		}
		pos, ok := i.maps[n][value]
		if !ok || pos < 0 {
			continue
		}
		return pos, i.rules[n], true
	}
	return -1, Rule{}, false
}

func (i *Index) HasAnyKey(probe Values) bool {
	if i == nil {
		return false
	}
	for n := range i.rules {
		if probe.Get(i.probeFields[n]) != "" {
			return true
		}
	}
	return false
}
