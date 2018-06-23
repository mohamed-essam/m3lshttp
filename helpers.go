package m3lshttp

func handlerMapHasKey(mp map[string]handler, key string) bool {
	if _, ok := mp[key]; ok {
		return true
	}
	return false
}

func mapToInterfaceMap(mp map[string][]string) (ret map[string]interface{}) {
	ret = make(map[string]interface{}, 0)
	for k, v := range mp {
		ret[k] = interface{}(v[0])
	}
	return ret
}
