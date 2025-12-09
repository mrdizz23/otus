package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var intervals = []string{"hour", "day", "week"}

var columns = []string{
	"datname",
	"username",
	"queryid",
	"query",
	"calls",
	"total_exec_time",
	"mean_exec_time",
	"max_exec_time",
	"rows",
	"shared_blks_hit",
	"shared_blks_read",
	"temp_blks_read",
	"temp_blks_written",
	"blk_read_time",
	"blk_write_time",
}

type data struct {
	datname           string
	username          string
	queryid           int64
	query             string
	calls             int64
	total_exec_time   float32
	mean_exec_time    float32
	max_exec_time     float32
	rows              int64
	shared_blks_hit   int64
	shared_blks_read  int64
	temp_blks_read    int64
	temp_blks_written int64
	blk_read_time     int64
	blk_write_time    int64
}

type Dictionary map[string]interface{}

func main() {
	connString := "postgres://statements_owner@statements:5432/statements?sslmode=disable"

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	delete_query := "DELETE FROM raw_data where date < now() - interval '1 week'"
	_, err = db.Exec(delete_query)
	if err != nil {
		panic(err)
	}

	for _, interval := range intervals {

		// trancate table with list of databases

		trunc_query := "truncate table databases_" + interval
		_, err := db.Exec(trunc_query)
		if err != nil {
			panic(err)
		}

		select_query := fmt.Sprintf(
			"select %s from raw_data where and date >= now() - interval '1 %s'",
			strings.Join(columns, ", "),
			interval)

		rows, err := db.Query(select_query)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		stats := make(map[string]map[int64]data)

		for rows.Next() {
			p := data{}
			err := rows.Scan(&p.datname, &p.username, &p.queryid, &p.query, &p.calls, &p.total_exec_time, &p.mean_exec_time, &p.max_exec_time, &p.rows, &p.shared_blks_hit, &p.shared_blks_read, &p.temp_blks_read, &p.temp_blks_written, &p.blk_read_time, &p.blk_write_time)
			if err != nil {
				fmt.Println(err)
				continue
			}

			_, ok := stats[p.datname]
			if !ok {
				stats[p.datname] = map[int64]data{}
			}

			entry, ok := stats[p.datname][p.queryid]

			// if interval == "hour" || !ok {

			if !ok {
				stats[p.datname][p.queryid] = data{
					datname:           p.datname,
					queryid:           p.queryid,
					username:          p.username,
					query:             p.query,
					calls:             p.calls,
					total_exec_time:   p.total_exec_time,
					max_exec_time:     p.max_exec_time,
					mean_exec_time:    p.mean_exec_time,
					rows:              p.rows,
					shared_blks_hit:   p.shared_blks_hit,
					shared_blks_read:  p.shared_blks_read,
					temp_blks_read:    p.temp_blks_read,
					temp_blks_written: p.temp_blks_written,
					blk_read_time:     p.blk_read_time,
					blk_write_time:    p.blk_write_time,
				}
			} else {
				entry.calls += p.calls
				entry.total_exec_time += p.total_exec_time
				entry.max_exec_time = max(entry.max_exec_time, p.max_exec_time)
				entry.mean_exec_time = (entry.mean_exec_time + p.mean_exec_time) / 2
				entry.rows += p.rows
				entry.shared_blks_hit += p.shared_blks_hit
				entry.shared_blks_read += p.shared_blks_read
				entry.temp_blks_read += p.temp_blks_read
				entry.temp_blks_written += p.temp_blks_written
				entry.blk_read_time += p.blk_read_time
				entry.blk_write_time += p.blk_write_time
				stats[p.datname][p.queryid] = entry
			}

		}

		// create or truncate target table

		create_table_query := "create table IF NOT EXISTS stats_" + interval + `(
				id bigserial PRIMARY KEY,
				datname text,
				queryid bigint,
				username text,
				query text,
				calls bigint,
				total_exec_time double precision,
				max_exec_time double precision,
				mean_exec_time double precision,
				rows bigint,
				shared_blks_hit bigint,
				shared_blks_read bigint,
				temp_blks_read bigint,
				temp_blks_written bigint,
				blk_read_time double precision,
				blk_write_time double precision
			);`
		_, err = db.Exec(create_table_query)
		if err != nil {
			panic(err)
		}

		trunc_query = "truncate table stats_" + interval
		_, err = db.Exec(trunc_query)
		if err != nil {
			panic(err)
		}

		// fmt.Println(len(stats))

		datname_keys := make([]string, len(stats))
		i := 0
		for k := range stats {
			datname_keys[i] = k
			i++
		}

		// insert databases to database interval table
		datname_table := "databases_" + interval
		// fmt.Println(datname_table)

		txn, err := db.Begin()
		if err != nil {
			panic(err)
		}

		stmt, err := txn.Prepare(pq.CopyIn(datname_table, "datname"))
		if err != nil {
			panic(err)
		}

		for _, datname := range datname_keys {

			_, err = stmt.Exec(datname)
			if err != nil {
				panic(err)
			}
		}
		_, err = stmt.Exec()
		if err != nil {
			panic(err)
		}

		err = stmt.Close()
		if err != nil {
			panic(err)
		}

		err = txn.Commit()
		if err != nil {
			panic(err)
		}

		// insert stats to stats interval table
		stats_table := "stats_" + interval
		// fmt.Println(stats_table)

		for _, database := range datname_keys {

			// fmt.Println(database)
			// fmt.Println(len(stats[database]))

			query_ids := make([]int64, len(stats[database]))
			i := 0
			for j := range stats[database] {
				query_ids[i] = j
				i++
			}

			txn, err = db.Begin()
			if err != nil {
				panic(err)
			}

			stmt, err = txn.Prepare(pq.CopyIn(stats_table, "datname", "queryid", "username", "query", "calls", "total_exec_time", "max_exec_time", "mean_exec_time", "rows", "shared_blks_hit", "shared_blks_read", "temp_blks_read", "temp_blks_written", "blk_read_time", "blk_write_time"))
			if err != nil {
				panic(err)
			}

			for _, query_id := range query_ids {

				// fmt.Println(query_id)

				_, err = stmt.Exec(
					stats[database][query_id].datname,
					stats[database][query_id].queryid,
					stats[database][query_id].username,
					stats[database][query_id].query,
					stats[database][query_id].calls,
					stats[database][query_id].total_exec_time,
					stats[database][query_id].max_exec_time,
					stats[database][query_id].mean_exec_time,
					stats[database][query_id].rows,
					stats[database][query_id].shared_blks_hit,
					stats[database][query_id].shared_blks_read,
					stats[database][query_id].temp_blks_read,
					stats[database][query_id].temp_blks_written,
					stats[database][query_id].blk_read_time,
					stats[database][query_id].blk_write_time)
				if err != nil {
					panic(err)
				}
			}

			_, err = stmt.Exec()
			if err != nil {
				panic(err)
			}

			err = stmt.Close()
			if err != nil {
				panic(err)
			}

			err = txn.Commit()
			if err != nil {
				panic(err)
			}

		}

	}

}
