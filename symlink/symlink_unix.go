// +build !windows

package symlink

func evalSymlinks(path string) (string, error) {
	return walkSymlinks(path)
}
