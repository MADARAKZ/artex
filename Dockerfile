# syntax=docker/dockerfile:1
#
# English text（English text）：English text，English text**English text Linux English text**。
# English text CI English text binaries job English text（English text Go、English text QEMU），English text
# English text dist/<TARGETARCH>/artex。English text arm64 English text apt English text，
# English text Next/Go English text，English text。
#
# English text，English text：
#   cd web && npm run build:static && cd ..
#   cp -r web/out server/webui/dist
#   CGO_ENABLED=0 GOARCH=amd64 go build -tags embedui -o dist/amd64/artex ./cmd/artex
#   docker build -t artex:local .
FROM python:3.12-slim-bookworm
ARG TARGETARCH
# English text：ripgrep / curl / vim，English text recon English text（English text）。
# Node English text NodeSource English text 20.x：bookworm English text apt nodejs English text 18，Playwright English text >=20。
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates ripgrep curl wget vim git jq unzip \
      dnsutils iputils-ping netcat-openbsd inetutils-telnet whois nmap \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && rm -rf /var/lib/apt/lists/*
# English text Playwright MCP English text CLI（English text），English text npx English text。
# @playwright/mcp：browser MCP English text `npx @playwright/mcp`（English text，English text -y/@latest）。
# @playwright/cli：English text playwright-cli，English text --help English text。
# English text playwright（English text），English text --with-deps English text chromium English text，
# English text MCP/CLI English text，English text。
RUN npm install -g @playwright/mcp@latest @playwright/cli@latest playwright@latest \
    && playwright-cli --help \
    && playwright install --with-deps chromium \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
# English text（dist/amd64/artex English text dist/arm64/artex）
COPY dist/${TARGETARCH}/artex /app/artex
# English text：English text，English text。
# English text SIGTERM English text artex —— docker stop English text PID 1，
# English text artex English text、English text，10 English text SIGKILL English text。
COPY start.sh /app/start.sh
RUN chmod +x /app/artex /app/start.sh
COPY skills/ /app/skills/
# data/（SQLite + jwt.key）English text
VOLUME ["/app/data"]
EXPOSE 8787 8788
ENTRYPOINT ["/app/start.sh"]
CMD ["-addr", ":8787", "-proxy", ":8788"]
