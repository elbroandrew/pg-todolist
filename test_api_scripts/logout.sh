#!/bin/bash


API_URL="http://localhost:8080"
LOGOUT_ENDPOINT="/auth/logout"
COOKIE_FILE="cookies.txt"



echo -e "\nВыходим из системы..."
logout_response=$(curl -s -X POST "$API_URL$LOGOUT_ENDPOINT" \
  -b $COOKIE_FILE)

echo "Результат выхода: $logout_response"

# Очищаем cookies
rm -f $COOKIE_FILE
echo -e "\nФайл cookies.txt удален."

echo -e "\nТестирование завершено!"