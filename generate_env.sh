#!/bin/bash

# Функция для создания .env файла
generate_env_file() {
    local file_name=".env"

    # Список переменных окружения
    cat <<EOL > $file_name
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=
MIGRATION_DIR=

PG_DSN="host= port= dbname= user= password= sslmode="
MIGRATION_DSN="host= port= dbname= user= password= sslmode="

GRPC_HOST=localhost
GRPC_PORT=50051
GRPC_OTHER_PORT=50052

HTTP_HOST=localhost
HTTP_PORT=8080

SWAGGER_HOST=localhost
SWAGGER_PORT=8090

REFRESH_TOKEN_SECRET_KEY=
ACCESS_TOKEN_SECRET_KEY=
REFRESH_TOKEN_EXPIRATION=60
ACCESS_TOKEN_EXPIRATION=5
AUTH_PREFIX=Bearer
EMAIL_TOKEN_SECRET_KEY=
EMAIL_TOKEN_EXPIRATION=60

SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASSWORD=
EOL

    echo "$file_name успешно создан."
}

# Запуск функции генерации
generate_env_file
