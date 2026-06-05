// TEMPLATE
package main

import (
	"fmt"
//	"os", "io", "net/http"
)

func doSomething(target string) {

	// [ logic การทำงานหลักอยู่ตรงนี้ ]
	
	// 3. เช็ก Error ทุกครั้งหลังจบคำสั่งสำคัญ (สไตล์ Go ดั้งเดิม)
	// if err != nil { ... }

}

func main() {
	fmt.Println("--- เริ่มทำงาน ---")

	// 4. เรียกใช้งานฟังก์ชันหลักจากใน main
	doSomething("ระบบ A")
	doSomething("ระบบ B")
}