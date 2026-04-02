
#!/bin/bash

TARGET="https://teknologisantuy.vercel.app/api/gemini"

echo "Memulai pengetesan API Gemini (3 kali request)..."

for i in {1..7}
do
  echo "Mengirim request ke-$i..."
  
  curl -s -X POST "$TARGET" \
    -H "Content-Type: application/json" \
    -d '{
      "message": "Test request ke-'$i'. anonim cuma ngetes limitasi lo ya, jangan marah."
    }' | jq '.' # Pake jq biar output JSON-nya rapi (kalau udah install)

  echo -e "\nCooldown 2 detik..."
  sleep 1
done

echo "Pengetesan selesai. Sekarang waktunya TIDUR, Bill!"
