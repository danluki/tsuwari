package moderationhelpers

import (
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

const denyUnicodeAwareBoundaryPrefix = `(?:^|\s|[^\p{L}\p{N}])`
const denyUnicodeAwareBoundarySuffix = `(?:$|\s|[^\p{L}\p{N}])`

type HasDeniedWordInput struct {
	Message             string
	RulesList           []string
	RegexpEnabled       bool
	WordBoundaryEnabled bool
	SensitivityEnabled  bool
}

func (c *ModerationHelpers) HasDeniedWord(input HasDeniedWordInput) bool {
	msg := input.Message
	if !input.SensitivityEnabled {
		msg = strings.ToLower(msg)
	}

	for _, rule := range input.RulesList {
		if rule == "" {
			continue
		}

		if !input.SensitivityEnabled {
			rule = strings.ToLower(rule)
		}

		// if regexp enabled - we handle regexp and just go through other words
		if input.RegexpEnabled {
			r, err := regexp2.Compile(rule, regexp2.RE2)
			if err == nil {
				matched, err := r.MatchString(msg)
				if matched && err == nil {
					return true
				}
			}
			continue
		}

		if !input.WordBoundaryEnabled {
			if strings.Contains(msg, rule) {
				return true
			}
			continue
		}

		wordRg := regexp.MustCompile(denyUnicodeAwareBoundaryPrefix + regexp.QuoteMeta(strings.ToLower(rule)) + denyUnicodeAwareBoundarySuffix)
		matched := wordRg.MatchString(msg)

		if matched {
			return true
		}
	}

	return false
}
