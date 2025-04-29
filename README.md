### Run guide:
1. Перейти в директорию dev-config
```
cd ./dev-config
```
2. Запустить [**Docker**](https://www.docker.com/) у себя на устройстве
3. Собрать образы и поднять контейнеры для клиента, сервера и базы данных (На Windows рекомендуется использовать [wsl](https://learn.microsoft.com/ru-ru/windows/wsl/install))
```
wsl
```
```
docker compose up --build
```

Клиент будет доступен в браузере по адресу `http://localhost:5173`
