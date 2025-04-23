package library

var validTenors = map[int]struct{}{
	1: {}, 2: {}, 3: {}, 6: {},
}

func IsValidTenor(tenor int) bool {
	_, ok := validTenors[tenor]

	return ok
}
