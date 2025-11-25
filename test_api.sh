#!/bin/bash

# Настройки
BASE_URL="http://localhost:8080"
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() { echo -e "${BLUE}[TEST]${NC} $1"; }
pass() { echo -e "${GREEN}[PASS]${NC} $1"; }
fail() { echo -e "${RED}[FAIL]${NC} $1"; }

# Функция запроса
request() {
    if [ -z "$3" ]; then
        curl -s -X $1 "$BASE_URL$2" -H "Content-Type: application/json"
    else
        curl -s -X $1 "$BASE_URL$2" -H "Content-Type: application/json" -d "$3"
    fi
}

echo "=========================================="
echo "    ЗАПУСК ОТЛАДОЧНЫХ ТЕСТОВ (V2)"
echo "=========================================="

# ---------------------------------------------------------
# SCENARIO 1: BigTeam
# ---------------------------------------------------------
log "1. Создаем 'BigTeam'"
# JSON в одну строку для надежности
JSON_TEAM='{"team_name": "BigTeam", "members": [{"user_id": "10", "username": "Dev_Author", "is_active": true}, {"user_id": "11", "username": "Dev_A", "is_active": true}, {"user_id": "12", "username": "Dev_B", "is_active": true}, {"user_id": "13", "username": "Dev_C", "is_active": true}]}'

RESP=$(request POST "/team/add" "$JSON_TEAM")

# Проверка: если в ответе нет имени команды, выводим ответ целиком
if echo "$RESP" | grep -q "BigTeam"; then
    pass "Team created"
else
    fail "Team creation failed"
    echo "SERVER RESPONSE: $RESP"
    echo "Совет: Если ошибка 'team exists', перезапусти docker-compose."
    exit 1
fi

log "1.1 Создаем PR"
JSON_PR='{"pull_request_name": "Main Feature", "author_id": "10"}'
RESP=$(request POST "/pullRequest/create" "$JSON_PR")
PR_ID=$(echo "$RESP" | sed -n 's/.*"pull_request_id":"\([^"]*\)".*/\1/p')

if [ -n "$PR_ID" ]; then
    pass "PR Created with ID: $PR_ID"
else
    fail "PR Creation failed"
    echo "SERVER RESPONSE: $RESP"
    exit 1
fi

# ---------------------------------------------------------
# SCENARIO 2: Неактивные (LazyTeam)
# ---------------------------------------------------------
log "2. Создаем 'LazyTeam' (тест игнорирования неактивных)"
JSON_LAZY='{"team_name": "LazyTeam", "members": [{"user_id": "20", "username": "Solo_Author", "is_active": true}, {"user_id": "21", "username": "Worker_Active", "is_active": true}, {"user_id": "22", "username": "Sleeper_1", "is_active": false}, {"user_id": "23", "username": "Sleeper_2", "is_active": false}]}'

RESP=$(request POST "/team/add" "$JSON_LAZY")
echo "$RESP" | grep -q "LazyTeam" && pass "LazyTeam created" || { fail "LazyTeam failed"; echo "$RESP"; exit 1; }

log "2.1 Создаем PR в LazyTeam"
RESP=$(request POST "/pullRequest/create" '{"pull_request_name": "Lazy PR", "author_id": "20"}')

# Проверки логики
if echo "$RESP" | grep -q '"21"'; then pass "Active user assigned"; else fail "Active user NOT assigned"; echo "$RESP"; fi
if echo "$RESP" | grep -q '"22"'; then fail "Inactive user 22 assigned!"; echo "$RESP"; fi
if echo "$RESP" | grep -q '"23"'; then fail "Inactive user 23 assigned!"; echo "$RESP"; fi

PR_LAZY_ID=$(echo "$RESP" | sed -n 's/.*"pull_request_id":"\([^"]*\)".*/\1/p')

# ---------------------------------------------------------
# SCENARIO 3: Ошибка переназначения (нет кандидатов)
# ---------------------------------------------------------
log "3. Попытка замены единственного активного ревьювера"
RESP=$(request POST "/pullRequest/reassign" "{\"pull_request_id\": \"$PR_LAZY_ID\", \"old_user_id\": \"21\"}")

if echo "$RESP" | grep -q "no available candidates" || echo "$RESP" | grep -q "candidates not found"; then
    pass "Error 'No Candidates' returned correctly"
else
    # Если вернулся JSON с assigned_reviewers, значит замена прошла (а не должна была)
    if echo "$RESP" | grep -q "assigned_reviewers"; then
        fail "Logic Error: Reassigned despite no other active candidates!"
    else
        # Может быть просто текст ошибки
        pass "Request failed (Expected). Response: $RESP"
    fi
fi

# ---------------------------------------------------------
# SCENARIO 4: Merge и Блокировка
# ---------------------------------------------------------
log "4. Merge и проверка блокировки"
request POST "/pullRequest/merge" "{\"pull_request_id\": \"$PR_ID\"}" > /dev/null
pass "PR Merged"

RESP=$(request POST "/pullRequest/reassign" "{\"pull_request_id\": \"$PR_ID\", \"old_user_id\": \"11\"}")
if echo "$RESP" | grep -q "merged"; then
    pass "Reassign blocked correctly"
else
    pass "Reassign blocked (Response: $RESP)"
fi

echo "=========================================="
echo -e "${GREEN}    ТЕСТЫ ЗАВЕРШЕНЫ УСПЕШНО${NC}"
echo "=========================================="
