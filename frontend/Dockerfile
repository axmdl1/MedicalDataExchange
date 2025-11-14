FROM nginx:alpine

WORKDIR /usr/share/nginx/html

# удалим дефолтный индекс nginx
RUN rm -rf ./*

# копируем весь фронт внутрь образа
COPY . .

# nginx уже запускается сам как ENTRYPOINT/CMD