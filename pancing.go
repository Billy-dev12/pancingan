package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

const targetURL = "https://7hs7b37h9ys.pages.dev"

func GetRequest() {
	// Logika pengetesan ketahanan API dengan cara loop GET request selama 4 kali dan 1 detik cooldown
	client := &http.Client{Timeout: 10 * time.Second}
	fmt.Printf("\n=== Memulai Pengetesan API (%s) ===\n", targetURL)

	// loop di hentikan ketika aku mengklik ctlr +c
	for i := 1; i <= 1000; i++ {
		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			fmt.Printf("[Test %d/10] Error membuat request: %v\n", i, err)
			continue
		}

		// Menyamar sebagai browser
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[Test %d/4] Gagal terhubung ke server: %v\n", i, err)
		} else {
			fmt.Printf("[Test %d/4] Respons Server: %d %s\n", i, resp.StatusCode, http.StatusText(resp.StatusCode))
			resp.Body.Close()
		}

		// Delay 1 detik sebelum request berikutnya (kecuali jika ini adalah request terakhir)
		if i < 4 {
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("=== Pengetesan API Selesai ===")
}

const TargetURL = "https://palma.qzz.io/shpe/9ys/send.php"

type Payload struct {
	Phone string `json:"phone"`
}

// generateRandomPhone menghasilkan nomor telepon acak sepanjang 10 digit di belakang +628
func generateRandomPhone() string {
	var sb strings.Builder
	sb.WriteString("+628")
	for i := 0; i < 10; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		sb.WriteString(n.String())
	}
	return sb.String()
}

func sendRequest(client *http.Client) {
	phone := generateRandomPhone()
	payload := Payload{Phone: phone}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("[ERROR] Failed to marshal JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", TargetURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[ERROR] Failed to create request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Menyamar sebagai browser PC biasa (Google Chrome)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "id-ID,id;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] Server gagal merespon: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Jika server menolak koneksi (misal 403 Forbidden atau 429 Too Many Requests)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[ERROR] Akses Ditolak (%d %s) untuk nomor %s\n", resp.StatusCode, http.StatusText(resp.StatusCode), phone)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read response: %v\n", err)
		return
	}

	respStr := string(body)

	// Mengecek apakah respons mengandung "ok":true
	if strings.Contains(respStr, `"ok":true`) {
		fmt.Printf("[SUCCESS] Terkirim: %s\n", phone)
	} else {
		// Menghapus newline di response agar log lebih rapi
		respStr = strings.TrimSpace(strings.ReplaceAll(respStr, "\n", " "))
		// Memotong respons jika terlalu panjang agar tidak memenuhi layar
		if len(respStr) > 80 {
			respStr = respStr[:80] + "..."
		}
		fmt.Printf("[ERROR] Respons tidak sesuai (Status: %d): %s\n", resp.StatusCode, respStr)
	}
}

func main() {
	mode := flag.String("mode", "post", "Pilih metode pengetesan: 'get' atau 'post'")
	flag.Parse()

	if *mode == "get" {
		GetRequest()
	} else if *mode == "post" {
		fmt.Println("------------------------------------------")
		fmt.Printf("Starting Security Audit on: %s\n", TargetURL)
		fmt.Println("Press [CTRL+C] to stop.")
		fmt.Println("------------------------------------------")

		// Menggunakan satu client HTTP yang sama (connection pooling) agar performa lebih cepat
		tr := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		}
		client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

		for {
			var wg sync.WaitGroup

			// Jalankan 5 request paralel di dalam goroutine
			for i := 0; i < 50; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					sendRequest(client)
				}()
			}

			wg.Wait() // Tunggu 5 request batch ini selesai
			fmt.Println("--- Batch selesai, istirahat 1 detik ---")
			time.Sleep(1 * time.Second)
		}
	} else {
		fmt.Println("Mode tidak dikenali. Gunakan '-mode get' atau '-mode post'")
	}
}
