package shared

// CollatzStoppingTimeF returns the result of the following recursive function f : N → N, as well as
// the standard total stopping time, as well as the largest value of x passed to f during the recursion.
//
// f(x) = { 0                 if x = 1
// .... = { 1 + f(C(x))       if x != 1
//
// # See README for explanation of C
//
// It returns the total stopping time twice to match the signature of the optimized functions, as well as the largest
// value of x passed to f recursively.
//
// https://oeis.org/A006577
func CollatzStoppingTimeF(n uint64) (uint64, uint64, uint64) {
	time := uint64(0)
	maxN := n
	for n != 1 {
		if n&1 == 1 {
			n = n<<1 + n + 1 // 3n+1
		} else {
			n >>= 1 // n/2
		}
		if n > maxN {
			maxN = n
		}
		time++
	}
	return time, time, maxN
}

// CollatzStoppingTimeG returns the result of the following recursive function g : N → N, as well as
// the standard total stopping time, as well as the largest value of x passed to g during the recursion.
//
// g(x) = { 0                  if x = 1
// .... = { 1 + g(x/2^m)       if x ≡ 0 (mod 2)
// .... = { 1 + g(R(x))        if x ≡ 1 (mod 2) and x ≠ 1
//
// Where $R$ is defined in the README and 2^m is highest power of 2 that divides x
//
// https://oeis.org/A286380
func CollatzStoppingTimeG(n uint64) (uint64, uint64, uint64) {
	maxN := n
	reducedTime, normalTime := uint64(0), uint64(0)
	// the main loop assumes we have an odd number
	if n&1 == 0 {
		for n&1 == 0 { // n/2^m
			n >>= 1
			normalTime++
		}
		reducedTime++
	}
	for n != 1 {
		if n > maxN {
			maxN = n
		}
		n = n<<1 + n + 1 // 3n+1
		normalTime++
		for n&1 == 0 { // n/2^m
			n >>= 1
			normalTime++
		}
		reducedTime++
	}
	return reducedTime, normalTime, maxN
}

// CollatzStoppingTimeH returns the result of the following recursive function h : N → N, as well as
// the standard total stopping time, as well as the largest value of x passed to h during the recursion.
//
// h(x) = { 0              if x = 1
// .... = { 1 + h(x/2^m)   if x ≡ 0 (mod 2)
// .... = { 1 + h(R'(x))   if x ≡ 1 (mod 2) and x ≠ 1
//
// Where R' is defined in the README and 2^m is the highest power of 2 that divides x
//
// https://oeis.org/A160541
func CollatzStoppingTimeH(n uint64) (uint64, uint64, uint64) {
	maxN := n
	reducedTime, normalTime := uint64(0), uint64(0)
	// the main loop assumes we have an odd number
	if n&1 == 0 {
		for n&1 == 0 { // n/2^m
			n >>= 1
			normalTime++
		}
		reducedTime++
	}
	for n != 1 {
		if n > maxN {
			maxN = n
		}
		// variation of ((3x+1)/2)+1 that doesn't get as large
		n = ((n - 1) >> 1) + 1
		x := n
		n = (n << 1) + n
		normalTime += 2
		// multiply n by (3/2) until x is no longer divisible by 2
		for x&1 == 0 {
			x >>= 1
			n = (n >> 1) + n
			normalTime += 2
		}
		n--
		for n&1 == 0 { // n/2^m
			n >>= 1
			normalTime++
		}
		reducedTime++
		if n > maxN {
			maxN = n
		}
	}
	return reducedTime, normalTime, maxN
}
