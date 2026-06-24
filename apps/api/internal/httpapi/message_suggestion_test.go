package httpapi

import (
	"testing"
	"time"
)

func TestInferMessageSuggestionInSpanish(t *testing.T) {
	now := time.Date(2026, time.June, 22, 18, 0, 0, 0, time.UTC)
	if _, _, ok := inferMessageTime("¿Puedes pasear a Tomy mañana a las 10 por 45 minutos?"); !ok {
		t.Fatal("expected message time")
	}
	if _, ok := inferMessageDay("¿Puedes pasear a Tomy mañana a las 10 por 45 minutos?", now); !ok {
		t.Fatal("expected message day")
	}
	suggestion := inferredMessageSuggestion("¿Puedes pasear a Tomy mañana a las 10 por 45 minutos?", now)
	if suggestion == nil || suggestion.ServiceType != "walk" || suggestion.Confidence != "high" {
		t.Fatalf("unexpected suggestion: %#v", suggestion)
	}
	if suggestion.StartAt == nil || *suggestion.StartAt != "2026-06-23T10:00:00-06:00" {
		t.Fatalf("unexpected start: %#v", suggestion.StartAt)
	}
	if suggestion.EndAt == nil || *suggestion.EndAt != "2026-06-23T10:45:00-06:00" {
		t.Fatalf("unexpected end: %#v", suggestion.EndAt)
	}
}

func TestInferMessageSuggestionDoesNotInventTime(t *testing.T) {
	suggestion := inferredMessageSuggestion("Quiero un paseo cuando puedas", time.Now())
	if suggestion == nil || suggestion.ServiceType != "walk" || suggestion.StartAt != nil {
		t.Fatalf("expected service-only suggestion, got %#v", suggestion)
	}
}
