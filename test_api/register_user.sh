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

# Регистрация пользователя
echo "1. Регистрируем пользователя $EMAIL"
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

if echo "$REGISTER_RESPONSE" | grep -q "error"; then
    if echo "$REGISTER_RESPONSE" | grep -q "already exists"; then
        echo "Пользователь уже существует, пропускаем регистрацию"
    else
        error "Ошибка регистрации: $REGISTER_RESPONSE"
    fi
else
    echo "Успешная регистрация: $REGISTER_RESPONSE"
fi