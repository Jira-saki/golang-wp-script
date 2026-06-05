package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func main() {
	// 1. specific folder to scan
	root := "./public_html"

	fmt.Println("🔍 [Golang Scanner] Starting security sweep in:", root)
	fmt.Println("--------------------------")

	// 2. WalkDir function; to scan all files and folders recursively
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // unaccess-folder then skip it
		}
		// skip if the folder itself
		if d.IsDir() {
			return nil
		}

		// 3. Catch; if found "uploads" in path and .php file
		if strings.Contains(path, "uploads") && strings.HasSuffix(strings.ToLower(path), ".php") {

			// 💡 ค่ายกลกรองภัย: ถ้าไฟล์ลงท้ายด้วย "index.php" ให้ข้ามไปเลย ไม่ต้องแจ้งเตือน
      if strings.HasSuffix(strings.ToLower(path), "index.php") {
        return nil
      }

			fmt.Printf("⚠️  [SUSPICIOUS] PHP file hiding in uploads directory: %s\n", path)
		}

		return nil
	})

	// check if there is error
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
	} else {
		fmt.Println("--------------------------")
		fmt.Println("✅ [Golang Scanner] Scan completed.")
	}
}
