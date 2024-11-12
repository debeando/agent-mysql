package main

import (
	"fmt"
	"strings"
	"time"
)

type Query struct {
	Counter   int64
	Interval  int
	Key       string
	Name      string
	Statement string
	UnPivot   bool
	Value     string
}

var Queries = []*Query{
	&Query{
		Name:      "mysql_variables",
		Statement: "SHOW GLOBAL VARIABLES",
		Key:       "Variable_name",
		Value:     "Value",
	},
	&Query{
		Name:      "mysql_status",
		Statement: "SHOW GLOBAL STATUS",
		Key:       "Variable_name",
		Value:     "Value",
	},
	&Query{
		Name:      "mysql_replica",
		Statement: "SHOW REPLICA STATUS",
		UnPivot:   true,
	},
	&Query{
		Statement: "SELECT name, count FROM information_schema.innodb_metrics WHERE status='enabled'",
		Name:      "mysql_innodb",
		Key:       "name",
		Value:     "count",
	},
	&Query{
		Name: "mysql_performance_schema",
		Statement: fmt.Sprintf(`
        SELECT
            ifnull(SCHEMA_NAME, 'NONE') AS SCHEMA_NAME,
            DIGEST,
            DIGEST_TEXT,
            COUNT_STAR,
            SUM_TIMER_WAIT,
            SUM_ERRORS,
            SUM_WARNINGS,
            SUM_ROWS_AFFECTED,
            SUM_ROWS_SENT,
            SUM_ROWS_EXAMINED,
            SUM_CREATED_TMP_DISK_TABLES,
            SUM_CREATED_TMP_TABLES,
            SUM_SORT_MERGE_PASSES,
            SUM_SORT_ROWS,
            SUM_NO_INDEX_USED
        FROM performance_schema.events_statements_summary_by_digest
        WHERE SCHEMA_NAME NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
          AND last_seen > DATE_SUB(NOW(), INTERVAL %d SECOND)
        ORDER BY SUM_TIMER_WAIT DESC;
        `, int(getInterval().Seconds())),
		UnPivot: true,
	},
	&Query{
		Name: "mysql_query_latency",
		Statement: fmt.Sprintf(`
        SELECT
            ifnull(SCHEMA_NAME, 'NONE') AS SCHEMA_NAME,
            sum(count_star) AS count_star,
            round(avg_timer_wait/1000000, 0) AS avg_time_us
        FROM performance_schema.events_statements_summary_by_digest
        WHERE SCHEMA_NAME NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
          AND last_seen > DATE_SUB(NOW(), INTERVAL %d SECOND)
        GROUP BY SCHEMA_NAME
        ORDER BY avg_time_us DESC;
        `, int(getInterval().Seconds())),
		UnPivot: true,
	},
	&Query{
		Interval: 3600,
		Name:     "mysql_overflow",
		Statement: `
        SELECT
            t.table_schema AS SCHEMA_NAME,
            t.table_name,
            t.table_rows,
            t.auto_increment,
            (SELECT column_type FROM information_schema.columns WHERE table_schema = t.table_schema AND table_name = t.table_name AND extra = 'auto_increment' LIMIT 1) AS auto_increment_data_type,
            (CASE 
               WHEN (SELECT column_type FROM information_schema.columns WHERE table_schema = t.table_schema AND table_name = t.table_name AND extra = 'auto_increment' LIMIT 1) IN ('int unsigned', "int(10) unsigned") THEN ROUND( (t.auto_increment/4294967295)*100 , 2)
               WHEN (SELECT column_type FROM information_schema.columns WHERE table_schema = t.table_schema AND table_name = t.table_name AND extra = 'auto_increment' LIMIT 1) IN ('int(11)', 'int') THEN ROUND( (t.auto_increment/2147483647)*100, 2)
               WHEN (SELECT column_type FROM information_schema.columns WHERE table_schema = t.table_schema AND table_name = t.table_name AND extra = 'auto_increment' LIMIT 1) IN ('bigint unsigned', 'bigint(20) unsigned') THEN ROUND( (t.auto_increment/(POWER(2, 64) -1))*100 , 2 )
               WHEN (SELECT column_type FROM information_schema.columns WHERE table_schema = t.table_schema AND table_name = t.table_name AND extra = 'auto_increment' LIMIT 1) IN ('bigint(20)', 'bigint' ) THEN ROUND( (t.auto_increment/(POWER(2, 64) -1))*100 , 2 )
            END) AS auto_increment_pct
        FROM information_schema.tables t
        WHERE t.table_schema NOT IN ('sys', 'mysql', 'performance_schema', 'information_schema')
          AND t.auto_increment IS NOT NULL;
        `,
		UnPivot: true,
	},
	&Query{
		Interval: 3600,
		Name:     "mysql_tables",
		Statement: `
        SELECT
            table_schema AS SCHEMA_NAME,
            table_name,
            COALESCE(data_length + index_length, 0) AS 'table_size',
            COALESCE(table_rows, 0) AS 'table_rows'
        FROM information_schema.tables
        WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema')
        ORDER BY table_schema, table_name;
        `,
		UnPivot: true,
	},
}

func (q *Query) Beautifier() string {
	q.Statement = strings.ReplaceAll(q.Statement, "\r\n", " ")
	q.Statement = strings.ReplaceAll(q.Statement, "\n", " ")
	q.Statement = strings.ReplaceAll(q.Statement, "\t", " ")
	q.Statement = strings.ReplaceAll(q.Statement, "  ", "")
	q.Statement = strings.Trim(q.Statement, " ")

	return q.Statement
}

func (q *Query) IsTime(i int) bool {
	if q.Interval == 0 {
		return true
	}

	if q.Counter == 0 || int(time.Since(time.Unix(q.Counter, 0)).Seconds()) >= i {
		(*q).Counter = int64(time.Now().Unix())

		return true
	}

	return false
}
