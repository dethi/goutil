package main

import (
	"fmt"
	"os"

	"github.com/dethi/goutil/fs"
)

func main() {
	const GB = 1024 * 1024 * 1024

	if fsinfo, err := fs.Stat(""); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("Free: %d GB\nTotal: %d GB\nUsage: %3.2f%%\n",
			fsinfo.Free/GB, fsinfo.Size/GB, fsinfo.Usage*100)
	}
}
