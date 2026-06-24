package httpapi

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	messageTimePattern     = regexp.MustCompile(`(?i)(?:a\s+las?|a\s+la|las?)\s*(\d{1,2})(?::(\d{2}))?\s*(am|pm)?`)
	messageISODatePattern  = regexp.MustCompile(`\b(20\d{2})-(\d{1,2})-(\d{1,2})\b`)
	messageShortDate       = regexp.MustCompile(`\b(\d{1,2})/(\d{1,2})(?:/(20\d{2}))?\b`)
	messageDurationPattern = regexp.MustCompile(`(?i)(?:por|durante|duraci[oó]n\s*)?\s*(\d+)\s*(min|minuto|minutos|h|hora|horas)\b`)
)

func inferredMessageSuggestion(body string, now time.Time) *detectedSuggestionInput {
	service := inferServiceType(body)
	start, hasStart := inferMessageStart(body, now)
	var startAt, endAt *string
	if hasStart {
		formatted := start.Format(time.RFC3339)
		startAt = &formatted
		if duration, ok := inferMessageDuration(body); ok {
			formattedEnd := start.Add(duration).Format(time.RFC3339)
			endAt = &formattedEnd
		}
	}
	if service == "" && startAt == nil {
		return nil
	}
	confidence := "low"
	if service != "" && startAt != nil {
		confidence = "high"
	} else if service != "" || startAt != nil {
		confidence = "medium"
	}
	return &detectedSuggestionInput{
		ServiceType: service, StartAt: startAt, EndAt: endAt, Confidence: confidence,
	}
}

func mergeMessageSuggestion(explicit *detectedSuggestionInput, inferred *detectedSuggestionInput) *detectedSuggestionInput {
	if explicit == nil {
		return inferred
	}
	if inferred == nil {
		return explicit
	}
	merged := *explicit
	if merged.ServiceType == "" {
		merged.ServiceType = inferred.ServiceType
	}
	if merged.StartAt == nil {
		merged.StartAt = inferred.StartAt
	}
	if merged.EndAt == nil {
		merged.EndAt = inferred.EndAt
	}
	if merged.Confidence == "" || merged.Confidence == "unknown" {
		merged.Confidence = inferred.Confidence
	}
	return &merged
}

func inferServiceType(body string) string {
	normalized := normalizeMessageText(body)
	services := []struct {
		value    string
		keywords []string
	}{
		{"walk", []string{"paseo", "pasear", "caminar", "caminata"}},
		{"boarding", []string{"hospedaje", "hospedar", "alojamiento", "pension"}},
		{"client_home_sitting", []string{"cuidar en casa", "quedarse en casa", "pet sitting"}},
		{"visit", []string{"visita", "dar de comer", "alimentar", "medicina"}},
		{"transport", []string{"transportar", "traslado", "recoger", "llevar al vet", "veterinario"}},
	}
	for _, service := range services {
		for _, keyword := range service.keywords {
			if strings.Contains(normalized, keyword) {
				return service.value
			}
		}
	}
	return ""
}

func inferMessageStart(body string, now time.Time) (time.Time, bool) {
	location, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		location = time.FixedZone("CST", -6*60*60)
	}
	localNow := now.In(location)
	day, hasDay := inferMessageDay(body, localNow)
	hour, minute, hasTime := inferMessageTime(body)
	if !hasDay || !hasTime {
		return time.Time{}, false
	}
	return time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, location), true
}

func inferMessageDay(body string, now time.Time) (time.Time, bool) {
	normalized := normalizeMessageText(body)
	switch {
	case strings.Contains(normalized, "pasado manana"):
		return now.AddDate(0, 0, 2), true
	case strings.Contains(normalized, "manana"):
		return now.AddDate(0, 0, 1), true
	case strings.Contains(normalized, "hoy"):
		return now, true
	}
	if match := messageISODatePattern.FindStringSubmatch(body); match != nil {
		year, _ := strconv.Atoi(match[1])
		month, _ := strconv.Atoi(match[2])
		day, _ := strconv.Atoi(match[3])
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, now.Location()), true
	}
	if match := messageShortDate.FindStringSubmatch(body); match != nil {
		day, _ := strconv.Atoi(match[1])
		month, _ := strconv.Atoi(match[2])
		year := now.Year()
		if match[3] != "" {
			year, _ = strconv.Atoi(match[3])
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, now.Location()), true
	}
	return time.Time{}, false
}

func inferMessageTime(body string) (int, int, bool) {
	match := messageTimePattern.FindStringSubmatch(body)
	if match == nil {
		return 0, 0, false
	}
	hour, _ := strconv.Atoi(match[1])
	minute := 0
	if match[2] != "" {
		minute, _ = strconv.Atoi(match[2])
	}
	period := strings.ToLower(match[3])
	if period == "pm" && hour < 12 {
		hour += 12
	}
	if period == "am" && hour == 12 {
		hour = 0
	}
	if hour > 23 || minute > 59 {
		return 0, 0, false
	}
	return hour, minute, true
}

func inferMessageDuration(body string) (time.Duration, bool) {
	match := messageDurationPattern.FindStringSubmatch(body)
	if match == nil {
		return 0, false
	}
	amount, _ := strconv.Atoi(match[1])
	if amount <= 0 {
		return 0, false
	}
	if strings.HasPrefix(strings.ToLower(match[2]), "h") {
		return time.Duration(amount) * time.Hour, true
	}
	return time.Duration(amount) * time.Minute, true
}

func normalizeMessageText(value string) string {
	value = strings.ToLower(value)
	replacer := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u", "ñ", "n")
	return replacer.Replace(value)
}
