FROM node:18-alpine as builder

WORKDIR /app

COPY ./package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:18-alpine as production

COPY --from=builder /app/build .
COPY --from=builder /app/package.json .
COPY --from=builder /app/package-lock.json .

RUN npm ci --omit dev

USER node:node

CMD ["node", "."]
