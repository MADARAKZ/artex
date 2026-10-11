"use client";

import * as React from "react";

import { CpuIcon, FlaskConicalIcon, KeyboardIcon, RadioTowerIcon, SearchIcon, ShieldAlertIcon } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { api } from "@/lib/api";
import { CHAT_SEND_MODE_OPTIONS, type ChatSendMode, setChatSendMode, useChatSendMode } from "@/lib/chat-send-mode";
import type { Settings } from "@/lib/types";

import { UpdateCard } from "./_components/update-card";

export default function SystemSettingsPage() {
  const [trafficCapture, setTrafficCapture] = React.useState(false);
  const [agentTrafficBinding, setAgentTrafficBinding] = React.useState(false);
  const [webSearch, setWebSearch] = React.useState(false);
  const [backend, setBackend] = React.useState("ddgs");
  const [braveKeySet, setBraveKeySet] = React.useState(false);
  const [braveKeyInput, setBraveKeyInput] = React.useState("");
  const [tavilyKeySet, setTavilyKeySet] = React.useState(false);
  const [tavilyKeyInput, setTavilyKeyInput] = React.useState("");
  const [savingTavilyKey, setSavingTavilyKey] = React.useState(false);
  const [proxyInput, setProxyInput] = React.useState("");
  const [savingProxy, setSavingProxy] = React.useState(false);
  const [globalProxyInput, setGlobalProxyInput] = React.useState("");
  const [savingGlobalProxy, setSavingGlobalProxy] = React.useState(false);
  const [testing, setTesting] = React.useState(false);
  const [loaded, setLoaded] = React.useState(false);
  const [saving, setSaving] = React.useState(false);
  const [savingKey, setSavingKey] = React.useState(false);
  const [pyInterp, setPyInterp] = React.useState("");
  const [workers, setWorkers] = React.useState("3");
  const [savingWorkers, setSavingWorkers] = React.useState(false);
  // ActionsEnglish text(DefaultEnglish text)。
  const [injectPlanner, setInjectPlanner] = React.useState(true);
  const [injectWorker, setInjectWorker] = React.useState(true);
  // English textFunctions:noa English textDigest(DefaultEnglish text)。
  const [noaCompaction, setNoaCompaction] = React.useState(false);
  // English text：English text /api/settings，English text localStorage。
  const sendMode = useChatSendMode();

  const apply = React.useCallback((s: Settings) => {
    setTrafficCapture(!!s.traffic_capture);
    setAgentTrafficBinding(!!s.agent_traffic_binding);
    setWebSearch(!!s.web_search_enabled);
    setBackend(s.web_search_backend || "ddgs");
    setBraveKeySet(!!s.brave_key_set);
    setTavilyKeySet(!!s.tavily_key_set);
    setProxyInput(s.web_search_proxy ?? "");
    setGlobalProxyInput(s.global_proxy ?? "");
    setPyInterp(s.python_interpreter ?? "");
    setWorkers(String(s.workers ?? 3));
    setInjectPlanner(s.constraints_inject_planner !== false);
    setInjectWorker(s.constraints_inject_worker !== false);
    setNoaCompaction(!!s.noa_compaction);
  }, []);

  const saveWorkers = () => {
    const n = Number(workers);
    if (!Number.isInteger(n) || n <= 0) {
      toast.error("English text 0 English text");
      return;
    }
    setSavingWorkers(true);
    api
      .setSettings({ workers: n })
      .then((s) => {
        apply(s);
        toast.success("English textSaveEnglish text agent English text（English textTasksEnglish text）");
      })
      .catch((e) => toast.error("Save failed：" + (e as Error).message))
      .finally(() => setSavingWorkers(false));
  };

  const savePython = () => {
    setSaving(true);
    api
      .setSettings({ python_interpreter: pyInterp.trim() })
      .then((s) => {
        apply(s);
        toast.success("English textSave Python English textConfiguration");
      })
      .catch((e) => toast.error("Save failed：" + (e as Error).message))
      .finally(() => setSaving(false));
  };
  const detectPython = () => {
    setSaving(true);
    api
      .detectPython()
      .then((r) => setPyInterp(r.python_interpreter))
      .catch(() => undefined)
      .finally(() => setSaving(false));
  };

  React.useEffect(() => {
    api
      .settings()
      .then(apply)
      .catch(() => undefined)
      .finally(() => setLoaded(true));
  }, [apply]);

  const toggleTraffic = (v: boolean) => {
    setTrafficCapture(v); // optimistic
    setSaving(true);
    api
      .setSettings({ traffic_capture: v })
      .then(apply)
      .catch(() => setTrafficCapture(!v)) // revert on failure
      .finally(() => setSaving(false));
  };

  const toggleInjectPlanner = (v: boolean) => {
    setInjectPlanner(v); // optimistic
    api
      .setSettings({ constraints_inject_planner: v })
      .then(apply)
      .catch(() => setInjectPlanner(!v)); // revert on failure
  };

  const toggleAgentTrafficBinding = (v: boolean) => {
    setAgentTrafficBinding(v);
    setSaving(true);
    api
      .setSettings({ agent_traffic_binding: v })
      .then((s) => {
        apply(s);
        toast.success(v ? "English text Agent English textTraffic" : "English textClose Agent English textTraffic");
      })
      .catch((e) => {
        setAgentTrafficBinding(!v);
        toast.error(`Save failed：${(e as Error).message}`);
      })
      .finally(() => setSaving(false));
  };

  const toggleInjectWorker = (v: boolean) => {
    setInjectWorker(v); // optimistic
    api
      .setSettings({ constraints_inject_worker: v })
      .then(apply)
      .catch(() => setInjectWorker(!v)); // revert on failure
  };

  const toggleNoaCompaction = (v: boolean) => {
    setNoaCompaction(v); // optimistic
    api
      .setSettings({ noa_compaction: v })
      .then((s) => {
        apply(s);
        toast.success(v ? "English text noa English textDigest（English text）" : "English textClose noa English textDigest（English textDigest）");
      })
      .catch((e) => {
        setNoaCompaction(!v); // revert on failure
        toast.error(`Save failed：${(e as Error).message}`);
      });
  };

  // Persist a web-search patch (enable and/or backend). Optimistic with refetch.
  const saveWebSearch = (patch: Partial<Settings>) => {
    setSaving(true);
    api
      .setSettings(patch)
      .then((s) => {
        apply(s);
        toast.success("English textSaveEnglish textSearchConfiguration");
      })
      .catch((e) => {
        toast.error("Save failed：" + (e as Error).message);
        api
          .settings()
          .then(apply)
          .catch(() => undefined);
      })
      .finally(() => setSaving(false));
  };

  const saveBraveKey = () => {
    setSavingKey(true);
    api
      .setSettings({ brave_search_api_key: braveKeyInput })
      .then((s) => {
        apply(s);
        setBraveKeyInput("");
        toast.success("English textSave Brave API Key");
      })
      .catch((e) => toast.error("Save failed：" + (e as Error).message))
      .finally(() => setSavingKey(false));
  };

  const saveTavilyKey = () => {
    setSavingTavilyKey(true);
    api
      .setSettings({ tavily_search_api_key: tavilyKeyInput })
      .then((s) => {
        apply(s);
        setTavilyKeyInput("");
        toast.success("English textSave Tavily API Key");
      })
      .catch((e) => toast.error("Save failed：" + (e as Error).message))
      .finally(() => setSavingTavilyKey(false));
  };

  const saveProxy = () => {
    setSavingProxy(true);
    api
      .setSettings({ web_search_proxy: proxyInput.trim() })
      .then((s) => {
        apply(s);
        toast.success(proxyInput.trim() ? "English textSaveEnglish text" : "English text（English text）");
      })
      .catch((e) => toast.error("Save failed：" + (e as Error).message))
      .finally(() => setSavingProxy(false));
  };

  const saveGlobalProxy = () => {
    setSavingGlobalProxy(true);
    api
      .setSettings({ global_proxy: globalProxyInput.trim() })
      .then((s) => {
        apply(s);
        toast.success(globalProxyInput.trim() ? "English textSaveEnglish text" : "English text（English text）");
      })
      .catch((e) => toast.error("Save failed：" + (e as Error).message))
      .finally(() => setSavingGlobalProxy(false));
  };

  // Run a real "test" search ("test") against the CURRENT form values (backend +
  // proxy + entered key), falling back to saved values server-side. Toasts result.
  const runTest = () => {
    setTesting(true);
    api
      .testWebSearch({
        web_search_backend: backend,
        web_search_proxy: proxyInput.trim(),
        brave_search_api_key: braveKeyInput,
        tavily_search_api_key: tavilyKeyInput,
      })
      .then((r) => {
        if (r.ok) toast.success(`SearchEnglish textSuccess · ${r.backend} English text ${r.count}  itemsEnglish text`);
        else toast.error("SearchEnglish textFailed：" + (r.error || "English text"));
      })
      .catch((e) => toast.error("SearchEnglish textFailed：" + (e as Error).message))
      .finally(() => setTesting(false));
  };

  // brave-free selected but no key stored and none being entered → tool stays off.
  const braveNeedsKey = webSearch && backend === "brave-free" && !braveKeySet;

  return (
    <div className="flex flex-1 flex-col gap-4 md:gap-6">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">System Settings</h1>
        <p className="text-muted-foreground text-sm">English text</p>
      </div>

      {/* English text grid：English textSearchEnglish text，English text（brave/tavily
          English text key InputEnglish text itemsEnglish text）。grid English text、English text，
          English text。English text mb English text gap——English text
          column-gap English text，English text。 */}
      <div className="columns-1 gap-4 md:gap-6 lg:columns-2">
        <UpdateCard />

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <RadioTowerIcon className="size-4" />
              TrafficEnglish text
            </CardTitle>
            <CardDescription>
              English text，English text Agent English text HTTP TrafficEnglish text，English text Agent English text traffic_search / traffic_get
              ToolsEnglish textConfiguration（HintEnglish text）。
              <br />
              Close（Default）English textTraffic：Agent
              <b>English text</b>English textConfigurationEnglish textTrafficTools，HintEnglish text<b>English text</b>English text。English text Agent English text。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between gap-4">
            <Label htmlFor="traffic-capture" className="text-sm font-normal text-muted-foreground">
              {trafficCapture ? "English text · English textTrafficEnglish text" : "English textClose · English text、English text"}
            </Label>
            <Switch
              id="traffic-capture"
              checked={trafficCapture}
              disabled={!loaded || saving}
              onCheckedChange={toggleTraffic}
            />
          </CardContent>
        </Card>

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <RadioTowerIcon className="size-4" />
              Agent English textTraffic
            </CardTitle>
            <CardDescription id="agent-traffic-binding-description">
              DefaultClose。English text，FindingEnglish text Agent English text HTTP English text/English text，English textTrafficEnglish text。
              <b>English textToolsEnglish text Token English text。</b>
              <br />
              TCP、English textTrafficEnglish text。English textTrafficEnglish text、English textSaveEnglish text。 English text Agent
              English text；CloseEnglish textDenyEnglish text。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between gap-4">
            <Label htmlFor="agent-traffic-binding" className="text-sm font-normal text-muted-foreground">
              {agentTrafficBinding ? "English text · English text Token English text" : "English textClose · English text"}
            </Label>
            <Switch
              id="agent-traffic-binding"
              aria-describedby="agent-traffic-binding-description"
              checked={agentTrafficBinding}
              disabled={!loaded || saving}
              onCheckedChange={toggleAgentTrafficBinding}
            />
          </CardContent>
        </Card>

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <RadioTowerIcon className="size-4" />
              English text
            </CardTitle>
            <CardDescription>
              English text Agent English text<b>GoalTraffic</b>English text（English text IP / English text）。English text <b>http / https / socks5</b>，English text{" "}
              <code>user:pass</code> English text。English text=English text。
              <br />
              English text<b>TrafficEnglish text</b>English text，English text<b>English text</b>（TrafficEnglish text，English text）；CloseEnglish text，English text
              Agent English text bash / WebFetch English text。English textSearchEnglish text、LLM English text。
              <br />
              <b>Hint</b>：socks5 English text<b>CloseEnglish text</b>English textToolsEnglish text <code>ALL_PROXY</code> English text（curl
              English text，English textToolsEnglish text）； English text socks5，English textTrafficEnglish text——English text MITM
              English text，ToolsEnglish text、English text。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-2">
            <Label htmlFor="global-proxy" className="text-sm font-normal text-muted-foreground">
              English text
            </Label>
            <div className="flex items-center gap-2">
              <Input
                id="global-proxy"
                autoComplete="off"
                placeholder="socks5://user:pass@host:1080 English text http://host:port（English text=English text）"
                value={globalProxyInput}
                disabled={!loaded || savingGlobalProxy}
                onChange={(e) => setGlobalProxyInput(e.target.value)}
              />
              <Button type="button" onClick={saveGlobalProxy} disabled={!loaded || savingGlobalProxy}>
                Save
              </Button>
            </div>
            <p className="text-muted-foreground text-xs">
              {globalProxyInput.trim() ? "English textConfiguration · English textGoalTrafficEnglish text" : "English textConfiguration · GoalTrafficEnglish text"}
            </p>
          </CardContent>
        </Card>

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <ShieldAlertIcon className="size-4" />
              ActionsEnglish text
            </CardTitle>
            <CardDescription>
              English text，English textTasksEnglish text<b>ActionsEnglish text</b>（English textTasksOverview「ActionsEnglish text」English text allow/deny  itemsEnglish text）English text Agent
              English textSystemHint，English textExploringEnglish text（English text「English textPort」「English text」）。
              <br />
              English text <b>English text（planner）</b>English text <b>English text（worker）</b>
              ；DefaultEnglish text。English text（English text），English text Agent。CloseEnglish text Agent English text。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-4">
            <div className="flex items-center justify-between gap-4">
              <Label htmlFor="inject-planner" className="text-sm font-normal text-muted-foreground">
                English text（planner）{injectPlanner ? " · English text" : " · English textClose"}
              </Label>
              <Switch
                id="inject-planner"
                checked={injectPlanner}
                disabled={!loaded}
                onCheckedChange={toggleInjectPlanner}
              />
            </div>
            <div className="flex items-center justify-between gap-4">
              <Label htmlFor="inject-worker" className="text-sm font-normal text-muted-foreground">
                English text（worker）{injectWorker ? " · English text" : " · English textClose"}
              </Label>
              <Switch
                id="inject-worker"
                checked={injectWorker}
                disabled={!loaded}
                onCheckedChange={toggleInjectWorker}
              />
            </div>
          </CardContent>
        </Card>

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <FlaskConicalIcon className="size-4" />
              English textFunctions
            </CardTitle>
            <CardDescription>
              English text，DefaultClose。English text Agent English text，English text。
              <br />
              <b>noa English textDigest</b>：English textDigestEnglish textChatEnglish text（norma v0.4.0）。English text Agent（
              <b>English text / English text / English text Agent / Chat</b>）English text noa English text，English textDigest，
              DigestEnglish textTasksEnglish text。English text（English text），English text Agent；
              CloseEnglish textDigest。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between gap-4">
            <Label htmlFor="noa-compaction" className="text-sm font-normal text-muted-foreground">
              noa English textDigest{noaCompaction ? " · English text" : " · English textClose"}
            </Label>
            <Switch
              id="noa-compaction"
              checked={noaCompaction}
              disabled={!loaded}
              onCheckedChange={toggleNoaCompaction}
            />
          </CardContent>
        </Card>

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <SearchIcon className="size-4" />
              English textSearch
            </CardTitle>
            <CardDescription>
              English textSearchEnglish text<b>English text + SourceConfiguration</b>。English text，English text<b>English text Agent English textConfiguration</b>English text
              <b>web_search</b>（English text/English text/Summary，English textBody；English text WebFetch English text）。English textSearch<b>English text</b>
              English text，English textTrafficEnglish text。
              <br />
              SourceEnglish text <b>DuckDuckGo（ddgs）</b>（English text Key）、<b>Brave（English text）</b>（English text Brave API Key）、{" "}
              <b>Tavily</b>（English text Tavily API Key）English text <b>DeepSeek</b>（English text LLM Configuration）。English textCloseEnglish text，English text
              Agent English textSearchEnglish text。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-4">
            <div className="flex items-center justify-between gap-4">
              <Label htmlFor="web-search" className="text-sm font-normal text-muted-foreground">
                {webSearch ? "English text · English text Agent ConfigurationEnglish text" : "English textClose · English text Agent English textSearch"}
              </Label>
              <Switch
                id="web-search"
                checked={webSearch}
                disabled={!loaded || saving}
                onCheckedChange={(v) => {
                  setWebSearch(v); // optimistic
                  saveWebSearch({ web_search_enabled: v });
                }}
              />
            </div>

            {webSearch && (
              <div className="flex items-center justify-between gap-4">
                <Label className="text-sm font-normal text-muted-foreground">SearchSource</Label>
                <Select
                  value={backend}
                  disabled={!loaded || saving}
                  onValueChange={(v) => {
                    setBackend(v); // optimistic
                    saveWebSearch({ web_search_backend: v });
                  }}
                >
                  <SelectTrigger className="w-48 shrink-0">
                    <SelectValue placeholder="English textSource" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="ddgs">DuckDuckGo（ddgs · English text Key）</SelectItem>
                    <SelectItem value="brave-free">Brave（English text · English text Key）</SelectItem>
                    <SelectItem value="tavily">Tavily（English text Key）</SelectItem>
                    <SelectItem value="deepseek">DeepSeek（English text）</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            )}

            {webSearch && backend === "deepseek" && (
              <div className="border-border/60 bg-muted/30 flex flex-col gap-2 rounded-md border p-3">
                <p className="text-sm font-medium">DeepSeek English textSearch</p>
                <p className="text-muted-foreground text-xs leading-relaxed">
                  English textSourceEnglish text<b>English text LLM Configuration</b>。English text
                  <b>English text DeepSeek English text</b>，English textConfiguration<b>English text anthropic English text</b>
                  ——DeepSeek English text OpenAI English textEndpointEnglish textServiceEnglish textSearch。English text LLM ConfigurationEnglish textSourceEnglish text。
                </p>
                <p className="text-muted-foreground text-xs leading-relaxed">
                  English textSourceEnglish text，SearchEnglish text <b>DeepSeek ServiceEnglish text</b>：English textSearchEnglish text（English text Token
                  English text），SearchEnglish text<b>English text</b>，English text<b>English textTrafficEnglish text</b>；English text<b>English text</b>
                  （English textSummary），English textBodyEnglish text WebFetch English text。
                </p>
                <p className="text-muted-foreground text-xs leading-relaxed">
                  English text itemsEnglish text，SystemEnglish text；English text「English textSearch」English text。
                </p>
              </div>
            )}

            {webSearch && backend === "brave-free" && (
              <div className="flex flex-col gap-2">
                <Label htmlFor="brave-key" className="text-sm font-normal text-muted-foreground">
                  Brave Search API Key
                  {braveKeySet && <span className="ml-2 text-xs text-emerald-500">English textConfiguration</span>}
                </Label>
                <div className="flex items-center gap-2">
                  <Input
                    id="brave-key"
                    type="password"
                    autoComplete="off"
                    placeholder={braveKeySet ? "English textConfiguration（English text）" : "Input Brave API Key"}
                    value={braveKeyInput}
                    disabled={!loaded || savingKey}
                    onChange={(e) => setBraveKeyInput(e.target.value)}
                  />
                  <Button
                    type="button"
                    onClick={saveBraveKey}
                    disabled={!loaded || savingKey || braveKeyInput.trim() === ""}
                  >
                    Save
                  </Button>
                </div>
                {braveNeedsKey && (
                  <p className="text-xs text-amber-500">
                    English text Brave English textConfiguration Key —— English textSave Key English text，SearchToolsEnglish text。
                  </p>
                )}
                <p className="text-muted-foreground text-xs">
                  English text 2,000 English text/English text。English text https://brave.com/search/api/ English text Key。
                </p>
              </div>
            )}

            {webSearch && backend === "tavily" && (
              <div className="flex flex-col gap-2">
                <Label htmlFor="tavily-key" className="text-sm font-normal text-muted-foreground">
                  Tavily Search API Key
                  {tavilyKeySet && <span className="ml-2 text-xs text-emerald-500">English textConfiguration</span>}
                </Label>
                <div className="flex items-center gap-2">
                  <Input
                    id="tavily-key"
                    type="password"
                    autoComplete="off"
                    placeholder={tavilyKeySet ? "English textConfiguration（English text）" : "Input Tavily API Key（tvly-…）"}
                    value={tavilyKeyInput}
                    disabled={!loaded || savingTavilyKey}
                    onChange={(e) => setTavilyKeyInput(e.target.value)}
                  />
                  <Button
                    type="button"
                    onClick={saveTavilyKey}
                    disabled={!loaded || savingTavilyKey || tavilyKeyInput.trim() === ""}
                  >
                    Save
                  </Button>
                </div>
                {webSearch && backend === "tavily" && !tavilyKeySet && (
                  <p className="text-xs text-amber-500">
                    English text Tavily English textConfiguration Key —— English textSave Key English text，SearchToolsEnglish text。
                  </p>
                )}
                <p className="text-muted-foreground text-xs">English text https://tavily.com English text API Key。</p>
              </div>
            )}

            {webSearch && (
              <div className="flex flex-col gap-2">
                <Label htmlFor="ws-proxy" className="text-sm font-normal text-muted-foreground">
                  English text（English text）
                </Label>
                <div className="flex items-center gap-2">
                  <Input
                    id="ws-proxy"
                    autoComplete="off"
                    placeholder="http://host:port English text socks5://host:port（English text=English text）"
                    value={proxyInput}
                    disabled={!loaded || savingProxy}
                    onChange={(e) => setProxyInput(e.target.value)}
                  />
                  <Button type="button" onClick={saveProxy} disabled={!loaded || savingProxy}>
                    Save
                  </Button>
                </div>
                <p className="text-muted-foreground text-xs">
                  English text，English textSearchEndpoint（VPN/SOCKS English text）。English textTrafficEnglish text MITM English text；English text。
                </p>
              </div>
            )}

            {webSearch && (
              <div className="flex items-center justify-between gap-4 border-t pt-4">
                <p className="text-muted-foreground text-xs">
                  English textConfiguration（Source + English text + Key）English textSearchEnglish text「test」，English text。
                </p>
                <Button
                  type="button"
                  variant="outline"
                  onClick={runTest}
                  disabled={!loaded || testing}
                  className="shrink-0"
                >
                  {testing ? "English text…" : "English textSearch"}
                </Button>
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <RadioTowerIcon className="size-4" />
              English text · Python English text
            </CardTitle>
            <CardDescription>
              English text <b>script</b> TypeToolsEnglish text Python。English text（python3 English text）；English text venv /
              English text，English text。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-3">
            <div className="flex items-center gap-2">
              <Input
                className="font-mono text-sm"
                placeholder="/usr/bin/python3（English text=English text）"
                value={pyInterp}
                disabled={!loaded || saving}
                onChange={(e) => setPyInterp(e.target.value)}
              />
              <Button variant="outline" onClick={detectPython} disabled={!loaded || saving}>
                English text
              </Button>
              <Button onClick={savePython} disabled={!loaded || saving}>
                Save
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <CpuIcon className="size-4" />
              English text · Work Agent English text
            </CardTitle>
            <CardDescription>
              English textTasksEnglish text agent English text（Default 3）。English text、English text。English text
              <b>English textTasksEnglish text</b>，English textTasksEnglish text。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col gap-3">
            <div className="flex items-center gap-2">
              <Input
                type="number"
                min={1}
                className="w-32 font-mono text-sm"
                placeholder="3"
                value={workers}
                disabled={!loaded || savingWorkers}
                onChange={(e) => setWorkers(e.target.value)}
              />
              <Button onClick={saveWorkers} disabled={!loaded || savingWorkers}>
                Save
              </Button>
            </div>
          </CardContent>
        </Card>

        <Card className="mb-4 break-inside-avoid md:mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <KeyboardIcon className="size-4" />
              English textInputEnglish textSendEnglish text
            </CardTitle>
            <CardDescription>
              ChatEnglish textTasksEnglish text Agent English textInputEnglish text，English text、English textSave。
              <br />
              English text<b>English text</b>，English text，English text。
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between gap-4">
            <Label htmlFor="chat-send-mode" className="text-sm font-normal text-muted-foreground">
              SendEnglish text
            </Label>
            <Select value={sendMode} onValueChange={(v) => setChatSendMode(v as ChatSendMode)}>
              <SelectTrigger id="chat-send-mode" className="w-72">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {CHAT_SEND_MODE_OPTIONS.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
