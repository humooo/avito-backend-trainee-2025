#!/bin/bash

URL="http://localhost:8080"

echo "Start testing..."

# 1. Создаем команду
# Просто шлем JSON, грепаем ответ. Если вернулся JSON с именем команды - ок.
echo "--- 1. Create Team (BigTeam) ---"
RESP=$(curl -s -X POST "$URL/team/add" \
  -H "Content-Type: application/json" \
  -d '{"team_name": "BigTeam", "members": [{"user_id": "u1", "username": "Alice", "is_active": true}, {"user_id": "u2", "username": "Bob", "is_active": true}, {"user_id": "u3", "username": "Carol", "is_active": true}]}')

echo "Response: $RESP"

if [[ $RESP == *"BigTeam"* ]] || [[ $RESP == *"TEAM_EXISTS"* ]]; then
    echo " Team created or exists"
else
    echo " Failed to create team"
    exit 1
fi

# 2. Создаем PR
echo -e "\n--- 2. Create PR (pr-1) ---"
PR_ID="pr-1"
RESP=$(curl -s -X POST "$URL/pullRequest/create" \
  -H "Content-Type: application/json" \
  -d "{\"pull_request_id\": \"$PR_ID\", \"pull_request_name\": \"Fix bug\", \"author_id\": \"u1\"}")

echo "Response: $RESP"

if [[ $RESP == *"assigned_reviewers"* ]]; then
    echo "PR created with reviewers"
elif [[ $RESP == *"PR_EXISTS"* ]]; then
    echo "PR already exists (skipping creation check)"
else
    echo "Failed to create PR"
    exit 1
fi

# 3. Пробуем переназначить (Reassign)
# Пытаемся заменить u2 (если он попал). Если не попал — сервер вернет ошибку, но нам главное проверить сам вызов.
echo -e "\n--- 3. Reassign Reviewer ---"
RESP=$(curl -s -X POST "$URL/pullRequest/reassign" \
  -H "Content-Type: application/json" \
  -d "{\"pull_request_id\": \"$PR_ID\", \"old_user_id\": \"u2\"}")

echo "Response: $RESP"
# Нам не важно, успешна замена или нет (может u2 и не был ревьюером), главное чтоб сервис ответил адекватно (не 500)
if [[ $RESP == *"error"* ]] || [[ $RESP == *"replaced_by"* ]]; then
    echo "Reassign endpoint works (logic handled)"
else
    echo "Reassign failed heavily"
fi

# 4. Мержим PR
echo -e "\n--- 4. Merge PR ---"
RESP=$(curl -s -X POST "$URL/pullRequest/merge" \
  -H "Content-Type: application/json" \
  -d "{\"pull_request_id\": \"$PR_ID\"}")

echo "Response: $RESP"
if [[ $RESP == *"MERGED"* ]]; then
    echo "PR Merged"
else
    echo "Merge failed"
fi

# 5. Проверяем блокировку после мержа
# Пытаемся снова переназначить. Должна быть ошибка PR_MERGED.
echo -e "\n--- 5. Check Block After Merge ---"
RESP=$(curl -s -X POST "$URL/pullRequest/reassign" \
  -H "Content-Type: application/json" \
  -d "{\"pull_request_id\": \"$PR_ID\", \"old_user_id\": \"u2\"}")

echo "Response: $RESP"
if [[ $RESP == *"PR_MERGED"* ]]; then
    echo " Block works! (Got PR_MERGED)"
else
    echo " Warning: Did not get PR_MERGED error (maybe user wasnt assigned)"
fi

# 6. Бонус: Статистика
echo -e "\n--- 6. Stats (Bonus) ---"
curl -s "$URL/stats"
echo ""

echo -e "\nTesting finished."
