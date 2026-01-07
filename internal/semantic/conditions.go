package semantic

func ResolveCondition(kondisi string) (string, bool) {
	kondisiMap := map[string]string{
		"rusak":     "status!=Running",
		"terdampar": "status=Pending",
		"bocor":     "reason=OOMKilled",
		"sehat":     "status=Running",
	}

	filter, ok := kondisiMap[kondisi]
	return filter, ok
}
