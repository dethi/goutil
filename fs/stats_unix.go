// +build !windows

package fs

import (
	"fmt"
	"os"

	sys "golang.org/x/sys/unix"
)

// File system statistics, in bytes and percentage (Usage only).
type FsStats struct {
	Free      uint64
	Available uint64
	Size      uint64
	Used      uint64
	Usage     float64
}

// GetFsStats returns an object holding the file system usage of path.
// If the path is empty, it uses the current working directory.
func GetFsStats(path string) (*FsStats, error) {
	if path == "" {
		if dir, err := os.Getwd(); err != nil {
			return nil, fmt.Errorf("statfs %q: %v", path, err)
		} else {
			path = dir
		}
	}

	var st sys.Statfs_t
	if err := sys.Statfs(path, &st); err != nil {
		return nil, fmt.Errorf("statfs %q: %v", path, err)
	}

	var fs = &FsStats{
		Free:      st.Bfree * uint64(st.Bsize),
		Available: st.Bavail * uint64(st.Bsize),
		Size:      st.Blocks * uint64(st.Bsize),
	}
	fs.Used = fs.Size - fs.Free
	fs.Usage = float64(fs.Used) / float64(fs.Size)

	return fs, nil
}
