package semantic

func ResolveCondition(kondisi string) (string, bool) {
	kondisiMap := map[string]string{
		"rusak":     "status.phase!=Running",
		"terdampar": "status.phase=Pending",
		"bocor":     "status.reason=OOMKilled",
		"sehat":     "status.phase=Running",
	}

	filter, ok := kondisiMap[kondisi]
	return filter, ok
}
