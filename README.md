## Быстрый старт

1. **Скопируйте и настройте переменные окружения:**

   Скопируйте файл `.env_example` в `.env` и отредактируйте значения под свои нужды.

   ```sh
   cp .env_example .env
   ```

2. **Запустите сервисы через Docker Compose:**

   ```sh
   docker compose up --build
   ```

   Это поднимет backend, frontend и базу данных.

3. **Запуск сплойтов:**

   ```sh
   python start_sploit.py your_sploit.* --token <team's api token(из docker compose logs)> -u <farm server url>
   ```

   В логах докера можно найти апи ключ для фермы. Для доступа к веб-оболочке фермы используйте данные из .env
