#!/bin/bash

API_URL="http://localhost:8080"
EMAIL="test1@test.com"
PASSWORD="test1"
TASK_TITLE="Тестовая задача $(date '+%Y-%m-%d %H:%M:%S')"

# Функция для вывода ошибок
error() {
    echo -e "\033[31mОШИБКА: $1\033[0m"
    exit 1
}


# Авторизация
echo -e "\n2. Авторизуемся как $EMAIL"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    error "Не удалось получить токен: $LOGIN_RESPONSE"
fi

echo "${TOKEN}"