"use client";

// LLM RetryConfigurationEnglish text：English textRetryEnglish text「English text + English text」。
//
// English text：English text(SDK) → English text(SDK) → English text provider English text → English text → IntentEnglish text。
// English textEndpointEnglish text，English textConfigurationEnglish textDefault；English text，English text。
//
// English textInputEnglish text「English text = English textConfiguration」English text，English text db.RetryRule English text：
//   English text   English text/0 = English textDefault | -1 = CloseEnglish textRetry | >0 = English text
//   English text   English text/0 = English text | >0 = English text

import * as React from "react";

import { Loader2Icon, SaveIcon } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { api } from "@/lib/api";
import type { LLMRetryOverride, LLMRetryPolicy, LLMRetryRule } from "@/lib/types";

export const ZERO_RULE: LLMRetryRule = { attempts: 0, interval_ms: 0 };
export const ZERO_OVERRIDE: LLMRetryOverride = {
  connect: ZERO_RULE,
  empty: ZERO_RULE,
  stream: ZERO_RULE,
};
const ZERO_POLICY: LLMRetryPolicy = {
  ...ZERO_OVERRIDE,
  breaker: ZERO_RULE,
  intent: ZERO_RULE,
};

type LayerMeta = {
  title: string;
  /** English textRetryEnglish text、English text */
  where: string;
  /** English text——English textStatusEnglish text，English text */
  trigger: string;
  /** English text【English text】English text，English text bug */
  skips?: string;
  desc: string;
  attemptsLabel: string;
  /** English textDefaultEnglish text，English text */
  defAttempts: number;
  /** English textDefaultEnglish text，English text */
  defInterval: string;
  /** English text -1 English text */
  offHint: string;
};

export const RETRY_LAYERS = {
  connect: {
    title: "English textRetry",
    where: "SDK · English text 200 English text",
    trigger:
      "English text 200：English text / English text / DNS FailedEnglish text，English text HTTP 408、429、500、502、503、504。",
    skips: "English textStatusEnglish text（400 / 401 / 403 / 404 / 413 / 422 English text）English textDeny，English textFailed，English text。",
    desc: "English text。English text（English text 200），English text。",
    attemptsLabel: "RetryEnglish text",
    defAttempts: 3,
    defInterval: "0.5s→1s→2s English text（English text 8s）",
    offHint: "-1 = English textRetry，FailedEnglish text",
  },
  empty: {
    title: "English textRetry",
    where: "SDK · English text openai English text",
    trigger:
      "HTTP 200、finish_reason English text stop，English text itemsEnglish text——English text、English text、English text。",
    skips: "English text max_tokens English text（English textOutputEnglish text，English text）。",
    desc: "English text prompt，English text，English text。",
    attemptsLabel: "RetryEnglish text",
    defAttempts: 2,
    defInterval: "0.5s→1s→2s English text（English text 8s）",
    offHint: "-1 = English text",
  },
  stream: {
    title: "English text provider English textRetry",
    where: "English text · English textOutputEnglish text",
    trigger:
      "English text（English text 200）English text：English text、English text overloaded、English text 429 / 5xx English text——English text token English text。",
    skips:
      "English text（402 / insufficient_quota，English textConfiguration）、English text（413 / context length，English textDigest）、400 / 401 / 403 / 404 / 422 English textDeny，English textRetry。",
    desc: "English textConfigurationEnglish text。English textOutput，English textOutputEnglish textTool Runs。",
    attemptsLabel: "RetryEnglish text",
    defAttempts: 2,
    defInterval: "0.5s→1s English text（English text 4s）",
    offHint: "-1 = English textIntentEnglish text",
  },
  breaker: {
    title: "English text",
    where: "English text · English text，English text",
    trigger:
      "English textFailed（429、5xx、English text）English text；English text（402）、English text（401 / 403）、English text（404）English textFailedEnglish text，Line English text。",
    skips: "SuccessEnglish text，English textConfigurationEnglish text。",
    desc: "English text，English textConfiguration。StatusEnglish text，English text。",
    attemptsLabel: "English textFailedEnglish text",
    defAttempts: 3,
    defInterval: "1min→5min→30min English text",
    offHint: "-1 = English textFailedEnglish text（English textFailedEnglish text）",
  },
  intent: {
    title: "IntentEnglish text",
    where: "English text · English text，English text",
    trigger:
      "English text：worker English text model_error English text——English textRetryAllEnglish text，English textOutputEnglish text（English text，English text itemsEnglish text）。",
    skips: "English textConfigurationEnglish text，English text；TasksEnglish textPaused / English text / English text，English textTime。",
    desc: "English text itemsIntentEnglish text。English text，English text。",
    attemptsLabel: "English text",
    defAttempts: 2,
    defInterval: "English text 3s",
    offHint: "-1 = English text，English textIntentEnglish text blocked",
  },
} satisfies Record<string, LayerMeta>;

type LayerKey = keyof typeof RETRY_LAYERS;

/** English text，English textInputEnglish text，English text。 */
function humanMs(ms: number) {
  if (!Number.isFinite(ms) || ms <= 0) return "";
  if (ms < 1000) return `${ms}ms`;
  if (ms < 60_000) return `${Number((ms / 1000).toFixed(2))}s`;
  return `${Number((ms / 60_000).toFixed(2))}min`;
}

/** English textInput：English text ↔ 0，English text（"-"、"1e"）English text，English text。 */
function NumField({
  id,
  value,
  onChange,
  placeholder,
  min,
}: {
  id: string;
  value: number;
  onChange: (n: number) => void;
  placeholder: string;
  min: number;
}) {
  const [text, setText] = React.useState(value === 0 ? "" : String(value));
  // English text（English text、English textConfiguration）English text；English text，
  // English text value English text parse English text。
  React.useEffect(() => {
    const incoming = value === 0 ? "" : String(value);
    setText((cur) => (Number(cur || 0) === value ? cur : incoming));
  }, [value]);
  return (
    <Input
      id={id}
      type="number"
      min={min}
      className="w-28 shrink-0"
      value={text}
      placeholder={placeholder}
      onChange={(e) => {
        setText(e.target.value);
        const n = Number(e.target.value);
        onChange(e.target.value.trim() === "" || !Number.isFinite(n) ? 0 : Math.trunc(n));
      }}
    />
  );
}

/** English textRetryEnglish text。idPrefix English text label English text htmlFor。 */
export function RetryRuleFields({
  layer,
  idPrefix,
  value,
  onChange,
  compact,
}: {
  layer: LayerKey;
  idPrefix: string;
  value: LLMRetryRule;
  onChange: (r: LLMRetryRule) => void;
  /** true = ConfigurationEnglish text：English textExpandEnglish text，English text「English text」English text */
  compact?: boolean;
}) {
  const meta = RETRY_LAYERS[layer];
  const human = humanMs(value.interval_ms);
  return (
    <div className={compact ? "grid gap-2" : "grid gap-3 rounded-lg border p-3"}>
      <div className="grid gap-0.5">
        <div className="flex flex-wrap items-baseline gap-2">
          <Label className="text-sm">{meta.title}</Label>
          <span className="text-muted-foreground text-xs">{meta.where}</span>
        </div>
        {/* English text，English textStatusEnglish text——English text，English text。 */}
        <p className="text-muted-foreground text-xs">
          <span className="font-medium text-foreground">English text</span>：{meta.trigger}
        </p>
        {!compact && meta.skips && (
          <p className="text-muted-foreground text-xs">
            <span className="font-medium text-foreground">English text</span>：{meta.skips}
          </p>
        )}
        {!compact && <p className="text-muted-foreground text-xs">{meta.desc}</p>}
      </div>
      <div className="flex flex-wrap items-center gap-x-4 gap-y-2">
        <div className="flex items-center gap-2">
          <Label htmlFor={`${idPrefix}-${layer}-n`} className="text-muted-foreground text-xs">
            {meta.attemptsLabel}
          </Label>
          <NumField
            id={`${idPrefix}-${layer}-n`}
            min={-1}
            value={value.attempts}
            placeholder={`Default ${meta.defAttempts}`}
            onChange={(n) => onChange({ ...value, attempts: n })}
          />
        </div>
        <div className="flex items-center gap-2">
          <Label htmlFor={`${idPrefix}-${layer}-ms`} className="text-muted-foreground text-xs">
            English text ms
          </Label>
          <NumField
            id={`${idPrefix}-${layer}-ms`}
            min={0}
            value={value.interval_ms}
            placeholder="DefaultEnglish text"
            onChange={(n) => onChange({ ...value, interval_ms: n })}
          />
          <span className="text-muted-foreground text-xs">{human ? `English text ${human}` : meta.defInterval}</span>
        </div>
      </div>
      {!compact && <p className="text-muted-foreground text-xs">English text = English textDefault；{meta.offHint}。</p>}
    </div>
  );
}

/** English textConfigurationEnglish text（English textEndpointEnglish text）。 */
export function ProfileRetryFields({
  value,
  onChange,
}: {
  value: LLMRetryOverride;
  onChange: (o: LLMRetryOverride) => void;
}) {
  return (
    <div className="grid gap-3 rounded-lg border p-3">
      <div className="grid gap-0.5">
        <Label className="text-sm">RetryEnglish text</Label>
        <p className="text-muted-foreground text-xs">
          English textConfigurationEnglish text，English text「RetryEnglish text」English textDefault。English text = English text；English text -1 = English textRetry；
          English text。English textIntentEnglish text，English text。
        </p>
      </div>
      {(["connect", "empty", "stream"] as const).map((k) => (
        <div key={k} className="border-t pt-3 first:border-t-0 first:pt-0">
          <RetryRuleFields
            compact
            layer={k}
            idPrefix="pf"
            value={value[k]}
            onChange={(r) => onChange({ ...value, [k]: r })}
          />
        </div>
      ))}
    </div>
  );
}

/** 「RetryEnglish text」tab：English textDefaultEnglish text。 */
export function RetryPolicyPanel() {
  const [policy, setPolicy] = React.useState<LLMRetryPolicy>(ZERO_POLICY);
  const [loading, setLoading] = React.useState(true);
  const [saving, setSaving] = React.useState(false);

  const load = React.useCallback(async () => {
    setLoading(true);
    try {
      const p = await api.llmRetryPolicy();
      setPolicy({ ...ZERO_POLICY, ...p });
    } catch (e) {
      toast.error(`English textRetryEnglish textFailed：${(e as Error).message}`);
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    void load();
  }, [load]);

  async function save() {
    if (saving) return;
    setSaving(true);
    try {
      // English text，English textRefresh，English text。
      const saved = await api.saveLLMRetryPolicy(policy);
      setPolicy({ ...ZERO_POLICY, ...saved });
      toast.success("English textSave，English text（English textParameters）");
    } catch (e) {
      toast.error(`Save failed：${(e as Error).message}`);
    } finally {
      setSaving(false);
    }
  }

  const set = (k: LayerKey) => (r: LLMRetryRule) => setPolicy((p) => ({ ...p, [k]: r }));

  if (loading) {
    return (
      <div className="flex items-center gap-2 rounded-lg border border-dashed p-10 text-muted-foreground text-sm">
        <Loader2Icon className="size-4 animate-spin" /> English textRetryEnglish text…
      </div>
    );
  }

  return (
    <div className="grid gap-4">
      <div className="rounded-lg border bg-muted/30 p-3 text-muted-foreground text-xs leading-relaxed">
        English textFailedEnglish textRetry，English text：
        <span className="text-foreground"> English text → English text → English text provider English text → English text → IntentEnglish text</span>
        。English text，English text
        <span className="text-foreground">English text</span>
        English text——English text，English text。
        AllEnglish textDefaultEnglish text，English text。English textConfigurationEnglish text。
      </div>

      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
        {(Object.keys(RETRY_LAYERS) as LayerKey[]).map((k) => (
          <RetryRuleFields key={k} layer={k} idPrefix="gl" value={policy[k]} onChange={set(k)} />
        ))}
      </div>

      <div className="flex gap-2">
        <Button onClick={save} disabled={saving}>
          {saving ? <Loader2Icon className="animate-spin" /> : <SaveIcon />}
          Save
        </Button>
        <Button variant="outline" onClick={() => setPolicy(ZERO_POLICY)} disabled={saving}>
          AllEnglish textDefault
        </Button>
      </div>
    </div>
  );
}
