package secrets

import (
	"fmt"
	"os"
	"strings"

	"github.com/zricethezav/gitleaks/v8/detect"
)

// Finding stores information about a potential secret found by the scanner.
type Finding struct {
	RuleID string
	Match  string
	Secret string
}

// Scanner defines the interface for secret scanning operations.
type Scanner interface {
	Scan(content string) ([]Finding, error)
	Redact(content string, findings []Finding) string
}

// GitleaksScanner implements the Scanner interface using the gitleaks library.
type GitleaksScanner struct {
	detector *detect.Detector
}

// NewGitleaksScanner creates a scanner using gitleaks' embedded default configuration.
func NewGitleaksScanner() (*GitleaksScanner, error) {
	detector, err := detect.NewDetectorDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gitleaks detector with default config: %w", err)
	}
	return &GitleaksScanner{detector: detector}, nil
}

// Scan uses the gitleaks detector to find secrets in the content.
func (s *GitleaksScanner) Scan(content string) ([]Finding, error) {
	gitleaksFindings := s.detector.DetectString(content)
	if len(gitleaksFindings) == 0 {
		return nil, nil
	}

	findings := make([]Finding, 0, len(gitleaksFindings))
	// processedFindings prevents processing the same redaction target multiple times.
	processedFindings := make(map[string]struct{})

	for _, glFinding := range gitleaksFindings {
		if glFinding.Secret == "" || glFinding.RuleID == "" || glFinding.Match == "" {
			continue
		}

		// Ensure the reported Secret is actually within the reported Match string.
		if !strings.Contains(glFinding.Match, glFinding.Secret) {
			fmt.Fprintf(os.Stderr, "Warning: Secret '%s' not found within Match '%s' for rule %s. Cannot perform value-only redaction. Skipping.\n", shorten(glFinding.Secret, 15), shorten(glFinding.Match, 15), glFinding.RuleID)
			continue
		}

		uniqueKey := fmt.Sprintf("%s::%s::%s", glFinding.RuleID, glFinding.Match, glFinding.Secret)
		if _, exists := processedFindings[uniqueKey]; exists {
			continue
		}

		findings = append(findings, Finding{
			RuleID: glFinding.RuleID,
			Match:  glFinding.Match,
			Secret: glFinding.Secret,
		})
		processedFindings[uniqueKey] = struct{}{}
	}

	if len(findings) == 0 {
		return nil, nil
	}
	return findings, nil
}

// Redact replaces *only the Secret part* within each occurrence of a Match string.
func (s *GitleaksScanner) Redact(content string, findings []Finding) string {
	if len(findings) == 0 {
		return content
	}

	redactedContent := content

	// Process unique Match strings to avoid redundant searches.
	uniqueReplacements := make(map[string]Finding)
	for _, f := range findings {
		// Using Match as the key is a simplification. If multiple rules match the
		// same text but identify different Secret subgroups, this picks one arbitrarily.
		if _, exists := uniqueReplacements[f.Match]; !exists {
			uniqueReplacements[f.Match] = f
		}
	}

	for _, f := range uniqueReplacements {
		placeholder := fmt.Sprintf("[REDACTED_%s]", f.RuleID)

		startIndex := 0
		for {
			index := strings.Index(redactedContent[startIndex:], f.Match)
			if index == -1 {
				break
			}
			absoluteMatchStart := startIndex + index

			// Guard against index out of bounds if content changed drastically.
			matchEnd := absoluteMatchStart + len(f.Match)
			if matchEnd > len(redactedContent) {
				fmt.Fprintf(os.Stderr, "Warning (Redact): Match '%s' extends beyond content length after previous redactions. Stopping search for this match.\n", shorten(f.Match, 15))
				break
			}
			matchText := redactedContent[absoluteMatchStart:matchEnd]

			relativeSecretIndex := strings.Index(matchText, f.Secret)

			if relativeSecretIndex != -1 {
				absoluteSecretStart := absoluteMatchStart + relativeSecretIndex
				absoluteSecretEnd := absoluteSecretStart + len(f.Secret)

				redactedContent = redactedContent[:absoluteSecretStart] + placeholder + redactedContent[absoluteSecretEnd:]

				// Move startIndex past the *placeholder* to avoid reprocessing or infinite loops.
				startIndex = absoluteSecretStart + len(placeholder)

			} else {
				// Secret wasn't found within this specific Match instance (unexpected).
				fmt.Fprintf(os.Stderr, "Warning (Redact): Secret '%s' not found within matched segment '%s' at index %d for rule %s.\n", shorten(f.Secret, 15), shorten(f.Match, 15), absoluteMatchStart, f.RuleID)
				startIndex = absoluteMatchStart + 1 // Advance past the start of this problematic match.
			}

			if startIndex >= len(redactedContent) {
				break
			}
		}
	}

	return redactedContent
}

// shorten truncates a string for cleaner logging.
func shorten(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
