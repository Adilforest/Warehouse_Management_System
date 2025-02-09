# Используем базовый образ для фронтенда (например, nginx)
FROM nginx:alpine

# Копируем файлы фронтенда в контейнер
COPY . /usr/share/nginx/html

# Открываем порт 80
EXPOSE 80

# Команда для запуска nginx
CMD ["nginx", "-g", "daemon off;"]