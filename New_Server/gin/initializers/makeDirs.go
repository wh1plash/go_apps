package initializers

import (
	"fmt"
	"os"
)

func InitDirs() {
	dirs := []string{"uploads", "unknown"}

	// Make dir if not exist
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Error creating folder: %s", dir)
			return
		}
		fmt.Printf("Make folder: %s", dir)
	}
}
