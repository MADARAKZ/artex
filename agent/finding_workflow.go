package agent

import (
	"encoding/json"
	"fmt"

	"github.com/Autumn-27/artex/db"
	actool "github.com/Autumn-27/norma/tool"
)

const findingIDGuidance = "\n\n**English text**：finding_id English text ID；finding_node_id English text ID。list_findings / list_task_findings / node_detail / get_task_node_detail English text id English text ID，English text finding_id English text。get_finding_traffic / bind_finding_traffic English text finding_id。English text update_finding_report English text finding_id English text finding_node_id。English text report_finding English text，English text。"

// The server supplies the persisted setting. A missing setting/host is off.
// Consulted at assembly and again on writes so an already-running session
// cannot keep binding after the user switches the feature off.
var FindingTrafficBindingEnabled func() bool

func findingTrafficBindingEnabled() bool {
	return FindingTrafficBindingEnabled != nil && FindingTrafficBindingEnabled()
}

// Applied after ToolResolve: user descriptions and prompts remain intact, while
// all actual reporters (including Planner and custom chat agents) see the same
// API contract. Disabled/unbound tools are never reintroduced here.
func findingWorkflowTools(agentKey string, tools []actool.CoreTool) ([]actool.CoreTool, string) {
	if !findingTrafficBindingEnabled() {
		out := make([]actool.CoreTool, 0, len(tools))
		for _, tool := range tools {
			if tool.Name() == "bind_finding_traffic" {
				continue
			}
			if agentKey == "reporter" && (tool.Name() == "traffic_search" || tool.Name() == "traffic_get" || tool.Name() == "traffic_blob") {
				continue
			}
			switch tool.Name() {
			case "report_finding", "add_hint", "add_task_hint":
				// Work on a copy: toggling back on must restore the original schema.
				raw, _ := json.Marshal(tool.InputSchema())
				var schema map[string]any
				if json.Unmarshal(raw, &schema) == nil {
					stripTrafficParameters(schema)
					tool = DecorateTool(tool, tool.Description(), schema)
				}
			}
			out = append(out, tool)
		}
		return out, ""
	}
	out := append([]actool.CoreTool(nil), tools...)
	has := map[string]bool{}
	for i, tool := range out {
		has[tool.Name()] = true
		note := ""
		switch tool.Name() {
		case "report_finding":
			note = "\nEnglish text Agent English text。English text evidence English text、English text、English text ID English text，English text Agent English text；English text。English text：traffic_refs English text evidence_hint_id English text，English text hint English text；English text。TCP/English text。English text finding_id English text finding_node_id English text。"
		case "add_hint", "add_task_hint":
			note = "\nEnglish text，English text traffic_refs English text ID、English text、English text（English text，English text hints English text），English text text English text。English text。English text。"
		case "get_finding_traffic", "bind_finding_traffic", "list_findings", "list_task_findings", "node_detail", "get_task_node_detail", "update_finding_report":
			note = findingIDGuidance
		}
		if note != "" {
			out[i] = DecorateTool(tool, tool.Description()+note, tool.InputSchema())
		}
	}
	guidance := ""
	if has["report_finding"] || has["add_task_hint"] || has["add_hint"] {
		guidance = "\n\n**English text（English text）**：English text Agent English text、English text。English text evidence English text、English text、English text ID English text，English text intent_id，English text Agent English text；English text。Auto / Planner English text。add_hint / add_task_hint English text traffic_refs English text；English text report_finding English text traffic_refs / evidence_hint_id。TCP English text，English text ID，English text。"
		if has["add_task_hint"] && !has["add_hint"] {
			guidance += "\nEnglish text，English text report_finding；English text add_task_hint English text，English text Agent English text，English text list_task_findings English text。"
		}
		if has["prove_goal"] || has["goal_met"] {
			guidance += "\nEnglish text，English text/English text。English text、English text Worker；English text。"
		}
	}
	if has["update_finding_report"] && has["bind_finding_traffic"] && has["get_finding_traffic"] {
		guidance += "\n\n**English text（English text）**：English text，English text。English text report_finding English text JSON English text get_task_node_detail / list_task_findings English text finding_id English text finding_node_id。English text、English text，English text ID。English text HTTP English text，English text traffic_search English text，English text traffic_get English text/English text；English text，English text。English text bind_finding_traffic(finding_id, traffic_refs) English text，English text baseline / proof / verification / supporting English text。English text，English text。English text get_finding_traffic English text version，English text，English text version English text evidence_version English text update_finding_report（English text finding_id English text finding_node_id）。English text。TCP、English text、English text，English text，English text/English text，English text。English text；English text。"
	}
	if guidance != "" || has["get_finding_traffic"] || has["update_finding_report"] {
		guidance += findingIDGuidance
	}
	return out, guidance
}

func stripTrafficParameters(schema map[string]any) {
	props, _ := schema["properties"].(map[string]any)
	delete(props, "traffic_refs")
	delete(props, "evidence_hint_id")
	if required, ok := schema["required"].([]any); ok {
		kept := required[:0]
		for _, key := range required {
			if key != "traffic_refs" && key != "evidence_hint_id" {
				kept = append(kept, key)
			}
		}
		schema["required"] = kept
	}
	if hints, ok := props["hints"].(map[string]any); ok {
		if items, ok := hints["items"].(map[string]any); ok {
			stripTrafficParameters(items)
		}
	}
}

// HintTrafficSchema is shared by the task-local and cross-task hint tools.
func HintTrafficSchema() map[string]any {
	return map[string]any{"type": "array", "description": "English text：English text，English text；English text report_finding English text evidence_hint_id English text。", "items": obj(map[string]any{"traffic_id": str("English text ID"), "role": str("baseline / proof / verification / supporting"), "note": str("English text")}, "traffic_id")}
}

func (t *ToolSet) findingRefsFromHint(hintID int64, explicit []db.TrafficRef) ([]db.TrafficRef, error) {
	if hintID <= 0 {
		return db.NormalizeTrafficRefs(explicit)
	}
	n, err := t.ts.GetNode(hintID) // local store only: inherited hints cannot supply evidence
	if err != nil {
		return nil, err
	}
	if n == nil || n.Kind != db.KindHint {
		return nil, fmt.Errorf("evidence_hint_id=%d English text（English text）", hintID)
	}
	var payload struct {
		Refs []db.TrafficRef `json:"traffic_refs"`
	}
	if err := json.Unmarshal(n.Payload, &payload); err != nil {
		return nil, err
	}
	return db.NormalizeTrafficRefs(append(append([]db.TrafficRef{}, explicit...), payload.Refs...))
}
