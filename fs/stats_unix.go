// +build !windows

package fs

import (
	"os"

	sys "golang.org/x/sys/unix"
)

// Stats hold the file system statistics, in bytes. Only the Usage is in
// percentage.
type Stats struct {
	Free      uint64
	Available uint64
	Size      uint64
	Used      uint64
	Usage     float64
}

// GetStats returns an object holding the file system usage of path.
// If the path is empty, it uses the current working directory.
func GetStats(path string) (*Stats, error) {
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

	var stats = &Stats{
		Free:      st.Bfree * uint64(st.Bsize),
		Available: st.Bavail * uint64(st.Bsize),
		Size:      st.Blocks * uint64(st.Bsize),
	}
	stats.Used = stats.Size - stats.Free
	stats.Usage = float64(stats.Used) / float64(stats.Size)

	return stats, nil
}
