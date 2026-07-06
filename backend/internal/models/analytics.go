package models

import "time"

type AnalyticsEvent struct {
	SessionID              string    `json:"session_id"`
	VisitorHash            string    `json:"-"`
	EventName              string    `json:"event_name"`
	Path                   string    `json:"path"`
	Referrer               string    `json:"referrer,omitempty"`
	Source                 string    `json:"source"`
	Country                string    `json:"country"`
	Region                 string    `json:"region,omitempty"`
	City                   string    `json:"city,omitempty"`
	Browser                string    `json:"browser"`
	OS                     string    `json:"os"`
	DeviceType             string    `json:"device_type"`
	ScreenResolution       string    `json:"screen_resolution,omitempty"`
	Language               string    `json:"language,omitempty"`
	SessionDurationSeconds int       `json:"session_duration_seconds"`
	IsBounce               bool      `json:"is_bounce"`
	CreatedAt              time.Time `json:"created_at"`
}

type AnalyticsMetricRow struct {
	Label string  `json:"label"`
	Value int64   `json:"value"`
	Share float64 `json:"share"`
}

type AnalyticsTimePoint struct {
	Date           string  `json:"date"`
	Label          string  `json:"label"`
	Views          int64   `json:"views"`
	UniqueVisitors int64   `json:"unique_visitors"`
	BounceRate     float64 `json:"bounce_rate"`
}

type AnalyticsDashboard struct {
	RangeDays           int                  `json:"range_days"`
	TotalViews          int64                `json:"total_views"`
	UniqueVisitors      int64                `json:"unique_visitors"`
	BounceRate          float64              `json:"bounce_rate"`
	AverageSessionSecs  float64              `json:"average_session_seconds"`
	Timeline            []AnalyticsTimePoint `json:"timeline"`
	TopCountries        []AnalyticsMetricRow `json:"top_countries"`
	TopRegions          []AnalyticsMetricRow `json:"top_regions"`
	TopCities           []AnalyticsMetricRow `json:"top_cities"`
	Browsers            []AnalyticsMetricRow `json:"browsers"`
	OperatingSystems    []AnalyticsMetricRow `json:"operating_systems"`
	DeviceTypes         []AnalyticsMetricRow `json:"device_types"`
	TopPages            []AnalyticsMetricRow `json:"top_pages"`
	TopEntryPages       []AnalyticsMetricRow `json:"top_entry_pages"`
	TopExitPages        []AnalyticsMetricRow `json:"top_exit_pages"`
	AcquisitionChannels []AnalyticsMetricRow `json:"acquisition_channels"`
	GeneratedAt         time.Time            `json:"generated_at"`
}
