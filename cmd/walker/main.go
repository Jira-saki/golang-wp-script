package main

import (
  "bufio"
  "fmt"
  "io/fs"
  "os"
  "path/filepath"
  "strings"
)

// Whitelist PHP in WP Root
var knownWPRootFiles = map[string]bool{
  "index.php": true, "wp-activate.php": true, "wp-blog-header.php": true,
  "wp-comments-post.php": true, "wp-config.php": true, "wp-config-sample.php": true,
  "wp-cron.php": true, "wp-links-opml.php": true, "wp-load.php": true,
  "wp-login.php": true, "wp-mail.php": true, "wp-settings.php": true,
  "wp-signup.php": true, "wp-trackback.php": true, "xmlrpc.php": true,
  "wordfence-waf.php": true, "wp-cli.phar": true,
	"aios-bootstrap.php": true, // 👈 Whitelist of  All-In-One Security
}

// Asset Directory NO PHP !!!
var assetDirKeywords = []string{
  "/favicons/", "/ogp/", "/images/", "/css/", "/js/",
  "/assets/", "/media/", "/img/", "/fonts/", "/scss/",
}

func isWordPressSite(rootPath string) bool {
  _, err1 := os.Stat(filepath.Join(rootPath, "wp-config.php"))
  _, err2 := os.Stat(filepath.Join(rootPath, "wp-includes"))
  return err1 == nil || err2 == nil
}

func checkPluginHeader(filePath string) (string, string) {
  file, err := os.Open(filePath)
  if err != nil {
    return "", ""
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  pluginName, version := "", ""

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

func isSafeIndexPHP(filePath string) bool {
  content, err := os.ReadFile(filePath)
  if err != nil {
    return false
  }
  text := strings.TrimSpace(string(content))

  if text == "" || text == "<?php" || strings.Contains(text, "Silence is golden.") {
    return true
  }
  return false
}

func main() {
  root := "./public_html"

  fmt.Println("🔍 [Golang Scanner v2.4] Starting security sweep in:", root)
  fmt.Println("--------------------------")

  isWP := isWordPressSite(root)

  err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
    if err != nil || d.IsDir() {
      return nil
    }

    cleanPath := filepath.ToSlash(path)
    lowerPath := strings.ToLower(cleanPath)
    isPHP := strings.HasSuffix(lowerPath, ".php")

    // 🚨 FEATURE 1: Asset Directory Threat Hunting (catch malware on favicons/, ogp/ etc)
    if isPHP {
      hasAssetKeyword := false
      for _, keyword := range assetDirKeywords {
        if strings.Contains(lowerPath, keyword) {
          hasAssetKeyword = true
          break
        }
      }

      if hasAssetKeyword {
        // 🛡️ filtering (False Positives)  WordPress Noise
        isGutenbergAsset := strings.HasSuffix(lowerPath, ".asset.php")
        isWPCoreAsset    := strings.Contains(lowerPath, "/wp-includes/")
        isWPPluginTheme  := strings.Contains(lowerPath, "/wp-content/plugins/") || strings.Contains(lowerPath, "/wp-content/themes/")
        isSafeIndex      := filepath.Base(cleanPath) == "index.php" && isSafeIndexPHP(path)

        // If not WordPress file, do notify
        if !isGutenbergAsset && !isWPCoreAsset && !isWPPluginTheme && !isSafeIndex {
          fmt.Printf("🚨 [CRITICAL] PHP file hiding in asset directory: %s\n", path)
          return nil
        }
      }
    }

    // 🚨 FEATURE 2: Static Site Anomaly Detection (catch PHP on Static Sites)
    if !isWP && isPHP {
      filename := filepath.Base(cleanPath)
      if filename == "index.php" && isSafeIndexPHP(path) {
        return nil
      }
      fmt.Printf("🚨 [STATIC ANOMALY] Unexpected PHP file on static site: %s\n", path)
      return nil
    }

    // 🚨 FEATURE 3: WordPress Uploads & Root Anomaly Detection
    if isWP && isPHP {
      if strings.Contains(lowerPath, "/wp-content/uploads/") {
        if strings.Contains(cleanPath, "uploads/aios/firewall-rules/") {
          return nil
        }
        if strings.HasSuffix(lowerPath, "index.php") {
          if isSafeIndexPHP(path) {
            return nil
          }
          fmt.Printf("🚨 [HIGH RISK] Malicious index.php found in uploads: %s\n", path)
          return nil
        }
        fmt.Printf("⚠️  [SUSPICIOUS] PHP file hiding in uploads directory: %s\n", path)
        return nil
      }

      // Root Level Non-standard Files Check (เช่น tsigcfnd.php, public_html.php)
      relPath, _ := filepath.Rel(root, cleanPath)
      if !strings.Contains(relPath, "/") && !knownWPRootFiles[relPath] {
        fmt.Printf("⚠️  [SUSPICIOUS ROOT PHP] Non-standard PHP file in root: %s\n", path)
      }
    }

    // 📦 FEATURE 4: Plugin Inventory Detection
    if strings.Contains(cleanPath, "wp-content/plugins/") && isPHP {
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