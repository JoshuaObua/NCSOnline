package handlers

import (
	"encoding/json"
	"testing"
)

func TestNormalizeSlidePayloadAcceptsCMSAnimationModes(t *testing.T) {
	for animationType := range slideAnimationTypes {
		req := &slideReq{
			Title:         "Test slide",
			MediaType:     "image",
			AnimationType: animationType,
			IsActive:      true,
		}
		slide, errs := normalizeSlidePayload(req)
		if len(errs) != 0 {
			t.Fatalf("animation %q was rejected: %v", animationType, errs)
		}
		if slide.AnimationType != animationType {
			t.Fatalf("animation %q was not preserved", animationType)
		}
	}
}

func TestNormalizeSlidePayloadRejectsUnknownAnimation(t *testing.T) {
	req := &slideReq{Title: "Test slide", MediaType: "image", AnimationType: "spin-forever"}
	if _, errs := normalizeSlidePayload(req); errs["animation_type"] == "" {
		t.Fatal("unknown animation type was accepted")
	}
}

func TestCMSSettingsValidationAndSecretRedaction(t *testing.T) {
	settings, errs := validateThirdPartySettings(thirdPartySettings{
		GoogleAnalyticsEnabled: true,
		GoogleAnalyticsID:      " g-ab12cd34 ",
	})
	if len(errs) != 0 || settings.GoogleAnalyticsID != "G-AB12CD34" {
		t.Fatalf("unexpected analytics settings: %#v, errors: %#v", settings, errs)
	}

	public := publicCaptchaSettings(captchaSettings{
		Provider:              "cloudflare_turnstile",
		CloudflareSecretKey:   "private-secret",
		RecaptchaSecretKey:    "another-secret",
		CloudflareSecretSaved: false,
	})
	if public.CloudflareSecretKey != "" || public.RecaptchaSecretKey != "" {
		t.Fatal("captcha secrets were exposed in the public settings payload")
	}
	if !public.CloudflareSecretSaved || !public.RecaptchaSecretSaved {
		t.Fatal("saved-secret indicators were not preserved")
	}
}

func TestInboundPayloadAlwaysReturnsValidJSON(t *testing.T) {
	valid := inboundPayload(map[string]interface{}{"channel": "web", "status": "unread"})
	if !json.Valid(valid) {
		t.Fatalf("valid inbound payload became invalid JSON: %s", valid)
	}

	fallback := inboundPayload(map[string]interface{}{"unsupported": make(chan int)})
	if string(fallback) != "{}" {
		t.Fatalf("unserializable payload fallback = %s, want {}", fallback)
	}
}

func TestSanitizePlainRemovesMarkupAndControlCharacters(t *testing.T) {
	got := sanitizePlain(" <script>alert(1)</script>\x00 hello ", 20)
	if got == "" || got == "<script>alert(1)</script> hello" {
		t.Fatalf("plain-text sanitizer did not clean unsafe input: %q", got)
	}
	if len(got) > 20 {
		t.Fatalf("sanitized value exceeded limit: %q", got)
	}
}
