// +build !windows

package fs

func evalSymlink(path string) (string, error) {
	return walkSymlinks(path)
}
