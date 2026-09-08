package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/atenimedia-llc/ncs-online/backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestExpenseValidation(t *testing.T) {
	for _, v := range []string{"0", "-1", "1.001", "NaN", "1e3", "1000000000000", ""} {
		if validExpenseAmount(v) {
			t.Errorf("accepted %q", v)
		}
	}
	for _, v := range []string{"0.01", "100", "999999999999.99"} {
		if !validExpenseAmount(v) {
			t.Errorf("rejected %q", v)
		}
	}
	if validExpenseDate("2026-02-30") || !validExpenseDate("2026-09-08") {
		t.Fatal("date validation")
	}
}

func TestExpensesIntegration(t *testing.T) {
	url := os.Getenv("NCS_EXPENSE_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set NCS_EXPENSE_TEST_DATABASE_URL; uses and removes an isolated schema")
	}
	ctx := context.Background()
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	schema := fmt.Sprintf("expense_test_%d", time.Now().UnixNano())
	if _, err = db.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer db.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	_, err = pool.Exec(ctx, `CREATE TABLE users(id TEXT PRIMARY KEY,first_name TEXT,last_name TEXT);CREATE TABLE departments(id TEXT PRIMARY KEY,name TEXT,is_active BOOLEAN);INSERT INTO users VALUES('stella','Stella','Akello'),('admin','Test','Administrator');INSERT INTO departments VALUES('finance','Finance & Accounts',true);`)
	if err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../migrations/079_create_expenses.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, string(migration)); err != nil {
		t.Fatal(err)
	}
	router := NewExpenseHandler(pool).Routes()
	request := func(method, path, role, contentType string, body io.Reader, status int) *httptest.ResponseRecorder {
		t.Helper()
		r := httptest.NewRequest(method, path, body)
		uid := "stella"
		if role == "admin" {
			uid = "admin"
		}
		r = r.WithContext(context.WithValue(context.WithValue(r.Context(), models.CtxUserID, uid), models.CtxUserRoles, []string{role}))
		r.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		if w.Code != status {
			t.Fatalf("%s %s as %s: %d expected %d: %s", method, path, role, w.Code, status, w.Body.String())
		}
		return w
	}
	data := func(w *httptest.ResponseRecorder) map[string]interface{} {
		t.Helper()
		var e struct {
			Data map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
			t.Fatal(err)
		}
		return e.Data
	}
	catBody := `{"name":"Office supplies","is_active":true}`
	request("POST", "/categories", "accountant", "application/json", strings.NewReader(catBody), 403)
	cat := data(request("POST", "/categories", "admin", "application/json", strings.NewReader(catBody), 201))["id"].(string)
	request("POST", "/categories", "admin", "application/json", strings.NewReader(catBody), 409)
	for _, role := range []string{"accountant", "senior_accountant", "chief_accountant", "asset_accountant", "finance_department", "admin", "super_admin"} {
		request("GET", "/", ""+role, "", nil, 200)
	}
	request("GET", "/", "staff", "", nil, 403)
	if data(request("GET", "/options", "accountant", "", nil, 200))["recorded_by_name"] != "Stella Akello" {
		t.Fatal("recorder options")
	}
	create := func(overrides map[string]string, receipt []byte, status int) *httptest.ResponseRecorder {
		t.Helper()
		fields := map[string]string{"category_id": cat, "department_id": "finance", "title": "Stationery purchase", "description": "Pens for office", "amount": "12500.50", "payee": "Office shop", "payment_method": "CASH", "recorded_by": "admin", "recorded_by_name": "Forged name"}
		for k, v := range overrides {
			fields[k] = v
		}
		var b bytes.Buffer
		mw := multipart.NewWriter(&b)
		for k, v := range fields {
			mw.WriteField(k, v)
		}
		if receipt != nil {
			part, _ := mw.CreateFormFile("file", "receipt.pdf")
			part.Write(receipt)
		}
		mw.Close()
		return request("POST", "/", "accountant", mw.FormDataContentType(), &b, status)
	}
	pdf := []byte("%PDF-1.4\nTest expense receipt\n%%EOF")
	id := data(create(nil, pdf, 201))["id"].(string)
	record := data(request("GET", "/"+id, "accountant", "", nil, 200))
	if record["recorded_by"] != "stella" || record["recorded_by_name"] != "Stella Akello" || record["department_name"] != "Finance & Accounts" || record["expense_date"] != expenseToday() || record["amount"] != 12500.50 {
		t.Fatalf("wrong saved record: %v", record)
	}
	download := request("GET", "/"+id+"/attachment", "accountant", "", nil, 200)
	if !bytes.Equal(download.Body.Bytes(), pdf) || download.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("receipt mismatch")
	}
	request("GET", "/"+id+"/attachment", "staff", "", nil, 403)
	for _, fields := range []map[string]string{{"amount": "0"}, {"expense_date": "2999-01-01"}, {"department_id": "missing"}, {"category_id": "missing"}} {
		create(fields, nil, 400)
	}
	create(nil, []byte("<script>alert(1)</script>"), 415)
	create(nil, bytes.Repeat([]byte("x"), (5<<20)+1), 400)
	list := data(request("GET", "/?search=Stella&department_id=finance", "accountant", "", nil, 200))
	if list["total"] != float64(1) || list["total_amount"] != "12500.50" {
		t.Fatalf("list totals: %v", list)
	}
	empty := data(request("GET", "/?offset=50", "accountant", "", nil, 200))
	if len(empty["expenses"].([]interface{})) != 0 {
		t.Fatal("pagination")
	}
	request("GET", "/?from=invalid", "accountant", "", nil, 400)
	withoutReceipt := data(create(nil, nil, 201))["id"].(string)
	request("GET", "/"+withoutReceipt+"/attachment", "accountant", "", nil, 404)
	request("PUT", "/categories/"+cat, "accountant", "application/json", strings.NewReader(catBody), 403)
	request("PUT", "/categories/"+cat, "admin", "application/json", strings.NewReader(`{"name":"Renamed category","is_active":false}`), 200)
	create(nil, nil, 400)
	record = data(request("GET", "/"+id, "accountant", "", nil, 200))
	if record["category_name"] != "Office supplies" {
		t.Fatal("history changed")
	}
	cats := request("GET", "/categories", "accountant", "", nil, 200)
	if strings.Contains(cats.Body.String(), "Renamed") {
		t.Fatal("inactive category exposed for entry")
	}
	var count int
	pool.QueryRow(ctx, "SELECT count(*) FROM expenses").Scan(&count)
	if count != 2 {
		t.Fatalf("failed submissions saved %d rows", count)
	}
}
