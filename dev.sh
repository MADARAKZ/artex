#!/usr/bin/env bash
# English text：English text(:8787) + English text(:8788) English text English text next dev(:5173) English text。
# English text /api English text；Ctrl-C English text。
#
# English text（English text）English text README「English text」English text，English text。
set -euo pipefail
cd "$(dirname "$0")"

# English text（English text + English text）。
cleanup() { kill 0 2>/dev/null || true; }
trap cleanup EXIT INT TERM

# English text（English text go run，English text）；English text work agent English text「English text」English text。
go run ./cmd/artex -addr :8787 -proxy 127.0.0.1:8788 &

# English text（Vite/Next dev server，/api English text :8787）。
( cd web && npm run dev ) &

echo "[dev] English text :8787 / English text :8788 / English text http://localhost:5173  (Ctrl-C English text)"
wait
