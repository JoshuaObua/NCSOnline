package handlers

import "testing"

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
