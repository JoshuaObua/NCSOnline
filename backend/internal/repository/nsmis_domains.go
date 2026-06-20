package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type domainDef struct {
	table, federationExpr string
	required              []string
	columns               map[string]string
}

var domainDefs = map[string]domainDef{
	"federation-officers":     {"federation_officers", "t.federation_id=ANY($3)", []string{"federation_id", "position", "full_name"}, map[string]string{"federation_id": "text", "position": "text", "position_label": "text", "full_name": "text", "email": "text", "phone": "text", "appointed_on": "date?", "term_ends_on": "date?", "is_active": "boolean"}},
	"federation-memberships":  {"federation_memberships", "t.federation_id=ANY($3)", []string{"federation_id", "user_id", "membership_role"}, map[string]string{"federation_id": "text", "user_id": "text", "membership_role": "text", "starts_at": "date", "ends_at": "date?", "is_active": "boolean"}},
	"athletes":                {"athletes", "EXISTS(SELECT 1 FROM athlete_affiliations aa WHERE aa.athlete_id=t.id AND aa.federation_id=ANY($3))", []string{"federation_id", "full_name", "gender", "date_of_birth", "discipline"}, map[string]string{"athlete_number": "text", "full_name": "text", "gender": "text", "date_of_birth": "date", "district": "text", "region": "text", "club": "text", "discipline": "text", "national_team_status": "text", "status": "text", "consent_basis": "text"}},
	"competitions":            {"competitions", "t.federation_id=ANY($3)", []string{"federation_id", "name", "venue", "level", "starts_on", "ends_on"}, map[string]string{"federation_id": "text", "name": "text", "venue": "text", "host_country": "text", "level": "text", "starts_on": "date", "ends_on": "date", "returned_on": "date?", "athlete_count": "integer", "official_count": "integer", "status": "text"}},
	"medals":                  {"medals", "t.federation_id=ANY($3)", []string{"competition_id", "federation_id", "athlete_name", "event", "country", "won_on", "medal_type"}, map[string]string{"competition_id": "text", "federation_id": "text", "athlete_id": "text?", "athlete_name": "text", "event": "text", "country": "text", "won_on": "date", "medal_type": "text", "coach_responsible": "text", "team_manager": "text", "funding_source": "text", "status": "text"}},
	"coaches":                 {"coaches", "t.federation_id=ANY($3)", []string{"federation_id", "full_name", "certification_level", "license_number"}, map[string]string{"federation_id": "text", "full_name": "text", "certification_level": "text", "license_number": "text", "expiry_date": "date?", "status": "text"}},
	"technical-officials":     {"technical_officials", "t.federation_id=ANY($3)", []string{"federation_id", "full_name", "official_type", "level", "certification"}, map[string]string{"federation_id": "text", "full_name": "text", "official_type": "text", "level": "text", "certification": "text", "valid_until": "date?", "status": "text"}},
	"talent":                  {"talent_records", "t.federation_id=ANY($3)", []string{"federation_id", "athlete_name", "age_at_identification", "district", "identified_by", "identified_on"}, map[string]string{"federation_id": "text", "athlete_id": "text?", "athlete_name": "text", "age_at_identification": "integer", "school": "text", "district": "text", "region": "text", "identified_by": "text", "identified_on": "date", "status": "text"}},
	"scholarships":            {"talent_scholarships", "EXISTS(SELECT 1 FROM talent_records tr WHERE tr.id=t.talent_record_id AND tr.federation_id=ANY($3))", []string{"talent_record_id", "provider", "starts_on"}, map[string]string{"talent_record_id": "text", "provider": "text", "starts_on": "date", "ends_on": "date?", "status": "text"}},
	"safeguarding-aggregates": {"safeguarding_period_metrics", "t.federation_id=ANY($3)", []string{"federation_id", "reporting_period_id"}, map[string]string{"federation_id": "text", "reporting_period_id": "text", "female_athletes": "integer", "female_coaches": "integer", "safeguarding_cases": "integer", "resolved_cases": "integer", "status": "text"}},
	"safeguarding-cases":      {"safeguarding_cases", "t.federation_id=ANY($3)", []string{"case_reference", "federation_id", "encrypted_details", "reported_at"}, map[string]string{"case_reference": "text", "federation_id": "text", "encrypted_details": "text", "status": "text", "reported_at": "timestamptz", "resolved_at": "timestamptz?", "assigned_to": "text?"}},
	"disbursements":           {"ncs_disbursements", "t.federation_id=ANY($3)", []string{"federation_id", "reference", "amount", "released_on", "accountability_due_on"}, map[string]string{"federation_id": "text", "reference": "text", "amount": "numeric", "currency": "text", "released_on": "date", "accountability_due_on": "date", "status": "text"}},
	"accountabilities":        {"financial_reports", "t.federation_id=ANY($3)", []string{"federation_id", "reporting_period_id"}, map[string]string{"federation_id": "text", "reporting_period_id": "text", "currency": "text", "government_grant": "numeric", "sponsorship": "numeric", "membership_fees": "numeric", "donations": "numeric", "competitions": "numeric", "training": "numeric", "equipment": "numeric", "administration": "numeric", "status": "text"}},
	"equipment":               {"equipment_items", "t.federation_id=ANY($3)", []string{"federation_id", "item_name"}, map[string]string{"federation_id": "text", "item_name": "text", "unit": "text", "quantity_received": "integer", "quantity_distributed": "integer"}},
}

func DomainNames() []string {
	out := make([]string, 0, len(domainDefs))
	for k := range domainDefs {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func domainDefinition(name string) (domainDef, bool) { d, ok := domainDefs[name]; return d, ok }

func (r *NSMISRepo) FederationScope(ctx context.Context, userID string, all bool) ([]string, error) {
	if all {
		return []string{}, nil
	}
	rows, err := r.db.Query(ctx, `SELECT federation_id FROM federation_memberships WHERE user_id=$1 AND is_active AND (ends_at IS NULL OR ends_at>=CURRENT_DATE)`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *NSMISRepo) ListDomain(ctx context.Context, resource, userID, search string, all bool, limit, offset int) ([]map[string]interface{}, int64, error) {
	d, ok := domainDefinition(resource)
	if !ok {
		return nil, 0, ErrNotFound
	}
	ids, err := r.FederationScope(ctx, userID, all)
	if err != nil {
		return nil, 0, err
	}
	scope := "TRUE"
	if !all {
		scope = d.federationExpr
	}
	where := fmt.Sprintf("(%s) AND ($1='' OR to_jsonb(t)::text ILIKE '%%'||$1||'%%')", scope)
	var total int64
	if err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM "+d.table+" t WHERE "+where, search, limit, ids).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, "SELECT to_jsonb(t) FROM "+d.table+" t WHERE "+where+" ORDER BY t.created_at DESC LIMIT $2 OFFSET $4", search, limit, ids, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []map[string]interface{}{}
	for rows.Next() {
		var raw []byte
		if err = rows.Scan(&raw); err != nil {
			return nil, 0, err
		}
		var item map[string]interface{}
		if err = json.Unmarshal(raw, &item); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func jsonExpr(field, typ string) string {
	base := fmt.Sprintf("$1::jsonb->>'%s'", field)
	if strings.HasSuffix(typ, "?") {
		return fmt.Sprintf("NULLIF(%s,'')::%s", base, strings.TrimSuffix(typ, "?"))
	}
	if typ == "text" {
		return base
	}
	return fmt.Sprintf("(%s)::%s", base, typ)
}

func (r *NSMISRepo) SaveDomain(ctx context.Context, resource, id, actorID string, payload map[string]interface{}, allowedFederations []string, all bool) (map[string]interface{}, error) {
	d, ok := domainDefinition(resource)
	if !ok {
		return nil, ErrNotFound
	}
	if id == "" {
		for _, f := range d.required {
			if v, exists := payload[f]; !exists || strings.TrimSpace(fmt.Sprint(v)) == "" {
				return nil, fmt.Errorf("%s is required", f)
			}
		}
	}
	if resource == "athletes" && strings.TrimSpace(fmt.Sprint(payload["athlete_number"])) == "" {
		payload["athlete_number"] = fmt.Sprintf("NCS-%d-%s", time.Now().Year(), strings.ToUpper(uuid.NewString()[:8]))
	}
	if fed, exists := payload["federation_id"]; exists && !all {
		found := false
		for _, x := range allowedFederations {
			if x == fmt.Sprint(fed) {
				found = true
			}
		}
		if !found {
			return nil, ErrNotFound
		}
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	cols := []string{}
	vals := []string{}
	sets := []string{}
	for name, typ := range d.columns {
		if _, exists := payload[name]; !exists {
			continue
		}
		cols = append(cols, name)
		vals = append(vals, jsonExpr(name, typ))
		sets = append(sets, name+"="+jsonExpr(name, typ))
	}
	sort.Strings(cols) // rebuild expressions in deterministic column order
	vals = vals[:0]
	sets = sets[:0]
	for _, name := range cols {
		typ := d.columns[name]
		vals = append(vals, jsonExpr(name, typ))
		sets = append(sets, name+"="+jsonExpr(name, typ))
	}
	var out []byte
	if id == "" {
		extraCols, extraVals := "", ""
		if _, ok := d.columns["created_by"]; ok {
			extraCols = ",created_by"
			extraVals = ",$2"
		} else if d.table == "athletes" || d.table == "competitions" || d.table == "medals" || d.table == "talent_records" || d.table == "ncs_disbursements" || d.table == "financial_reports" || d.table == "safeguarding_cases" {
			extraCols = ",created_by"
			extraVals = ",$2"
		} else if d.table == "federation_memberships" {
			extraCols = ",assigned_by"
			extraVals = ",$2"
		}
		tx, e := r.db.Begin(ctx)
		if e != nil {
			return nil, e
		}
		defer tx.Rollback(ctx)
		q := "INSERT INTO " + d.table + " (" + strings.Join(cols, ",") + extraCols + ") VALUES (" + strings.Join(vals, ",") + extraVals + ") RETURNING to_jsonb(" + d.table + ".*)"
		if e = tx.QueryRow(ctx, q, raw, actorID).Scan(&out); e != nil {
			if isDuplicate(e) {
				return nil, ErrDuplicate
			}
			return nil, e
		}
		if resource == "athletes" {
			var row map[string]interface{}
			_ = json.Unmarshal(out, &row)
			_, e = tx.Exec(ctx, `INSERT INTO athlete_affiliations(athlete_id,federation_id,starts_on) VALUES($1,$2,CURRENT_DATE)`, row["id"], payload["federation_id"])
			if e != nil {
				return nil, e
			}
		}
		if e = tx.Commit(ctx); e != nil {
			return nil, e
		}
	} else {
		if len(sets) == 0 {
			return nil, fmt.Errorf("no editable fields supplied")
		}
		scope := "TRUE"
		if !all {
			scope = d.federationExpr
		}
		q := "UPDATE " + d.table + " t SET " + strings.Join(sets, ",") + ",updated_at=NOW() WHERE t.id=$2 AND (" + scope + ") RETURNING to_jsonb(t.*)"
		e := r.db.QueryRow(ctx, q, raw, id, allowedFederations).Scan(&out)
		if e == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		if e != nil {
			if isDuplicate(e) {
				return nil, ErrDuplicate
			}
			return nil, e
		}
	}
	var item map[string]interface{}
	if err = json.Unmarshal(out, &item); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *NSMISRepo) DeleteDomain(ctx context.Context, resource, id string, scope []string, all bool) error {
	d, ok := domainDefinition(resource)
	if !ok {
		return ErrNotFound
	}
	condition := "TRUE"
	if !all {
		condition = d.federationExpr
	}
	q := "DELETE FROM " + d.table + " t WHERE t.id=$1 AND (" + condition + ")"
	res, err := r.db.Exec(ctx, q, id, nil, scope)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
