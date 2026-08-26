package template

import (
	"fmt"
	"html"
	"strings"
)

func VerificationCodeTemplate(name, code string) string {
	return fmt.Sprintf(
		"<p>Hello %s,</p><p>Your GoTask verification code is:</p><p><strong>%s</strong></p><p>This code expires in 10 minutes.</p>",
		escape(name),
		escape(code),
	)
}

func EmailChangedTemplate(name, newEmail string) string {
	return fmt.Sprintf(
		"<p>Hello %s,</p><p>Your GoTask account email was changed to <strong>%s</strong>.</p><p>If you did not make this change, please contact support immediately.</p>",
		escape(name),
		escape(newEmail),
	)
}

func escape(value string) string {
	return html.EscapeString(strings.TrimSpace(value))
}
