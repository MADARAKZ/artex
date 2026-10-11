// Centralised status → color/label semantics, reused across the whole app.
// Spec §8.3: Intent / English text / Tasks / CriticalEnglish text each have a consistent color set.

export type Tone = "neutral" | "blue" | "green" | "amber" | "red" | "rose" | "violet" | "slate";

export const toneClasses: Record<Tone, string> = {
  neutral: "bg-muted text-muted-foreground border-transparent",
  blue: "bg-blue-500/15 text-blue-600 dark:text-blue-400 border-blue-500/20",
  green: "bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/20",
  amber: "bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/20",
  red: "bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/20",
  // rose English text「Critical」——English text,English text「High」English text。
  rose: "bg-rose-600 text-white border-rose-600 dark:bg-rose-600 dark:text-white",
  violet: "bg-violet-500/15 text-violet-600 dark:text-violet-400 border-violet-500/20",
  slate: "bg-slate-500/15 text-slate-600 dark:text-slate-400 border-slate-500/20",
};

export const toneDot: Record<Tone, string> = {
  neutral: "bg-muted-foreground",
  blue: "bg-blue-500",
  green: "bg-emerald-500",
  amber: "bg-amber-500",
  red: "bg-red-500",
  rose: "bg-white",
  violet: "bg-violet-500",
  slate: "bg-slate-500",
};

interface StatusMeta {
  label: string;
  tone: Tone;
}

const intent: Record<string, StatusMeta> = {
  open: { label: "English text", tone: "slate" },
  running: { label: "English text", tone: "blue" },
  paused: { label: "English textPaused", tone: "amber" },
  done: { label: "Completed", tone: "green" },
  // blocked = English text/API/English textRetryEnglish text，English text itemsIntentEnglish text（English textGoalEnglish text）。
  blocked: { label: "English text", tone: "red" },
  // exhausted = English text/TimeEnglish text、English text（English text）。
  exhausted: { label: "English text", tone: "violet" },
  // stopped = English textDeleteStatus（English text,English text）。
  stopped: { label: "Stopped", tone: "slate" },
  // deleted = English textDeleteEnglish textIntent（English text，DeleteEnglish text delete_reason English text）。
  deleted: { label: "English textDelete", tone: "slate" },
};

const task: Record<string, StatusMeta> = {
  created: { label: "English text", tone: "slate" },
  queued: { label: "English text", tone: "amber" },
  running: { label: "Running", tone: "blue" },
  paused: { label: "English textPaused", tone: "amber" },
  done: { label: "Completed", tone: "green" },
  failed: { label: "Failed", tone: "red" },
  timeout: { label: "Timed out", tone: "amber" },
};

const severity: Record<string, StatusMeta> = {
  critical: { label: "Critical", tone: "rose" },
  high: { label: "High", tone: "red" },
  medium: { label: "Medium", tone: "amber" },
  low: { label: "Low", tone: "slate" },
};

const finding: Record<string, StatusMeta> = {
  pending: { label: "English text", tone: "amber" },
  in_progress: { label: "English text", tone: "blue" },
  confirmed: { label: "English text", tone: "red" },
  resolved: { label: "English text", tone: "green" },
  fixed: { label: "English text", tone: "green" },
  false_positive: { label: "English text", tone: "slate" },
  ignored: { label: "English text", tone: "neutral" },
  duplicate: { label: "English text", tone: "neutral" },
  risk_accepted: { label: "English text", tone: "violet" },
};

const engine: Record<string, StatusMeta> = {
  exploring: { label: "ExploringEnglish text", tone: "blue" },
  paused: { label: "English textPaused", tone: "amber" },
  stalled: { label: "English text", tone: "red" },
  idle: { label: "English text", tone: "neutral" },
};

const goal: Record<string, StatusMeta> = {
  open: { label: "English text", tone: "blue" },
  met: { label: "English text", tone: "green" },
  abandoned: { label: "English text", tone: "slate" },
};

const audit: Record<string, StatusMeta> = {
  allow: { label: "English text", tone: "green" },
  block: { label: "English text", tone: "red" },
};

const node: Record<string, StatusMeta> = {
  observed: { label: "English text", tone: "slate" },
  confirmed: { label: "English text", tone: "green" },
  tombstoned: { label: "English text", tone: "neutral" },
};

// English textStatus。sending English text blue English text amber：English text「English text」，
// English text「English text、English text」，English text pending English text。
const delivery: Record<string, StatusMeta> = {
  pending: { label: "English textSend", tone: "amber" },
  sending: { label: "SendEnglish text", tone: "blue" },
  sent: { label: "English text", tone: "green" },
  failed: { label: "Failed", tone: "red" },
  skipped: { label: "English text", tone: "neutral" },
};

const maps = {
  intent,
  task,
  severity,
  finding,
  engine,
  goal,
  audit,
  node,
  delivery,
} as const;

export type StatusDomain = keyof typeof maps;

export function statusMeta(domain: StatusDomain, key: string): StatusMeta {
  return maps[domain][key] ?? { label: key, tone: "neutral" };
}
