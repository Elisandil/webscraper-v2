package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/domain/repository"
	"webscraper-v2/internal/infrastructure/database"
)

type scrapingRepository struct {
	db *database.SQLiteDB
}

func NewScrapingRepository(db *database.SQLiteDB) repository.ScrapingRepository {
	return &scrapingRepository{db: db}
}

func (r *scrapingRepository) Save(result *entity.ScrapingResult) error {
	query := `INSERT INTO scraping_results (
		user_id, url, title, description, keywords, author, language, favicon, 
		image_url, site_name, links, images, headers, status_code, 
		content_type, word_count, load_time_ms, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	linksJSON, err := json.Marshal(result.Links)

	if err != nil {
		return fmt.Errorf("error marshaling links: %w", err)
	}
	imagesJSON, err := json.Marshal(result.Images)
	if err != nil {
		return fmt.Errorf("error marshaling images: %w", err)
	}

	headersJSON, err := json.Marshal(result.Headers)
	if err != nil {
		return fmt.Errorf("error marshaling headers: %w", err)
	}
	res, err := r.db.Exec(query,
		result.UserID, result.URL, result.Title, result.Description,
		result.Keywords, result.Author, result.Language, result.Favicon,
		result.ImageURL, result.SiteName, string(linksJSON), string(imagesJSON),
		string(headersJSON), result.StatusCode, result.ContentType,
		result.WordCount, result.LoadTime, result.CreatedAt)

	if err != nil {
		return fmt.Errorf("error executing insert: %w", err)
	}
	id, err := res.LastInsertId()

	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	result.ID = id
	return nil
}

func (r *scrapingRepository) FindAll() ([]*entity.ScrapingResult, error) {
	query := `SELECT 
		id, url, title, description, keywords, author, language, favicon,
		image_url, site_name, links, images, headers, status_code,
		content_type, word_count, load_time_ms, created_at 
	FROM scraping_results 
	ORDER BY created_at DESC`

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error querying results: %w", err)
	}
	defer rows.Close()
	var results []*entity.ScrapingResult
	for rows.Next() {
		result := &entity.ScrapingResult{}
		var linksJSON, imagesJSON, headersJSON, createdAt string

		err := rows.Scan(
			&result.ID, &result.URL, &result.Title, &result.Description,
			&result.Keywords, &result.Author, &result.Language, &result.Favicon,
			&result.ImageURL, &result.SiteName, &linksJSON, &imagesJSON,
			&headersJSON, &result.StatusCode, &result.ContentType,
			&result.WordCount, &result.LoadTime, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		if err := r.unmarshalJSONField(linksJSON, &result.Links); err != nil {
			result.Links = []string{}
		}
		if err := r.unmarshalJSONField(imagesJSON, &result.Images); err != nil {
			result.Images = []string{}
		}
		if err := r.unmarshalJSONField(headersJSON, &result.Headers); err != nil {
			result.Headers = []entity.Header{}
		}
		result.CreatedAt, err = r.parseDateTime(createdAt)

		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return results, nil
}

func (r *scrapingRepository) FindAllByUserID(userID int64) ([]*entity.ScrapingResult, error) {
	query := `
    SELECT 
        id, user_id, url, title, description, keywords, author, language,
        favicon, image_url, site_name, links, images, headers, status_code,
        content_type, word_count, load_time_ms, created_at
    FROM scraping_results
    WHERE user_id = ?
    ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, userID)

	if err != nil {
		return nil, fmt.Errorf("error querying results: %w", err)
	}
	defer rows.Close()
	var results []*entity.ScrapingResult
	for rows.Next() {
		result := &entity.ScrapingResult{}
		var linksJSON, imagesJSON, headersJSON, createdAt string
		err := rows.Scan(
			&result.ID, &result.UserID, &result.URL, &result.Title, &result.Description,
			&result.Keywords, &result.Author, &result.Language, &result.Favicon, &result.ImageURL,
			&result.SiteName, &linksJSON, &imagesJSON, &headersJSON, &result.StatusCode,
			&result.ContentType, &result.WordCount, &result.LoadTime, &createdAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		if err := r.unmarshalJSONField(linksJSON, &result.Links); err != nil {
			result.Links = []string{}
		}
		if err := r.unmarshalJSONField(imagesJSON, &result.Images); err != nil {
			result.Images = []string{}
		}
		if err := r.unmarshalJSONField(headersJSON, &result.Headers); err != nil {
			result.Headers = []entity.Header{}
		}
		result.CreatedAt, err = r.parseDateTime(createdAt)

		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return results, nil
}

func (r *scrapingRepository) FindByID(id int64) (*entity.ScrapingResult, error) {
	query := `SELECT 
		id, url, title, description, keywords, author, language, favicon,
		image_url, site_name, links, images, headers, status_code,
		content_type, word_count, load_time_ms, created_at 
	FROM scraping_results 
	WHERE id = ?`

	result := &entity.ScrapingResult{}
	var linksJSON, imagesJSON, headersJSON, createdAt string
	err := r.db.QueryRow(query, id).Scan(
		&result.ID, &result.URL, &result.Title, &result.Description,
		&result.Keywords, &result.Author, &result.Language, &result.Favicon,
		&result.ImageURL, &result.SiteName, &linksJSON, &imagesJSON,
		&headersJSON, &result.StatusCode, &result.ContentType,
		&result.WordCount, &result.LoadTime, &createdAt)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying result by id: %w", err)
	}
	if err := r.unmarshalJSONField(linksJSON, &result.Links); err != nil {
		result.Links = []string{}
	}
	if err := r.unmarshalJSONField(imagesJSON, &result.Images); err != nil {
		result.Images = []string{}
	}
	if err := r.unmarshalJSONField(headersJSON, &result.Headers); err != nil {
		result.Headers = []entity.Header{}
	}
	result.CreatedAt, err = r.parseDateTime(createdAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}
	return result, nil
}

func (r *scrapingRepository) Delete(id int64) error {
	query := `DELETE FROM scraping_results WHERE id = ?`
	_, err := r.db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("error deleting result: %w", err)
	}
	return nil
}

func (r *scrapingRepository) unmarshalJSONField(jsonStr string, target interface{}) error {

	if jsonStr == "" {
		return fmt.Errorf("empty json string")
	}
	return json.Unmarshal([]byte(jsonStr), target)
}

func (r *scrapingRepository) parseDateTime(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		time.DateTime,
	}
	for _, format := range formats {

		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", dateStr)
}
