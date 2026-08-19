package usecases

const maxReportedFailedSKUs = 50

type failedSKUs struct {
	skus  []string
	total int
}

func (f *failedSKUs) add(sku string) {
	f.total++
	if sku == "" || len(f.skus) >= maxReportedFailedSKUs {
		return
	}
	f.skus = append(f.skus, sku)
}

func (f *failedSKUs) count() int {
	return f.total
}

func (f *failedSKUs) list() []string {
	if f.skus == nil {
		return []string{}
	}
	return f.skus
}

func (f *failedSKUs) truncated() int {
	hidden := f.total - len(f.skus)
	if hidden < 0 {
		return 0
	}
	return hidden
}
