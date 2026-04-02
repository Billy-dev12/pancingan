package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

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

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] Server gagal merespon: %v\n", err)
		return
	}
	defer resp.Body.Close()

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
		fmt.Printf("[ERROR] Server mulai engap/error: %s\n", respStr)
	}
}

func main() {
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
		for i := 0; i < 5; i++ {
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
}
