# NAMIS Analytics & Reports Plan

The Analytics and Reporting Module builds charts and aggregates data on athlete demographics, medals won, talent identification progress, and geographical representation for the National Council of Sports.

## 1. Database View & Aggregations

We implement dashboard queries in the repository to compile high-performing reports.

### Key Queries
1. **Athletes Demographics Report**:
   ```sql
   SELECT gender, COUNT(*) FROM athletes GROUP BY gender;
   SELECT region, COUNT(*) FROM athletes GROUP BY region;
   SELECT age_category, COUNT(*) FROM athletes GROUP BY age_category;
   ```
2. **Medal Tables by Federation**:
   ```sql
   SELECT f.name, 
          COUNT(CASE WHEN m.medal_type = 'GOLD' THEN 1 END) as gold_count,
          COUNT(CASE WHEN m.medal_type = 'SILVER' THEN 1 END) as silver_count,
          COUNT(CASE WHEN m.medal_type = 'BRONZE' THEN 1 END) as bronze_count
   FROM medals m
   JOIN federations f ON f.id = m.federation_id
   WHERE m.status = 'APPROVED'
   GROUP BY f.name ORDER BY gold_count DESC, silver_count DESC, bronze_count DESC;
   ```
3. **Talent Pathway Statistics**:
   ```sql
   SELECT status, COUNT(*) FROM talent_records GROUP BY status;
   ```

## 2. Generic Backend Mapping

We extend the repository endpoint `h.NSMIS.AthleteDashboard` inside `backend/internal/repository/nsmis.go` to return these dynamic aggregate datasets in JSON.

## 3. UI Plan

- **Athlete Analytics Panel**: Dashboard page containing:
  - Pie charts of gender distribution and bar graphs of age groups.
  - Interactive maps of athlete regions/districts of origin.
- **NCS KPI Metrics Board**: High-level indicator widgets displaying:
  - Total Registered Athletes in Uganda
  - Active Elite Athletes
  - National Team Representatives
  - International Medals won
- **Medal Standings Table**: Searchable and sortable rankings showing medals by federation.
