#!/bin/sh
# ARTEX English text（Linux / macOS / Docker ENTRYPOINT）
#
# English text：
#   ./start.sh                       English text（Ctrl-C English text）
#   nohup ./start.sh >artex.log 2>&1 &   English text
#   ./start.sh -addr :9000           English text artex
#
# English text：English text artex English text，English text。
#
#   0      English text        → English text
#   75     English text        → English text（English text"English text"English text"English text"）
#   English text   English text                → English text（1→2→4…English text 60 English text）
#
# English text、SHA256 English text：English text sh English text bat English text，
# English text——English text，English text
# English text，English text。English text/English text Go English text（selfupdate English text），
# English text artex English text，English text。
set -u

cd "$(dirname "$0")" || exit 1

BIN=./artex
[ -x "$BIN" ] || { echo "[artex] English text $BIN" >&2; exit 1; }

RESTART_CODE=75
MAX_DELAY=60

child=0
stopping=0

# English text artex English text。
#
# Docker English text：docker stop English text SIGTERM English text PID 1（English text），
# English text。English text artex English text、English text，10 English text SIGKILL
# English text，English text。
forward() {
	stopping=1
	if [ "$child" -ne 0 ]; then
		kill -TERM "$child" 2>/dev/null || true
	fi
}
trap forward INT TERM

delay=1
while :; do
	"$BIN" "$@" &
	child=$!

	# English text wait English text >128。English text，
	# English text wait English text。
	wait "$child"
	code=$?
	if [ "$code" -gt 128 ]; then
		wait "$child"
		code=$?
	fi
	child=0

	if [ "$stopping" -eq 1 ]; then
		echo "[artex] English text"
		exit 0
	fi

	case "$code" in
		0)
			echo "[artex] English text"
			exit 0
			;;
		"$RESTART_CODE")
			# English text/English text：English text artex English text（English text selfupdate.Bootstrap）。
			echo "[artex] English text（English text）…"
			delay=1
			;;
		*)
			echo "[artex] English text (code=$code)，${delay}s English text" >&2
			sleep "$delay"
			delay=$((delay * 2))
			[ "$delay" -gt "$MAX_DELAY" ] && delay=$MAX_DELAY
			;;
	esac
done
