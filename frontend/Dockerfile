FROM oven/bun:1.2.15-alpine AS base

WORKDIR /app

ENV NEXT_TELEMETRY_DISABLED=1

FROM base AS deps

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

FROM base AS builder

COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN bun --bun next build

FROM oven/bun:1.2.15-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV HOSTNAME=0.0.0.0
ENV PORT=8000

COPY --from=builder /app ./

EXPOSE 8000

CMD ["bun", "--bun", "next", "start", "-p", "8000"]
