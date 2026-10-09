# Build ssc-mcp for HTTP/stdio. Prefer Aliyun mirror (same as app Dockerfile).
FROM registry.cn-hangzhou.aliyuncs.com/mengfei0053/node:22 AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY tsconfig.json ./
COPY src ./src
RUN npm run build

FROM registry.cn-hangzhou.aliyuncs.com/mengfei0053/node:22

WORKDIR /app

ENV NODE_ENV=production

COPY package.json package-lock.json ./
RUN npm ci --omit=dev

COPY --from=builder /app/dist ./dist

EXPOSE 3100

CMD ["node", "dist/index.js"]
