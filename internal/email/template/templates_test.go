package template

import (
	"strings"
	"testing"
)

func TestVerificationCodeTemplateEscapesUserInput(t *testing.T) {
	body := VerificationCodeTemplate("<Ada>", "123<&")
	if strings.Contains(body, "<Ada>") || strings.Contains(body, "123<&") {
		t.Fatalf("template did not escape input: %s", body)
	}
	if !strings.Contains(body, "&lt;Ada&gt;") || !strings.Contains(body, "123&lt;&amp;") {
		t.Fatalf("template is missing escaped input: %s", body)
	}
}

func TestEmailChangedTemplateIncludesNewEmail(t *testing.T) {
	body := EmailChangedTemplate("Ada", "new@example.com")
	if !strings.Contains(body, "new@example.com") {
		t.Fatalf("template does not include new email: %s", body)
	}
}
