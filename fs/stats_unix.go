// +build !windows

package fs

import (
	"os"

	sys "golang.org/x/sys/unix"
)

// Info holds the filesystem statistics. All the fields are in bytes, except for
// Usage which is a percentage.
type Info struct {
	Free      uint64
	Available uint64
	Size      uint64
	Used      uint64
	Usage     float64
}

// Stat returns the Info structure describing filesystem usage. If the path is
// empty, it uses the current working directory.
func Stat(path string) (*Info, error) {
	if path == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		path = dir
	}

	var st sys.Statfs_t
	if err := sys.Statfs(path, &st); err != nil {
		return nil, err
	}

	stats := &Info{
		Free:      st.Bfree * uint64(st.Bsize),
		Available: st.Bavail * uint64(st.Bsize),
		Size:      st.Blocks * uint64(st.Bsize),
	}
	stats.Used = stats.Size - stats.Free
	stats.Usage = float64(stats.Used) / float64(stats.Size)

	return stats, nil
}
