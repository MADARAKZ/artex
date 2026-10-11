#!/usr/bin/env bash
# ARTEX English text：① Docker English text（English text）  ② English text（English text）
# English text install.sh English text：install English text，update English text。
# DB English text——artex English text schema.sql（English text ADD COLUMN/CREATE
# INDEX IF NOT EXISTS），English text“English text”。English text（pgdata English text、./data、./skills）English text。
set -euo pipefail
cd "$(cd "$(dirname "$0")" && pwd)"

info(){ printf '\033[36m[*]\033[0m %s\n' "$*"; }
ok(){   printf '\033[32m[+]\033[0m %s\n' "$*"; }
warn(){ printf '\033[33m[!]\033[0m %s\n' "$*"; }
die(){  printf '\033[31m[x]\033[0m %s\n' "$*" >&2; exit 1; }
ask(){  local p="$1" d="${2:-}" a; read -rp "$p${d:+ [$d]}: " a; echo "${a:-$d}"; }

# ── English text：English text（compose/English text/English text）───────
sync_repo(){
  [ -d .git ] && command -v git >/dev/null 2>&1 || { warn "English text git English text，English text git pull"; return; }
  [ "$(ask 'English text (git pull --ff-only)? (y/n)' y)" = y ] || return
  if ! git pull --ff-only; then
    warn "git pull English text（English text）——English text，English text"
  fi
}

# ── ① Docker English text ───────────────────────────────
update_docker(){
  command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1 \
    || die "English text docker / docker compose，English text ./install.sh English text"
  [ -f .env ] || die "English text .env，English text ./install.sh English text"

  # English text：English text tag（English text .env English text ARTEX_TAG，English text latest）
  local tag; tag="$(ask 'English text tag（English text .env / latest）' '')"
  if [ -n "$tag" ]; then
    if grep -q '^ARTEX_TAG=' .env; then
      sed -i.bak "s|^ARTEX_TAG=.*|ARTEX_TAG=${tag}|" .env && rm -f .env.bak
    else
      printf '\nARTEX_TAG=%s\n' "$tag" >> .env
    fi
    ok "English text ARTEX_TAG English text ${tag}"
  fi

  # English text artex：postgres English text 16-alpine，English text（English text，
  # English text）。artex English text depends_on postgres，English text
  # up English text pg English text，English text、English text。
  info "English text（English text artex）…"
  docker compose pull artex
  info "English text（artex English text schema）…"
  docker compose up -d artex
  ok "English text → http://localhost:8787"
  info "English text：docker compose logs -f artex"
  info "English text（English text）：docker image prune -f"
}

# ── ② English text ──────────────────────────────
update_local(){
  command -v go >/dev/null 2>&1 || die "English text Go（>=1.26）：https://go.dev/dl/"
  [ -f config.json ] || warn "English text config.json——English text ./install.sh"
  ok "Go: $(go version)"

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
  warn "English text artex English text（English text schema）"
}

echo "=============================="
echo "  ARTEX English text"
echo "  1) Docker English text（English text）"
echo "  2) English text（go English text）"
echo "=============================="
case "$(ask 'English text' 1)" in
  1) sync_repo; update_docker ;;
  2) sync_repo; update_local ;;
  *) die "English text" ;;
esac
