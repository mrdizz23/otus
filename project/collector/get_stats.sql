SELECT format(
   'INSERT INTO raw_data (datname, username, queryid, query, calls, total_exec_time, min_exec_time, max_exec_time, mean_exec_time, stddev_exec_time, rows, shared_blks_hit, shared_blks_read, temp_blks_read, temp_blks_written, blk_read_time, blk_write_time) VALUES (''%s'', ''%s'', %s, ''%s'', %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);',
        pg_database.datname,
        pg_get_userbyid(userid),
        pg_stat_statements.queryid,
        replace(pg_stat_statements.query,'''','"'),
        pg_stat_statements.calls,
        pg_stat_statements.total_time,
        pg_stat_statements.min_time,
        pg_stat_statements.max_time,
        pg_stat_statements.mean_time,
        pg_stat_statements.rows,
        pg_stat_statements.shared_blks_hit,
        pg_stat_statements.shared_blks_read,
        pg_stat_statements.temp_blks_read,
        pg_stat_statements.temp_blks_written,
        pg_stat_statements.blk_read_time,
        pg_stat_statements.blk_write_time
)
FROM pg_stat_statements
JOIN pg_database ON pg_database.oid = pg_stat_statements.dbid
WHERE pg_database.datname != 'postgres' and pg_database.datname !~* '__';
