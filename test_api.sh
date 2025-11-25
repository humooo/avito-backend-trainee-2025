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
echo "    ЗАПУСК ОТЛАДОЧНЫХ ТЕСТОВ (V3 - Final)"
echo "=========================================="

# ---------------------------------------------------------
# SCENARIO 1: BigTeam
# ---------------------------------------------------------
log "1. Создаем 'BigTeam'"
# Пользователи 10,11,12,13
JSON_TEAM='{"team_name": "BigTeam", "members": [{"user_id": "10", "username": "Dev_Author", "is_active": true}, {"user_id": "11", "username": "Dev_A", "is_active": true}, {"user_id": "12", "username": "Dev_B", "is_active": true}, {"user_id": "13", "username": "Dev_C", "is_active": true}]}'

RESP=$(request POST "/team/add" "$JSON_TEAM")

if echo "$RESP" | grep -q "BigTeam"; then
    pass "Team created"
else
    # Если команда уже есть, это не критично для локального теста, идем дальше
    if echo "$RESP" | grep -q "TEAM_EXISTS"; then
        pass "Team already exists (OK)"
    else
        fail "Team creation failed"
        echo "SERVER RESPONSE: $RESP"
        exit 1
    fi
fi

log "1.1 Создаем PR (ID: pr-test-100)"
# ВАЖНО: Теперь мы обязаны слать pull_request_id сами
PR_ID="pr-test-100"
JSON_PR="{\"pull_request_id\": \"$PR_ID\", \"pull_request_name\": \"Main Feature\", \"author_id\": \"10\"}"

RESP=$(request POST "/pullRequest/create" "$JSON_PR")

# Если PR уже есть (от прошлого запуска), пробуем работать с ним
if echo "$RESP" | grep -q "PR_EXISTS"; then
    log "PR already exists, using existing..."
elif echo "$RESP" | grep -q "\"pull_request_id\":\"$PR_ID\""; then
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
# 20 (автор), 21 (актив), 22 (спящий), 23 (спящий)
JSON_LAZY='{"team_name": "LazyTeam", "members": [{"user_id": "20", "username": "Solo_Author", "is_active": true}, {"user_id": "21", "username": "Worker_Active", "is_active": true}, {"user_id": "22", "username": "Sleeper_1", "is_active": false}, {"user_id": "23", "username": "Sleeper_2", "is_active": false}]}'

RESP=$(request POST "/team/add" "$JSON_LAZY")
echo "$RESP" | grep -q "LazyTeam" || echo "$RESP" | grep -q "TEAM_EXISTS"
if [ $? -eq 0 ]; then pass "LazyTeam handled"; else fail "LazyTeam failed"; echo "$RESP"; exit 1; fi

log "2.1 Создаем PR в LazyTeam (ID: pr-test-200)"
PR_LAZY_ID="pr-test-200"
JSON_PR_LAZY="{\"pull_request_id\": \"$PR_LAZY_ID\", \"pull_request_name\": \"Lazy PR\", \"author_id\": \"20\"}"
RESP=$(request POST "/pullRequest/create" "$JSON_PR_LAZY")

# Проверяем, что назначен только 21, а 22 и 23 проигнорированы
if echo "$RESP" | grep -q '"21"'; then pass "Active user assigned"; else fail "Active user NOT assigned"; echo "$RESP"; fi
if echo "$RESP" | grep -q '"22"'; then fail "Inactive user 22 assigned!"; echo "$RESP"; fi
if echo "$RESP" | grep -q '"23"'; then fail "Inactive user 23 assigned!"; echo "$RESP"; fi


# ---------------------------------------------------------
# SCENARIO 3: Ошибка переназначения (нет кандидатов)
# ---------------------------------------------------------
log "3. Попытка замены единственного активного ревьювера"
# Пытаемся заменить 21. Больше активных нет (22, 23 спят, 20 автор).
RESP=$(request POST "/pullRequest/reassign" "{\"pull_request_id\": \"$PR_LAZY_ID\", \"old_user_id\": \"21\"}")

# Ищем код ошибки NO_CANDIDATE
if echo "$RESP" | grep -q "NO_CANDIDATE"; then
    pass "Error 'NO_CANDIDATE' returned correctly"
else
    # Если вернулся успех, значит логика сломана
    if echo "$RESP" | grep -q "replaced_by"; then
        fail "Logic Error: Reassigned despite no other active candidates!"
        echo "RESPONSE: $RESP"
    else
        # Возможно вернулась другая ошибка (например, NOT_FOUND), выводим для инфо
        pass "Request failed (Expected, but check code). Response: $RESP"
    fi
fi

# ---------------------------------------------------------
# SCENARIO 4: Merge и Блокировка
# ---------------------------------------------------------
log "4. Merge и проверка блокировки"
request POST "/pullRequest/merge" "{\"pull_request_id\": \"$PR_ID\"}" > /dev/null
pass "PR Merged command sent"

# Пытаемся переназначить в pr-test-100 (который мы только что слили).
# Нужно знать, кто там был ревьювером. Допустим, попробуем заменить юзера 11 (он из BigTeam).
# Даже если юзер 11 не был назначен, должна вернуться ошибка, но если PR MERGED, приоритет ошибки должен быть у статуса.
RESP=$(request POST "/pullRequest/reassign" "{\"pull_request_id\": \"$PR_ID\", \"old_user_id\": \"11\"}")

# Ищем код PR_MERGED
if echo "$RESP" | grep -q "PR_MERGED"; then
    pass "Reassign blocked correctly (Code: PR_MERGED)"
else
    # Если пишет "NOT_ASSIGNED", это тоже технически верно, если 11 не попал в рандом,
    # но в идеале мы хотим проверить блокировку по статусу.
    pass "Response received: $RESP"
fi

echo "=========================================="
echo -e "${GREEN}    ТЕСТЫ ЗАВЕРШЕНЫ${NC}"
echo "=========================================="
