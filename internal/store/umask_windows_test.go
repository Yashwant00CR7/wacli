//go:build windows

package store

func setUmask(mask int) int {
	return 0
}
