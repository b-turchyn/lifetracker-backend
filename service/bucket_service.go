package service

import (
	"context"
	"database/sql"
	"time"
)

type BucketResultRow struct {
	Key   string     `json:"key"`
	Over  BucketSide `json:"over"`
	Under BucketSide `json:"under"`
}

func NewBucketResultRow() BucketResultRow {
	return BucketResultRow{Over: BucketSide{}, Under: BucketSide{}}
}

type BucketSide struct {
	Average float64 `json:"average"`
	Count   int     `json:"count"`
}

type BucketService struct {
	db *sql.DB
}

func NewBucketService(db *sql.DB) *BucketService {
	return &BucketService{db: db}
}

func (b *BucketService) IndexBucket() ([]string, error) {
	query := ` SELECT distinct key FROM raw_data WHERE type IN ('range', 'boolean', 'number') ORDER BY key `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := b.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)

	for rows.Next() {
		var row string
		rows.Scan(&row)

		result = append(result, row)
	}

	return result, nil
}

func (b *BucketService) ShowBucket(compareTo string, pivotPoint float64) ([]BucketResultRow, error) {
	query := `
    WITH tmpData(bucket, other_key, avg, count) AS (select (rd.value::numeric >= $1) AS bucket,
      nrd.key as other_key,
      avg(nrd.value::numeric) as avg_value,
      count(nrd.id) as count
    FROM raw_data rd
    INNER JOIN raw_data nrd
      ON (nrd.type != 'text')
      AND (abs(rd.timestamp - nrd.timestamp) < (1000 * 60 * 60 * 18))
    WHERE rd.key = $2
    GROUP BY bucket, other_key)
    SELECT o.other_key, o.avg, o.count, u.avg, u.count FROM tmpData o
    LEFT OUTER JOIN tmpData u ON u.bucket = FALSE AND u.other_key = o.other_key
    WHERE o.bucket = TRUE
    ORDER BY o.other_key
`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := b.db.QueryContext(ctx, query, pivotPoint, compareTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]BucketResultRow, 0)

	for rows.Next() {
		row := NewBucketResultRow()
		rows.Scan(&row.Key, &row.Over.Average, &row.Over.Count, &row.Under.Average, &row.Under.Count)

		result = append(result, row)
	}

	return result, nil
}
