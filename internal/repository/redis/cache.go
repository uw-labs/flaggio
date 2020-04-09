package redis

func shouldCacheFindAll(search *string, offset, limit *int64) bool {
	if search == nil && offset == nil && limit == nil {
		return true
	}
	return false
}
