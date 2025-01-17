package types

func NewFilterLogger(filters [][]byte) FilterLogger {
	return FilterLogger{
		Filters: filters,
	}
}

func (flog *FilterLogger) Contains(filter []byte) bool {
	for _, f := range flog.Filters {
		if string(f) == string(filter) {
			return true
		}
	}
	return false
}

func (flog *FilterLogger) Add(filter []byte) {
	if !flog.Contains(filter) {
		flog.Filters = append(flog.Filters, filter)
	}
}

func (flog *FilterLogger) Delete(filter []byte) {
	for i, f := range flog.Filters {
		if string(f) == string(filter) {
			flog.Filters = append(flog.Filters[:i], flog.Filters[i+1:]...)
			return
		}
	}
}
