package main

import (
	"fmt"
	"net"
	"time"
)

func checkPort(host string, port int) {
	// 1. รวม host และ port เข้าด้วยกันในรูปแบบ "host:port" เช่น "google.com:443"
	address := fmt.Sprintf("%s:%d", host, port)

	// 2.set timeout 2.0 sec and TCP handshake
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)

	// 3. ถ้ามี error แสดงว่า port นั้นไม่เปิด (close)
	if err != nil {
		fmt.Printf("🔴 [DEAD] เครื่อง %s พอร์ต %d ปิดอยู่ หรือติดต่อไม่ได้! (Error: %v)\n", host, port, err)
		return
	}

	// 4. close connection after talk (defer)
	defer conn.Close()

	fmt.Printf("🟢 [ALIVE] เครื่อง %s พอร์ต %d เปิดอยู่!\n", host, port)
}

func main() {
	fmt.Println("--- เริ่มทำการตรวจสอบระบบเครือข่าย ---")

	checkPort("127.0.0.1", 80)      // เช็กเว็บเซิร์ฟเวอร์ในเครื่องตัวเอง (Loopback)
	checkPort("google.com", 443)    // เช็กพอร์ต HTTPS ของ Google บนโลกอินเทอร์เน็ต
	checkPort("192.168.1.50", 9200) // สมมติลองเช็กเครื่อง OpenSearch ในวงแล็บภายใน
}
