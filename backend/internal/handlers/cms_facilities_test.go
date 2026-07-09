package handlers

import "testing"

func TestNormalizeFacilitySeparatesRegionAndCategory(t *testing.T) {
	facility, errs := normalizeFacility(facilityReq{
		Name:               "  Hoima City Stadium  ",
		Category:           "Stadium",
		Region:             "Western Region",
		Location:           "Hoima City",
		AvailabilityStatus: "Available",
		IsActive:           true,
	})
	if len(errs) != 0 {
		t.Fatalf("valid facility rejected: %v", errs)
	}
	if facility.Slug != "hoima-city-stadium" {
		t.Fatalf("slug = %q", facility.Slug)
	}
	if facility.Region != "western" {
		t.Fatalf("region = %q", facility.Region)
	}
	if facility.Category != "Stadium" {
		t.Fatalf("category = %q", facility.Category)
	}
}

func TestNormalizeFacilityRejectsMissingRegionAndInvalidStatus(t *testing.T) {
	_, errs := normalizeFacility(facilityReq{
		Name:               "Test Facility",
		AvailabilityStatus: "Open sometimes",
	})
	if errs["region"] == "" || errs["availability_status"] == "" {
		t.Fatalf("expected region and status errors, got %v", errs)
	}
}

func TestNormalizeFacilityRegionBuildsStableSlug(t *testing.T) {
	region, errs := normalizeFacilityRegion(facilityRegionReq{
		Name:     "  West Nile Region ",
		IsActive: true,
	})
	if len(errs) != 0 {
		t.Fatalf("valid region rejected: %v", errs)
	}
	if region.Name != "West Nile Region" || region.Slug != "west-nile" {
		t.Fatalf("unexpected normalized region: %#v", region)
	}
}
