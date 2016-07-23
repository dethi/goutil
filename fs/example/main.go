package main

import (
	"fmt"
	"os"

	"github.com/dethi/goutil/fs"
)

const GB = 1024 * 1024 * 1024

func main() {
	if st, err := fs.GetFsStats(""); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("Free: %d GB\nTotal: %d GB\nUsage: %3.2f%%\n",
			st.Free/GB, st.Size/GB, st.Usage*100)
	}
}
