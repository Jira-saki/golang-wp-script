import socket

def check_port(host, port):

  s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    
    # 2. ตั้งเวลา Timeout 2.0 วินาที (ถ้าติดต่อไม่ได้ใน 2 วินาทีให้ตัดบททันที ไม่ต้องนั่งรอจนสคริปต์ค้าง)
  s.settimeout(2.0)
    
  try:
        # 3. ลองยื่นมือไปจับมือ (TCP Handshake) กับปลายทาง
        s.connect((host, port))
        print(f"🟢 [ALIVE] เครื่อง {host} พอร์ต {port} เปิดทำงานปกติ!")
        
        # 4. ปิดท่อหลังจากคุยเสร็จเพื่อคืนทรัพยากรให้ระบบ
        s.close()
  except Exception as e:
        print(f"🔴 [DEAD] เครื่อง {host} พอร์ต {port} ปิดอยู่ หรือติดต่อไม่ได้! (Error: {e})")

# --- โหมดทดสอบรันขำๆ ---
print("--- เริ่มทำการตรวจสอบระบบเครือข่าย ---")
check_port("127.0.0.1", 80)        # เช็กเว็บเซิร์ฟเวอร์ในเครื่องตัวเอง (Loopback)
check_port("google.com", 443)     # เช็กพอร์ต HTTPS ของ Google บนโลกอินเทอร์เน็ต
check_port("192.168.1.50", 9200)   # สมมติลองเช็กเครื่อง OpenSearch ในวงแล็บภายใน