#!/bin/bash

API_URL="http://localhost:8080"

# Получаем токен из первого скрипта
TOKEN=$(./test_api/auth_user.sh)
if [ $? -ne 0 ]; then
  exit 1
fi

# Параметры задачи
TASK_TITLE="Новая задача $(date '+%Y-%m-%d %H:%M')"

# Создаем задачу
RESPONSE=$(curl -s -X POST "$API_URL/tasks" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"$TASK_TITLE"}')

# Проверяем результат
if echo "$RESPONSE" | grep -q '"ID"'; then
  ID=$(echo "$RESPONSE" | grep -o '"ID":[0-9]*' | cut -d':' -f2)
  echo "Успешно! ID задачи: $ID"
else
  echo "Ошибка создания задачи:" >&2
  echo "$RESPONSE" >&2
  exit 1
fi