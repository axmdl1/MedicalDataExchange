# Stage 1 — просто nginx + статические файлы
FROM nginx:1.27-alpine

# Удалим дефолтный конфиг и добавим свой
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем всю статику фронтенда
COPY . /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
