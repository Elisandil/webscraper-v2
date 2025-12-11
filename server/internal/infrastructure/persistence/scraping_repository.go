package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/domain/repository"
	"webscraper-v2/internal/infrastructure/database"
	"webscraper-v2/pkg/datetime"
)

const (
	queryScrapingSave = `INSERT INTO scraping_results (
		user_id, url, title, description, keywords, author, language, favicon, 
		image_url, site_name, links, images, headers, status_code, 
		content_type, word_count, load_time_ms, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	queryScrapingFindAll = `SELECT 
		id, url, title, description, keywords, author, language, favicon,
		image_url, site_name, links, images, headers, status_code,
		content_type, word_count, load_time_ms, created_at 
	FROM scraping_results 
	ORDER BY created_at DESC`

	queryScrapingFindByUserID = `
    SELECT 
        id, user_id, url, title, description, keywords, author, language,
        favicon, image_url, site_name, links, images, headers, status_code,
        content_type, word_count, load_time_ms, created_at
    FROM scraping_results
    WHERE user_id = ?
    ORDER BY created_at DESC`

	queryScrapingFindByID = `SELECT 
		id, url, title, description, keywords, author, language, favicon,
		image_url, site_name, links, images, headers, status_code,
		content_type, word_count, load_time_ms, created_at 
	FROM scraping_results 
	WHERE id = ?`

	queryScrapingDelete = `DELETE FROM scraping_results WHERE id = ?`

	queryScrapingFindPaginated = `
    SELECT 
        id, user_id, url, title, description, keywords, author, language,
        favicon, image_url, site_name, links, images, headers, status_code,
        content_type, word_count, load_time_ms, created_at
    FROM scraping_results
    WHERE user_id = ?
    ORDER BY created_at DESC
    LIMIT ? OFFSET ?`

	queryScrapingCount = `SELECT COUNT(*) FROM scraping_results WHERE user_id = ?`
)

type scrapingRepository struct {
	db *database.SQLiteDB
}

func NewScrapingRepository(db *database.SQLiteDB) repository.ScrapingRepository {
	return &scrapingRepository{db: db}
}

func (r *scrapingRepository) Save(result *entity.ScrapingResult) error {
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
	res, err := r.db.Exec(queryScrapingSave,
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
	rows, err := r.db.Query(queryScrapingFindAll)

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
		result.CreatedAt, err = datetime.Parse(createdAt)

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
	rows, err := r.db.Query(queryScrapingFindByUserID, userID)

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
		result.CreatedAt, err = datetime.Parse(createdAt)

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
	result := &entity.ScrapingResult{}
	var linksJSON, imagesJSON, headersJSON, createdAt string
	err := r.db.QueryRow(queryScrapingFindByID, id).Scan(
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
	result.CreatedAt, err = datetime.Parse(createdAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}
	return result, nil
}

func (r *scrapingRepository) Delete(id int64) error {
	_, err := r.db.Exec(queryScrapingDelete, id)

	if err != nil {
		return fmt.Errorf("error deleting result: %w", err)
	}
	return nil
}

func (r *scrapingRepository) FindAllByUserIDPaginated(userID int64, pagination *entity.PaginationRequest) ([]*entity.ScrapingResult, int64, error) {
	totalCount, err := r.CountByUserID(userID)

	if err != nil {
		return nil, 0, fmt.Errorf("error counting results: %w", err)
	}
	if totalCount == 0 {
		return []*entity.ScrapingResult{}, 0, nil
	}
	rows, err := r.db.Query(queryScrapingFindPaginated, userID, pagination.PerPage, pagination.Offset())

	if err != nil {
		return nil, 0, fmt.Errorf("error querying paginated results: %w", err)
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
			return nil, 0, fmt.Errorf("error scanning row: %w", err)
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
		result.CreatedAt, err = datetime.Parse(createdAt)

		if err != nil {
			return nil, 0, fmt.Errorf("error parsing created_at: %w", err)
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}
	return results, totalCount, nil
}

func (r *scrapingRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := r.db.QueryRow(queryScrapingCount, userID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("error counting results by user ID: %w", err)
	}
	return count, nil
}

func (r *scrapingRepository) unmarshalJSONField(jsonStr string, target interface{}) error {

	if jsonStr == "" {
		return fmt.Errorf("empty json string")
	}
	return json.Unmarshal([]byte(jsonStr), target)
}
