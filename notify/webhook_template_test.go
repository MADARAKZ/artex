package notify

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// English text Webhook English text**English text**。
//
// English text「English text」English text，English text、
// English text，English text——English text
// English text、English text FuncMap English text readFile，English text，English text diff English text
// English text。

// TestTemplateContextHasNoMethods English text。
//
// text/template English text（{{.Foo}} English text）。English text
// English text**English text**English text，English text。
// English text（English text、English text）。
//
// English text：English text webhookTemplateData / webhookItem English text。
// English text，English text。
func TestTemplateContextHasNoMethods(t *testing.T) {
	for _, v := range []any{webhookTemplateData{}, webhookItem{}} {
		typ := reflect.TypeOf(v)
		if n := typ.NumMethod(); n != 0 {
			var names []string
			for i := 0; i < n; i++ {
				names = append(names, typ.Method(i).Name)
			}
			t.Fatalf("%s English text %d English text（%s）：text/template English text，"+
				"English text", typ.Name(), n, strings.Join(names, ", "))
		}
	}
}

// TestTemplateFuncsAreMinimal English text。
//
// FuncMap English text。English text json / jsons，English text
// English text JSON English text——English text、English text、English text。
func TestTemplateFuncsAreMinimal(t *testing.T) {
	var got []string
	for name := range webhookTemplateFuncs {
		got = append(got, name)
	}
	sort.Strings(got)
	want := []string{"json", "jsons"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("English text：English text %v，English text %v。English text"+
			"（English text、English text、English text）", got, want)
	}
}

// TestTemplateCannotReachUnknownData English text：
// English text，English text；English text。
func TestTemplateCannotReachUnknownData(t *testing.T) {
	_, err := renderWebhookBody(`{"x": {{.Environment}}, "y": {{.Env}}}`, singleMsg())
	if err == nil {
		t.Fatal("English text")
	}
	// English text（English text/English text）。
	for _, leak := range []string{"SQLEnglish text", "English text id"} {
		if strings.Contains(err.Error(), leak) {
			t.Errorf("English text %q: %v", leak, err)
		}
	}
}

// TestTemplateRenderFailsPermanently English text，English text。
// English text，English text。
func TestTemplateRenderFailsPermanently(t *testing.T) {
	cfg := map[string]any{
		"url":           "https://example.com/hook",
		"body_template": `{{.Items.`,
	}
	if err := (webhookChannel{}).Validate(cfg); err == nil {
		t.Fatal("English text")
	}
	// English text，English text。
	_, err := (webhookChannel{}).Send(context.Background(), cfg, singleMsg())
	if err == nil || !IsPermanent(err) {
		t.Fatalf("English text，English text %v", err)
	}
}

// TestTemplateCanOnlyProduceJSON English text「English text JSON」English text。
// English text「English text」English text。
func TestTemplateCanOnlyProduceJSON(t *testing.T) {
	// English text。
	ok := map[string]any{"url": "https://example.com/hook", "body_template": `{"t":{{json .Title}}}`}
	if err := (webhookChannel{}).Validate(ok); err != nil {
		t.Fatalf("English text: %v", err)
	}
	// English text JSON English text（English text）。
	bad := map[string]any{"url": "http://127.0.0.1:1/hook", "body_template": `not json {{.Count}}`}
	_, err := (webhookChannel{}).Send(context.Background(), bad, singleMsg())
	if err == nil || !IsPermanent(err) {
		t.Fatalf("English text JSON English text，English text %v", err)
	}
	if !strings.Contains(err.Error(), "English text JSON") {
		t.Errorf("English text JSON English text，English text %v", err)
	}
}
