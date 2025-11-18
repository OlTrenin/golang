## Примеры использования

```bash
# 1. Добавляем число 3
curl -X POST http://localhost:8080/numbers -H "Content-Type: application/json" -d '{"number": 3}'
# Response: {"numbers":[3]}

# 2. Добавляем число 2
curl -X POST http://localhost:8080/numbers -H "Content-Type: application/json" -d '{"number": 2}'
# Response: {"numbers":[2,3]}

# 3. Добавляем число 1
curl -X POST http://localhost:8080/numbers -H "Content-Type: application/json" -d '{"number": 1}'
# Response: {"numbers":[1,2,3]}

# 4. Добавляем число 5
curl -X POST http://localhost:8080/numbers -H "Content-Type: application/json" -d '{"number": 5}'
# Response: {"numbers":[1,2,3,5]}
```