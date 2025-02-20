FROM node:lts-alpine

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

EXPOSE 5173

CMD ["npm", "start"]