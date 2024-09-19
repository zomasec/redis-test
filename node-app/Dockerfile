FROM node:14

WORKDIR /app
COPY index.js .
RUN npm install redis

CMD ["node", "index.js"]
