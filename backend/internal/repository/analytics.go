package repository

import (
	"context"
	"math"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepo struct{ db *pgxpool.Pool }

func (r *AnalyticsRepo) InsertEvent(ctx context.Context, e *models.AnalyticsEvent) error {
	const q = `INSERT INTO page_views_raw (
		session_id, visitor_hash, event_name, path, referrer, source, country, region, city,
		browser, os, device_type, screen_resolution, language, session_duration_seconds, is_bounce
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err := r.db.Exec(ctx, q,
		e.SessionID, e.VisitorHash, e.EventName, e.Path, e.Referrer, e.Source, e.Country, e.Region, e.City,
		e.Browser, e.OS, e.DeviceType, e.ScreenResolution, e.Language, e.SessionDurationSeconds, e.IsBounce,
	)
	return err
}

func (r *AnalyticsRepo) Dashboard(ctx context.Context, days int) (*models.AnalyticsDashboard, error) {
	if days <= 0 || days > 365 {
		days = 30
	}
	start := time.Now().UTC().AddDate(0, 0, -(days - 1)).Truncate(24 * time.Hour)
	out := &models.AnalyticsDashboard{RangeDays: days, GeneratedAt: time.Now().UTC()}
	const overviewQ = `SELECT
		COUNT(*) FILTER (WHERE event_name='page_view'),
		COUNT(DISTINCT visitor_hash) FILTER (WHERE event_name='page_view'),
		COUNT(*) FILTER (WHERE event_name='page_exit' AND is_bounce),
		COALESCE(AVG(NULLIF(session_duration_seconds, 0)) FILTER (WHERE event_name='page_exit'), 0)
		FROM page_views_raw WHERE created_at >= $1`
	var bounces int64
	if err := r.db.QueryRow(ctx, overviewQ, start).Scan(&out.TotalViews, &out.UniqueVisitors, &bounces, &out.AverageSessionSecs); err != nil {
		return nil, err
	}
	if out.TotalViews > 0 {
		out.BounceRate = roundPercent(float64(bounces) / float64(out.TotalViews) * 100)
	}
	timeline, err := r.timeline(ctx, start, days)
	if err != nil {
		return nil, err
	}
	out.Timeline = timeline
	if out.TopCountries, err = r.breakdown(ctx, start, "country", "event_name='page_view' AND country <> ''", 8); err != nil {
		return nil, err
	}
	if out.TopRegions, err = r.breakdown(ctx, start, "region", "event_name='page_view' AND region <> ''", 8); err != nil {
		return nil, err
	}
	if out.TopCities, err = r.breakdown(ctx, start, "city", "event_name='page_view' AND city <> ''", 10); err != nil {
		return nil, err
	}
	if out.Browsers, err = r.breakdown(ctx, start, "browser", "event_name='page_view' AND browser <> ''", 8); err != nil {
		return nil, err
	}
	if out.OperatingSystems, err = r.breakdown(ctx, start, "os", "event_name='page_view' AND os <> ''", 8); err != nil {
		return nil, err
	}
	if out.DeviceTypes, err = r.breakdown(ctx, start, "device_type", "event_name='page_view' AND device_type <> ''", 8); err != nil {
		return nil, err
	}
	if out.TopPages, err = r.breakdown(ctx, start, "path", "event_name='page_view' AND path <> ''", 12); err != nil {
		return nil, err
	}
	if out.TopEntryPages, err = r.entryPages(ctx, start, 10); err != nil {
		return nil, err
	}
	if out.TopExitPages, err = r.breakdown(ctx, start, "path", "event_name='page_exit' AND path <> ''", 10); err != nil {
		return nil, err
	}
	if out.AcquisitionChannels, err = r.breakdown(ctx, start, "source", "event_name='page_view' AND source <> ''", 8); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *AnalyticsRepo) timeline(ctx context.Context, start time.Time, days int) ([]models.AnalyticsTimePoint, error) {
	const q = `SELECT day::date,
		COUNT(p.id) FILTER (WHERE p.event_name='page_view') AS views,
		COUNT(DISTINCT p.visitor_hash) FILTER (WHERE p.event_name='page_view') AS uniques,
		COUNT(p.id) FILTER (WHERE p.event_name='page_exit' AND p.is_bounce) AS bounces
		FROM generate_series($1::date, CURRENT_DATE, INTERVAL '1 day') day
		LEFT JOIN page_views_raw p ON p.created_at::date = day::date
		GROUP BY day ORDER BY day`
	rows, err := r.db.Query(ctx, q, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]models.AnalyticsTimePoint, 0, days)
	for rows.Next() {
		var day time.Time
		var views, uniques, bounces int64
		if err := rows.Scan(&day, &views, &uniques, &bounces); err != nil {
			return nil, err
		}
		point := models.AnalyticsTimePoint{Date: day.Format("2006-01-02"), Label: day.Format("Jan 2"), Views: views, UniqueVisitors: uniques}
		if views > 0 {
			point.BounceRate = roundPercent(float64(bounces) / float64(views) * 100)
		}
		out = append(out, point)
	}
	return out, rows.Err()
}

func (r *AnalyticsRepo) breakdown(ctx context.Context, start time.Time, column, filter string, limit int) ([]models.AnalyticsMetricRow, error) {
	q := `WITH rows AS (
		SELECT ` + column + ` AS label, COUNT(*)::bigint AS value
		FROM page_views_raw WHERE created_at >= $1 AND ` + filter + `
		GROUP BY ` + column + `
	), total AS (SELECT COALESCE(SUM(value), 0)::float AS total FROM rows)
	SELECT label, value, CASE WHEN total.total > 0 THEN ROUND((value::numeric / total.total::numeric) * 100, 2) ELSE 0 END
	FROM rows, total ORDER BY value DESC, label ASC LIMIT $2`
	rows, err := r.db.Query(ctx, q, start, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.AnalyticsMetricRow{}
	for rows.Next() {
		var item models.AnalyticsMetricRow
		if err := rows.Scan(&item.Label, &item.Value, &item.Share); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *AnalyticsRepo) entryPages(ctx context.Context, start time.Time, limit int) ([]models.AnalyticsMetricRow, error) {
	const q = `WITH first_views AS (
		SELECT DISTINCT ON (session_id) session_id, path
		FROM page_views_raw
		WHERE created_at >= $1 AND event_name='page_view'
		ORDER BY session_id, created_at ASC
	), rows AS (
		SELECT path AS label, COUNT(*)::bigint AS value FROM first_views GROUP BY path
	), total AS (SELECT COALESCE(SUM(value), 0)::float AS total FROM rows)
	SELECT label, value, CASE WHEN total.total > 0 THEN ROUND((value::numeric / total.total::numeric) * 100, 2) ELSE 0 END
	FROM rows, total ORDER BY value DESC, label ASC LIMIT $2`
	rows, err := r.db.Query(ctx, q, start, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.AnalyticsMetricRow{}
	for rows.Next() {
		var item models.AnalyticsMetricRow
		if err := rows.Scan(&item.Label, &item.Value, &item.Share); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func roundPercent(value float64) float64 {
	return math.Round(value*100) / 100
}
