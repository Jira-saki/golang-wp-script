package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ฟังก์ชันแกะรอยหาป้ายชื่อและเวอร์ชันปลั๊กอิน
func checkPluginHeader(filePath string) (string, string) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	pluginName := ""
	version := ""

	for i := 0; i < 30 && scanner.Scan(); i++ {
		line := scanner.Text()
		if strings.Contains(line, "Plugin Name:") {
			parts := strings.Split(line, "Plugin Name:")
			pluginName = strings.TrimSpace(parts[1])
		}
		if strings.Contains(line, "Version:") {
			parts := strings.Split(line, "Version:")
			version = strings.TrimSpace(parts[1])
		}
	}

	return pluginName, version
}

// 🛡️ เช็กเนื้อหาไฟล์ index.php ว่าเป็นไฟล์ว่าง/Silence is golden ของจริงหรือไม่
func isSafeIndexPHP(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	text := strings.TrimSpace(string(content))
	
	// อนุญาตเฉพาะไฟล์ว่าง หรือโค้ดมาตรฐาน WordPress
	if text == "" || text == "<?php" || strings.Contains(text, "Silence is golden.") {
		return true
	}
	return false
}

func main() {
	root := "./public_html"

	fmt.Println("🔍 [Golang Scanner] Starting security sweep in:", root)
	fmt.Println("--------------------------")

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		cleanPath := filepath.ToSlash(path)

		// 🛡️ ท่าที่ 1: ตรวจจับไฟล์ PHP ใน wp-content/uploads/
		if strings.Contains(cleanPath, "/wp-content/uploads/") && strings.HasSuffix(strings.ToLower(cleanPath), ".php") {
			
			// 1. ยกเว้นโฟลเดอร์คอนฟิกไฟร์วอลล์ AIOS (Known Safe)
			if strings.Contains(cleanPath, "uploads/aios/firewall-rules/") {
				return nil
			}

			// 2. ตรวจสอบ index.php: ถ้าเป็น Silence is golden ข้ามได้ แต่ถ้ามีโค้ดแปลกปลอมให้แจ้งเตือนทันที!
			if strings.HasSuffix(strings.ToLower(cleanPath), "index.php") {
				if isSafeIndexPHP(path) {
					return nil
				}
				fmt.Printf("🚨 [HIGH RISK] Malicious index.php found in uploads: %s\n", path)
				return nil
			}

			fmt.Printf("⚠️  [SUSPICIOUS] PHP file hiding in uploads directory: %s\n", path)
		}

		// 👑 ท่าที่ 2: สแกนหาปลั๊กอิน
		if strings.Contains(cleanPath, "wp-content/plugins/") && strings.HasSuffix(strings.ToLower(cleanPath), ".php") {
			segments := strings.Split(cleanPath, "wp-content/plugins/")
			if len(segments) > 1 {
				subSegments := strings.Split(segments[1], "/")
				if len(subSegments) > 2 {
					return nil
				}
			}

			pName, vNum := checkPluginHeader(path)
			if pName != "" && vNum != "" {
				if strings.Contains(pName, "</td>") || strings.Contains(pName, "wp_kses") {
					return nil
				}
				fmt.Printf("📦 [PLUGIN DETECTED] %s (v%s)\n", pName, vNum)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
	} else {
		fmt.Println("--------------------------")
		fmt.Println("✅ [Golang Scanner] Scan completed.")
	}
}