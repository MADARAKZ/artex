#!/usr/bin/env bash
# ARTEX English text：① English text Docker  ② English text
set -euo pipefail
cd "$(cd "$(dirname "$0")" && pwd)"

info(){ printf '\033[36m[*]\033[0m %s\n' "$*"; }
ok(){   printf '\033[32m[+]\033[0m %s\n' "$*"; }
warn(){ printf '\033[33m[!]\033[0m %s\n' "$*"; }
die(){  printf '\033[31m[x]\033[0m %s\n' "$*" >&2; exit 1; }
ask(){  local p="$1" d="${2:-}" a; read -rp "$p${d:+ [$d]}: " a; echo "${a:-$d}"; }
rand(){ head -c 18 /dev/urandom | base64 | tr -dc 'A-Za-z0-9' | head -c 24; }

# ── docker English text / English text ───────────────────
ensure_docker(){
  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    ok "English text docker English text docker compose"; return
  fi
  warn "English text docker / docker compose"
  case "$(uname -s)" in
    Linux)
      if [ "$(ask 'English text Docker? (y/n)' y)" = y ]; then
        curl -fsSL https://get.docker.com | sh
        sudo usermod -aG docker "$USER" || true
        ok "Docker English text（English text sudo）"
      else
        die "English text docker English text"
      fi ;;
    Darwin) die "macOS English text Docker Desktop：https://www.docker.com/products/docker-desktop/" ;;
    *)      die "English text docker English text" ;;
  esac
}

# ── ① English text Docker ───────────────────────────────
install_docker(){
  ensure_docker
  if [ ! -f .env ]; then
    cp .env.example .env 2>/dev/null || true
    local pw key
    pw="$(ask 'Postgres English text（English text）' "$(rand)")"
    key="$(ask 'ANTHROPIC_API_KEY（English text，English text UI English text）' '')"
    sed -i.bak "s|^POSTGRES_PASSWORD=.*|POSTGRES_PASSWORD=${pw}|" .env
    sed -i.bak "s|^ANTHROPIC_API_KEY=.*|ANTHROPIC_API_KEY=${key}|" .env
    rm -f .env.bak
    ok "English text .env（POSTGRES_PASSWORD English text）"
  else
    info "English text .env"
  fi
  info "English text…"
  docker compose pull || true
  docker compose up -d
  ok "English text → http://localhost:8787"
  info "English text：docker compose logs -f artex"
}

# ── ② English text ──────────────────────────────
install_local(){
  echo "English text："
  echo "  1) English text PostgreSQL"
  echo "  2) English text Docker English text PostgreSQL（English text docker）"
  case "$(ask 'English text' 1)" in
    2)
      ensure_docker
      local pw; pw="$(ask 'Postgres English text（English text）' "$(rand)")"
      docker run -d --name artex-pg -p 5432:5432 \
        -e POSTGRES_USER=artex -e POSTGRES_PASSWORD="$pw" -e POSTGRES_DB=artex \
        -v artex-pg:/var/lib/postgresql/data postgres:16-alpine
      DB_HOST=127.0.0.1 DB_PORT=5432 DB_USER=artex DB_PASS="$pw" DB_NAME=artex DB_SSL=disable ;;
    *)
      DB_HOST="$(ask 'English text' 127.0.0.1)"
      DB_PORT="$(ask 'English text' 5432)"
      DB_USER="$(ask 'English text' artex)"
      DB_PASS="$(ask 'English text' '')"
      DB_NAME="$(ask 'English text' artex)"
      DB_SSL="$(ask 'sslmode (disable/require)' disable)" ;;
  esac

  # English text config.json
  cat > config.json <<JSON
{
  "database": {
    "host": "${DB_HOST}",
    "port": ${DB_PORT},
    "user": "${DB_USER}",
    "password": "${DB_PASS}",
    "dbname": "${DB_NAME}",
    "sslmode": "${DB_SSL}"
  }
}
JSON
  ok "English text config.json"

  # go English text
  command -v go >/dev/null 2>&1 || die "English text Go，English text Go（>=1.26）：https://go.dev/dl/"
  ok "Go: $(go version)"

  # English text node English text
  if command -v npm >/dev/null 2>&1; then
    info "English text…"
    ( cd web && npm ci && npm run build:static )
    rm -rf server/webui/dist && cp -r web/out server/webui/dist
    info "English text…"
    CGO_ENABLED=0 go build -tags embedui -trimpath -o artex ./cmd/artex
  else
    warn "English text npm：English text**English text**English text（English text npm run dev）"
    CGO_ENABLED=0 go build -o artex ./cmd/artex
  fi
  ok "English text → ./artex"

  info "English text…（Ctrl-C English text）"
  ./artex
}

echo "=============================="
echo "  ARTEX English text"
echo "  1) English text Docker English text"
echo "  2) English text（go English text）"
echo "=============================="
case "$(ask 'English text' 1)" in
  1) install_docker ;;
  2) install_local ;;
  *) die "English text" ;;
esac
