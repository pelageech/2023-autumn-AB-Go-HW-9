package slices

// Split splits a source slice into several slices,where all slices,
// except the last one, have the length portion.
//
// Combining all the resulting slices gives the original slice.
//
// Panics if portion <= 0.
func Split[S ~[]E, E any](s S, portion int) []S {
	if portion <= 0 {
		panic("portion is less than zero")
	}

	if len(s) == 0 {
		return []S{}
	}

	ssLen := len(s) / portion
	if len(s)%portion > 0 {
		ssLen++
	}

	ss := make([]S, ssLen)

	for i := 0; i < len(s)/portion; i++ {
		ss[i] = make(S, portion)
		copy(ss[i], s[i*portion:(i+1)*portion])
	}

	if len(s)%portion > 0 {
		ss[len(ss)-1] = make(S, len(s)-(len(ss)-1)*portion)
		copy(ss[len(ss)-1], s[(len(ss)-1)*portion:])
	}

	return ss
}
