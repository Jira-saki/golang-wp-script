package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var alertBuffer bytes.Buffer
	files, err := os.ReadDir("/proc")
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		pid := file.Name()
		if _, err := fmt.Sscanf(pid, "%d", new(int)); err != nil {
			continue
		}

		// อ่านชื่อ process จาก /proc/[PID]/comm
		commBytes, err := os.ReadFile(filepath.Join("/proc", pid, "comm"))
		if err != nil {
			continue // ข้ามถ้าไม่มีสิทธิ์อ่าน (ของ user อื่นหรือจบโปรเซสไปแล้ว)
		}
		procName := strings.TrimSpace(string(commBytes))

		// Whitelist: ยกเว้นโปรเซสปกติที่ไม่ต้องการเสียเวลาสแกน
		// เช่น ตัว mem_scanner เอง, bash ที่เราล็อกอิน SSH อยู่, sshd, rsync
		switch procName {
		case "mem_scanner", "mem_notify.sh", "bash", "sh", "sshd", "ssh", "rsync", "grep", "ps", "sendmail", "cat":
			continue
		}

		// เปิดอ่าน virtual memory maps (/proc/[PID]/maps)
		mapsPath := filepath.Join("/proc", pid, "maps")
		mapsFile, err := os.Open(mapsPath)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(mapsFile)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}

			perms := fields[1] // e.g. r-xp, rw-p, rwxp
			hasRWX := strings.HasPrefix(perms, "rwx")

			// Anomaly conditions:
			// 1. Memory page is RWX (W^X violation)
			// 2. Mapped file from /tmp or wp-content/uploads
			// 3. Executable mapping whose file was deleted from disk
			isSuspiciousPath := strings.Contains(line, "/tmp") || strings.Contains(line, "wp-content/uploads")
			isDeletedExec := strings.Contains(perms, "x") && strings.Contains(line, "(deleted)")

			if hasRWX || isSuspiciousPath || isDeletedExec {
				alertBuffer.WriteString(fmt.Sprintf("[ALERT] PID: %s (%s)\nMapping: %s\n\n", pid, procName, line))
			}
		}
		mapsFile.Close()
	}

	// if find ANOMALY , then send mail
	if alertBuffer.Len() > 0 {
		triggerNotification(alertBuffer.String())
	}
}

func triggerNotification(report string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/home/xb880282"
	}
	scriptPath := filepath.Join(homeDir, "bin", "mem_notify.sh")

	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Stdin = strings.NewReader(report)
	_ = cmd.Run()
}
