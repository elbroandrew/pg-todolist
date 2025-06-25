#!/bin/bash

API_URL="http://localhost:8080"
LOGIN_ENDPOINT="/auth/login"
TASKS_ENDPOINT="/tasks"
LOGOUT_ENDPOINT="/auth/logout"
COOKIE_FILE="cookies.txt"
EMAIL="test1@test.com"
PASSWORD="test1"

# Функция для извлечения данных из JSON ответа
extract_json_value() {
    echo "$1" | grep -o "\"$2\":\"[^\"]*\"" | cut -d':' -f2 | tr -d '"'
}

# 1. Логин и сохранение cookies
echo "Логинимся пользователем $EMAIL..."
login_response=$(curl -s -X POST "$API_URL$LOGIN_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" \
  -c $COOKIE_FILE)

access_token=$(extract_json_value "$login_response" "access_token")
if [ -z "$access_token" ]; then
    echo "Ошибка при логине: $login_response"
    exit 1
fi

echo "Успешный логин! Access Token получен."


# 2. Запрос к защищенному ресурсу
echo -e "\nЗапрашиваем список задач..."
tasks_response=$(curl -s "$API_URL$TASKS_ENDPOINT" \
  -H "Authorization: Bearer $access_token" \
  -b $COOKIE_FILE)

echo "Ответ от сервера:"
echo "$tasks_response"