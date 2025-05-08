# Stage 1: build
# Node image is used on build stage because vite requires some node.js API missing in bun
FROM node:24-alpine3.20 AS build

WORKDIR /app

COPY . .

RUN npm install && cd vue && npm run build-only

# Stage 2: image
FROM oven/bun:1.0.5-alpine

WORKDIR /app

COPY --from=build /app/bff /app

EXPOSE 3000

ENTRYPOINT ["bun", "run", "src/index.ts"]