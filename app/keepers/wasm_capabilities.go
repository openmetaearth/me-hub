package keepers

func filterWasmCapabilities(capabilities []string, disabled ...string) []string {
	disabledSet := make(map[string]struct{}, len(disabled))
	for _, capability := range disabled {
		disabledSet[capability] = struct{}{}
	}

	filtered := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		if _, ok := disabledSet[capability]; ok {
			continue
		}
		filtered = append(filtered, capability)
	}
	return filtered
}
