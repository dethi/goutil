package main

import (
	"fmt"
	"os"

	"github.com/dethi/goutil/fs"
)

func main() {
	const GB = 1024 * 1024 * 1024

	if st, err := fs.GetStats(""); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Printf("Free: %d GB\nTotal: %d GB\nUsage: %3.2f%%\n",
			st.Free/GB, st.Size/GB, st.Usage*100)
	}
}
