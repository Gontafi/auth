Тестовое задание

Запуск проекта:
```bash
cp .env.example .env
docker compose up -d --build
```

Routes:<br>
```
  /login POST - получить Access и Refresh токены
  /check GET - проверка Access токена
  /refresh POST - обновить Access и Refresh токены
```

Postman коллекция внутри репозитории
