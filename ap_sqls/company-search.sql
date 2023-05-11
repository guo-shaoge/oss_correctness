set @@tidb_isolation_read_engines="tiflash";
set @@collation_connection = 'utf8_bin';
WITH company_ci AS (
    SELECT gu.organization AS company, COUNT(1) AS cnt
    FROM github_users gu
    WHERE
        gu.organization LIKE CONCAT('%', 'PingCAP', '%')
    GROUP BY company
)
SELECT
    TRIM(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(company, 'www.', ''), '.com', ''), '!', ''), ',', ''), '-', ''), '@', ''), '.', ''), 'ltd', ''), 'inc', ''), 'corporation', '')) AS `name`,
    SUM(cnt) AS total
FROM company_ci
GROUP BY `name`
ORDER BY total DESC, name
LIMIT 10
