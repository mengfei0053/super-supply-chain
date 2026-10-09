# HTTP MCP server for Super Supply Chain.
# Compose publishes 172.17.0.1:3100 only. NPM terminates TLS in front of it.
FROM registry.cn-hangzhou.aliyuncs.com/mengfei0053/node:22

WORKDIR /usr/src/ssc-mcp

COPY package.json package-lock.json ./
# Include devDependencies so `tsc` is present even if the base image sets NODE_ENV=production.
RUN npm ci --include=dev

COPY tsconfig.json ./
COPY src ./src
RUN npm run build \
 && npm prune --omit=dev

ENV NODE_ENV=production \
    SSC_MCP_TRANSPORT=http \
    SSC_MCP_HOST=0.0.0.0 \
    SSC_MCP_PORT=3100 \
    SSC_BASE_URL=https://ssc.mengfei.tech

EXPOSE 3100

CMD ["node", "dist/index.js"]
