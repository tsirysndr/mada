package graph

func ParseSkipAndSize(skip *int, size *int) (int, int) {
	var (
		skipInt = 0
		sizeInt = 100
	)
	if skip != nil {
		skipInt = *skip
	}
	if size != nil {
		sizeInt = *size
	}
	return skipInt, sizeInt
}
