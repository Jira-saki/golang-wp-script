package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ฟังก์ชันแกะรอยหาป้ายชื่อและเวอร์ชันปลั๊กอิน (คงเดิม ไว้ใจได้)
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

func main() {
	root := "./public_html"

	fmt.Println("🔍 [Golang Scanner] Starting security sweep in:", root)
	fmt.Println("--------------------------")

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		// 🛡️ ท่าที่ 1: ค่ายกลกรองภัยในโฟลเดอร์ Uploads (โค้ดดั้งเดิมของพี่ Jirasak)
		if strings.Contains(path, "uploads") && strings.HasSuffix(strings.ToLower(path), ".php") {
			if strings.HasSuffix(strings.ToLower(path), "index.php") {
				return nil
			}
			fmt.Printf("⚠️  [SUSPICIOUS] PHP file hiding in uploads directory: %s\n", path)
		}

		// 👑 ท่าที่ 2: สแกนหาปลั๊กอินเวอร์ชันล่าสุด (v2.1 อุดบั๊ก Wordfence)
		if strings.Contains(path, "wp-content/plugins/") && strings.HasSuffix(strings.ToLower(path), ".php") {
			
			// 💡 ทริคตัดสัญญาณรบกวน: แบ่งพาธออกเป็นท่อนๆ เพื่อเช็กความลึก
			// พาธมาตรฐาน: public_html/wp-content/plugins/plugin-dir/main-file.php
			cleanPath := filepath.ToSlash(path)
			segments := strings.Split(cleanPath, "wp-content/plugins/")
			
			if len(segments) > 1 {
				subSegments := strings.Split(segments[1], "/")
				// ถ้าไฟล์ .php นอนอยู่ในโฟลเดอร์ย่อยลึกลงไปอีก (subSegments มากกว่า 2) ให้ข้ามเลย!
				// เช่น "wordfence/models/common/wfConfig.php" -> ยาว 4 ท่อน = ข้าม!!
				if len(subSegments) > 2 {
					return nil
				}
			}

			pName, vNum := checkPluginHeader(path)
			if pName != "" && vNum != "" {
				// ตรวจสอบเพิ่มเติม: ถ้าชื่อปลั๊กอินมีเครื่องหมายแปลกๆ หลุดมา ให้ตีเป็นขยะแล้วข้ามไป
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