package saol

import "strings"

// Small generic helpers for a functional style.

func mapi[T, U any](xs []T, f func(int, T) U) []U {
	ys := make([]U, len(xs))
	for i, x := range xs {
		ys[i] = f(i, x)
	}
	return ys
}

func mapf[T, U any](xs []T, f func(T) U) []U {
	return mapi(xs, func(_ int, x T) U { return f(x) })
}

func flatMap[T, U any](xs []T, f func(T) []U) []U {
	var ys []U
	for _, x := range xs {
		ys = append(ys, f(x)...)
	}
	return ys
}

// mapErr is mapf for fallible functions; it stops at the first error.
func mapErr[T, U any](xs []T, f func(T) (U, error)) ([]U, error) {
	ys := make([]U, 0, len(xs))
	for _, x := range xs {
		y, err := f(x)
		if err != nil {
			return nil, err
		}
		ys = append(ys, y)
	}
	return ys, nil
}

func filter[T any](xs []T, keep func(T) bool) []T {
	var ys []T
	for _, x := range xs {
		if keep(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

// uniqueBy keeps the first element for each key, preserving order.
func uniqueBy[T any, K comparable](xs []T, key func(T) K) []T {
	seen := map[K]bool{}
	return filter(xs, func(x T) bool {
		k := key(x)
		dup := seen[k]
		seen[k] = true
		return !dup
	})
}

func nonEmpty(s string) bool { return s != "" }

// join concatenates the non-empty parts with sep.
func join(sep string, parts ...string) string {
	return strings.Join(filter(parts, nonEmpty), sep)
}

// wrap encloses s in l and r, unless s is empty.
func wrap(l, s, r string) string {
	if s == "" {
		return ""
	}
	return l + s + r
}
