package repositories

import (
	"database/sql"
	"fmt"
	"frontdesk/models"
)

type QueriesRepository interface {
	SaveQuery(query *models.Query) error
	FetchQueries() ([]models.Query, error)
	UpdateQueryStatus(queryStatus *models.QueryStatus, id int) error
	FetchFAQs() ([]models.FAQ, error)
}

type queriesRepository struct {
	db *sql.DB
}

func NewQueriesRepository(db *sql.DB) *queriesRepository {
	return &queriesRepository{db: db}
}

func (r *queriesRepository) SaveQuery(query *models.Query) error {
	_, err := r.db.Exec("SELECT insert_customer_query($1::TEXT, $2::TEXT, $3::TEXT, $4::SMALLINT)",
		query.CustomerId, query.QueryText, query.Answer, query.QueryStatus)

	if err != nil {
		return fmt.Errorf("error executing function: %w", err)
	}

	return nil
}

func (r *queriesRepository) FetchQueries() ([]models.Query, error) {
	rows, err := r.db.Query(`SELECT id, customer_id, created_at, query, answer, query_status FROM "CustomerQueries"`)
	if err != nil {
		return nil, fmt.Errorf("error querying customer queries: %w", err)
	}
	defer rows.Close()

	var queries []models.Query

	for rows.Next() {
		var query models.Query
		if err := rows.Scan(&query.Id, &query.CustomerId, &query.CreatedAt, &query.QueryText, &query.Answer, &query.QueryStatus); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		queries = append(queries, query)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return queries, nil
}

func (r *queriesRepository) UpdateQueryStatus(queryStatus *models.QueryStatus, id int) error {
	_, err := r.db.Exec(`UPDATE "CustomerQueries" SET answer = $1, query_status = $2 WHERE id = $3`,
		queryStatus.Answer, queryStatus.QueryStatus, id)

	if err != nil {
		return fmt.Errorf("error updating query status: %w", err)
	}

	return nil
}

func (r *queriesRepository) FetchFAQs() ([]models.FAQ, error) {
	rows, err := r.db.Query(`SELECT query, answer FROM "CustomerQueries" WHERE query_status = $1`, models.RESOLVED)
	if err != nil {
		return nil, fmt.Errorf("error fetching FAQs: %w", err)
	}
	defer rows.Close()

	var faqs []models.FAQ

	for rows.Next() {
		var faq models.FAQ
		if err := rows.Scan(&faq.Question, &faq.Answer); err != nil {
			return nil, fmt.Errorf("error scanning FAQ row: %w", err)
		}
		faqs = append(faqs, faq)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return faqs, nil
}
