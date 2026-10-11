package agent

import (
	"log"
	"path/filepath"

	"github.com/Autumn-27/norma/agentcore"
	"github.com/Autumn-27/norma/noaadapter"
)

// noaWarn returns a diagnostics sink tagging non-fatal noa messages with the
// session, routed through the package logger (agents have no per-instance one).
func noaWarn(session string) func(string) {
	return func(msg string) { log.Printf("[noa] %s: %s", session, msg) }
}

// noa English text norma v0.4.0 English text「English text」English text,English text
// English text。English text compaction English text:noaadapter.Enable English text,English text
// English text(Compactor)、Compress English text,English text Enable English text
// (English text compaction English text)。English text agent English text noaEnabledFn English text,English text run English text
// English text,English text run,English text agent。

// enableNoa English text noa English text opts。archiveRoot English text
// (English text workDir,English text agent English text <workDir>/noa English text,English text/English text),sessionID
// English text(English text,English text)。
//
// noa English text:English text。English text onWarn English text。
// English text opts.Compaction,English text agentcore English text「English text」English text。
func enableNoa(opts *agentcore.Options, enabled func() bool, archiveRoot, sessionID string, onWarn func(string)) {
	if enabled == nil || !enabled() {
		return
	}
	if opts.OnWarn == nil {
		opts.OnWarn = onWarn
	}
	if err := noaadapter.Enable(opts, noaadapter.Options{
		ArchiveBaseDir: filepath.Join(archiveRoot, "noa"),
		SessionID:      sessionID,
		OnWarn:         onWarn,
	}); err != nil {
		if onWarn != nil {
			onWarn("noa English text,English text:" + err.Error())
		}
		return
	}
	// Compactor English text Compaction,English text agentcore English text;English text。
	opts.Compaction = nil
}
