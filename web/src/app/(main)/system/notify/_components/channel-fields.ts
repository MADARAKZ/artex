// English textConfigurationEnglish textTools。
//
// English text**English text**English text：English textDescriptionEnglish text、
// English text，English textConfigurationEnglish text（JSON）English text。
// English text，English text，English text。
// English textTypeEnglish text。English text，English text。
export const KIND_LABEL: Record<string, string> = {
  dingtalk: "English text",
  feishu: "English text",
  wecom: "English text",
  webhook: "English text Webhook",
  telegram: "Telegram",
  email: "English text",
};

// English textConfigurationEnglish text。
//
// English text，English text schema：English text
// Validate（English text/English text），UI English textType，English text。
// English text secret_keys —— English textPasswordEnglish text，
// English text（English text Webhook English text，
// English text secret）。English text itemsEnglish text，
// English text（English text hasFields English textHint）。
export type FieldKind = "text" | "password" | "number" | "select" | "textarea" | "switch" | "kv" | "list";
export interface FieldDef {
  key: string;
  label: string;
  kind: FieldKind;
  placeholder?: string;
  help?: string;
  options?: { value: string; label: string }[];
}
export const CHANNEL_FIELDS: Record<string, FieldDef[]> = {
  dingtalk: [
    {
      key: "webhook",
      label: "Webhook English text",
      kind: "text",
      placeholder: "https://oapi.dingtalk.com/robot/send?access_token=...",
    },
    {
      key: "secret",
      label: "English text",
      kind: "password",
      help: "English text「English text」English text；English text「English textKeyword」English text",
    },
  ],
  feishu: [
    {
      key: "webhook",
      label: "Webhook English text",
      kind: "text",
      placeholder: "https://open.feishu.cn/open-apis/bot/v2/hook/...",
    },
    { key: "secret", label: "English text", kind: "password", help: "English text「English text」English text，English text" },
  ],
  wecom: [
    {
      key: "webhook",
      label: "Webhook English text",
      kind: "text",
      placeholder: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=...",
    },
  ],
  webhook: [
    { key: "url", label: "Goal URL", kind: "text", placeholder: "https://your-endpoint.example.com/hook" },
    {
      key: "method",
      label: "English textMethod",
      kind: "select",
      options: [
        { value: "POST", label: "POST（English text）" },
        { value: "PUT", label: "PUT（English text）" },
        { value: "PATCH", label: "PATCH（English text）" },
        { value: "GET", label: "GET（English text）" },
      ],
    },
    { key: "headers", label: "English text", kind: "kv", help: "English text KEY=VALUE，English text Authorization=Bearer xxx" },
    {
      key: "body_template",
      label: "English text",
      kind: "textarea",
      help:
        "English textDefaultEnglish text。English text：{{.Title}} {{.Batch}} {{.Count}} {{.HomeURL}} {{.SentAt}}，" +
        "English text range .Items English text .Name/.VulnClass/.Severity/.Summary/.Assets/.DetailURL/.StatusLabel。" +
        "English text {{json .Xxx}} English text {{.Xxx}}，English text JSON。",
    },
  ],
  telegram: [
    { key: "bot_token", label: "Bot Token", kind: "password", placeholder: "123456:ABC-DEF..." },
    { key: "chat_id", label: "Chat ID", kind: "text", placeholder: "-1001234567890" },
    {
      key: "base_url",
      label: "API English text",
      kind: "text",
      placeholder: "https://api.telegram.org",
      help: "English text；English text Bot API English text",
    },
  ],
  email: [
    { key: "host", label: "SMTP ServiceEnglish text", kind: "text", placeholder: "smtp.example.com" },
    {
      key: "port",
      label: "Port",
      kind: "number",
      placeholder: "587",
      help: "587 English text STARTTLS；465 English text「English text TLS」English text",
    },
    { key: "username", label: "English text", kind: "text" },
    { key: "password", label: "Password / English text", kind: "password" },
    { key: "from", label: "English text", kind: "text", placeholder: "artex@example.com" },
    { key: "to", label: "English text", kind: "list", help: "English text" },
    { key: "tls", label: "English text TLS", kind: "switch", help: "465 PortEnglish text；587 English textClose（English text STARTTLS）" },
  ],
};

export const SEVERITY_OPTIONS = [
  { value: "", label: "English text" },
  { value: "low", label: "LowEnglish text" },
  { value: "medium", label: "MediumEnglish text" },
  { value: "high", label: "HighEnglish text" },
  { value: "critical", label: "English textCritical" },
];

export type ChannelForm = {
  name: string;
  kind: string;
  mode: "realtime" | "digest";
  enabled: boolean;
  ratePerMin: string;
  config: Record<string, unknown>;
  minSeverity: string;
  includeText: string;
  excludeText: string;
  taskIDsText: string;
  assetIDsText: string;
  onStatusChange: boolean;
};

export const emptyForm = (kind: string): ChannelForm => ({
  name: "",
  kind,
  mode: "realtime",
  enabled: true,
  ratePerMin: "",
  config: {},
  minSeverity: "",
  includeText: "",
  excludeText: "",
  taskIDsText: "",
  assetIDsText: "",
  onStatusChange: false,
});

// parseKV English text「English text KEY=VALUE」English text。
export function parseKV(text: string): Record<string, string> {
  const out: Record<string, string> = {};
  for (const line of text.split("\n")) {
    const t = line.trim();
    if (!t) continue;
    const i = t.indexOf("=");
    if (i > 0) out[t.slice(0, i).trim()] = t.slice(i + 1).trim();
  }
  return out;
}
// parseIDs English text/English text id English text。
export function parseIDs(text: string): number[] {
  return text
    .split(/[\s,，]+/)
    .map((s) => s.trim())
    .filter(Boolean)
    .map((s) => Number(s))
    .filter((n) => Number.isFinite(n) && n > 0);
}
// parseKeywords English text/English textKeywordEnglish text（FindingTypeEnglish text，English text）。
export function parseKeywords(text: string): string[] {
  return text
    .split(/[\n,，]+/)
    .map((s) => s.trim())
    .filter(Boolean);
}
