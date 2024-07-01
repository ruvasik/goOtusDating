Запуск
======
```bash
curl -O https://raw.githubusercontent.com/OtusTeam/highload/master/homework/people.v2.csv
docker-compose up -d
```

Генерация пользователей
-----------------------
```bash
docker-compose run --rm app ./generate
```

Отчеты
======
1.

Postman
=======
[postman/go-otus-dating-api.postman_collection.json](https://github.com/ruvasik/goOtusDating/blob/master/postman/go-otus-dating-api.postman_collection.json)
