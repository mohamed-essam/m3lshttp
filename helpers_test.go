package m3lshttp

func urlNodeArrayContains(arr []*urlNode, node *urlNode) bool {
	for _, v := range arr {
		if v == node {
			return true
		}
	}
	return false
}
