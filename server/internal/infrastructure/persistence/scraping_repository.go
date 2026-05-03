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

// Columnas SELECT en el mismo orden que populateResult las escanea (31 columnas).
const selectCols = `
	id, user_id, url,
	title, description, keywords,
	author, language, favicon,
	image_url, site_name,
	links, images, headers,
	status_code, content_type,
	word_count, load_time_ms, created_at,
	canonical_url, robots_directive, x_robots_tag,
	viewport, og_data, twitter_card,
	schema_org, redirect_chain, final_url,
	h1_count, has_multiple_h1, seo_score`

const (
	queryScrapingSave = `INSERT INTO scraping_results (
		user_id, url, title, description, keywords, author, language, favicon,
		image_url, site_name, links, images, headers, status_code,
		content_type, word_count, load_time_ms, created_at,
		canonical_url, robots_directive, x_robots_tag, viewport,
		og_data, twitter_card, schema_org, redirect_chain,
		final_url, h1_count, has_multiple_h1, seo_score
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	queryScrapingFindAll = `SELECT` + selectCols + `
	FROM scraping_results ORDER BY created_at DESC`

	queryScrapingFindByUserID = `SELECT` + selectCols + `
	FROM scraping_results WHERE user_id = ? ORDER BY created_at DESC`

	queryScrapingFindByID = `SELECT` + selectCols + `
	FROM scraping_results WHERE id = ?`

	queryScrapingDelete = `DELETE FROM scraping_results WHERE id = ?`

	queryScrapingFindPaginated = `SELECT` + selectCols + `
	FROM scraping_results WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`

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
	ogDataJSON, err := json.Marshal(result.OGData)
	if err != nil {
		return fmt.Errorf("error marshaling og_data: %w", err)
	}
	twitterCardJSON, err := json.Marshal(result.TwitterCard)
	if err != nil {
		return fmt.Errorf("error marshaling twitter_card: %w", err)
	}
	schemaOrgJSON, err := json.Marshal(result.SchemaOrg)
	if err != nil {
		return fmt.Errorf("error marshaling schema_org: %w", err)
	}
	redirectChainJSON, err := json.Marshal(result.RedirectChain)
	if err != nil {
		return fmt.Errorf("error marshaling redirect_chain: %w", err)
	}

	res, err := r.db.Exec(queryScrapingSave,
		result.UserID, result.URL, result.Title, result.Description,
		result.Keywords, result.Author, result.Language, result.Favicon,
		result.ImageURL, result.SiteName,
		string(linksJSON), string(imagesJSON), string(headersJSON),
		result.StatusCode, result.ContentType, result.WordCount, result.LoadTime, result.CreatedAt,
		result.CanonicalURL, result.RobotsDirective, result.XRobotsTag, result.Viewport,
		string(ogDataJSON), string(twitterCardJSON), string(schemaOrgJSON), string(redirectChainJSON),
		result.FinalURL, result.H1Count, result.HasMultipleH1, result.SEOScore,
	)
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
	return r.collectRows(rows)
}

func (r *scrapingRepository) FindAllByUserID(userID int64) ([]*entity.ScrapingResult, error) {
	rows, err := r.db.Query(queryScrapingFindByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying results: %w", err)
	}
	defer rows.Close()
	return r.collectRows(rows)
}

func (r *scrapingRepository) FindByID(id int64) (*entity.ScrapingResult, error) {
	result, err := r.populateResult(r.db.QueryRow(queryScrapingFindByID, id).Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying result by id: %w", err)
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
	results, err := r.collectRows(rows)
	if err != nil {
		return nil, 0, err
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

// — Helpers —

type scanFunc func(dest ...interface{}) error

// populateResult escanea una fila en ScrapingResult. Acepta rows.Scan o row.Scan.
func (r *scrapingRepository) populateResult(scan scanFunc) (*entity.ScrapingResult, error) {
	result := &entity.ScrapingResult{}
	var (
		linksJSON, imagesJSON, headersJSON string
		ogDataJSON, twitterCardJSON        string
		schemaOrgJSON, redirectChainJSON   string
		createdAt                          string
	)

	if err := scan(
		&result.ID, &result.UserID, &result.URL,
		&result.Title, &result.Description, &result.Keywords,
		&result.Author, &result.Language, &result.Favicon,
		&result.ImageURL, &result.SiteName,
		&linksJSON, &imagesJSON, &headersJSON,
		&result.StatusCode, &result.ContentType,
		&result.WordCount, &result.LoadTime, &createdAt,
		&result.CanonicalURL, &result.RobotsDirective, &result.XRobotsTag,
		&result.Viewport, &ogDataJSON, &twitterCardJSON,
		&schemaOrgJSON, &redirectChainJSON, &result.FinalURL,
		&result.H1Count, &result.HasMultipleH1, &result.SEOScore,
	); err != nil {
		return nil, err
	}

	result.Links = r.unmarshalLinks(linksJSON)
	result.Images = r.unmarshalImages(imagesJSON)

	if err := r.unmarshalJSONField(headersJSON, &result.Headers); err != nil {
		result.Headers = []entity.Header{}
	}
	if err := json.Unmarshal([]byte(orDefault(ogDataJSON, "{}")), &result.OGData); err != nil {
		result.OGData = entity.OGData{}
	}
	if err := json.Unmarshal([]byte(orDefault(twitterCardJSON, "{}")), &result.TwitterCard); err != nil {
		result.TwitterCard = entity.TwitterCard{}
	}
	if err := r.unmarshalJSONField(schemaOrgJSON, &result.SchemaOrg); err != nil {
		result.SchemaOrg = []string{}
	}
	if err := r.unmarshalJSONField(redirectChainJSON, &result.RedirectChain); err != nil {
		result.RedirectChain = []string{}
	}

	var err error
	result.CreatedAt, err = datetime.Parse(createdAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}
	return result, nil
}

func (r *scrapingRepository) collectRows(rows *sql.Rows) ([]*entity.ScrapingResult, error) {
	var results []*entity.ScrapingResult
	for rows.Next() {
		result, err := r.populateResult(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return results, nil
}

// unmarshalLinks acepta el formato nuevo ([]Link) y el antiguo ([]string).
func (r *scrapingRepository) unmarshalLinks(jsonStr string) []entity.Link {
	if jsonStr == "" || jsonStr == "null" || jsonStr == "[]" {
		return []entity.Link{}
	}
	var links []entity.Link
	if err := json.Unmarshal([]byte(jsonStr), &links); err == nil {
		return links
	}
	// Fallback: formato viejo era []string
	var urls []string
	if err := json.Unmarshal([]byte(jsonStr), &urls); err != nil {
		return []entity.Link{}
	}
	result := make([]entity.Link, len(urls))
	for i, u := range urls {
		result[i] = entity.Link{URL: u}
	}
	return result
}

// unmarshalImages acepta el formato nuevo ([]Image) y el antiguo ([]string).
func (r *scrapingRepository) unmarshalImages(jsonStr string) []entity.Image {
	if jsonStr == "" || jsonStr == "null" || jsonStr == "[]" {
		return []entity.Image{}
	}
	var images []entity.Image
	if err := json.Unmarshal([]byte(jsonStr), &images); err == nil {
		return images
	}
	// Fallback: formato viejo era []string
	var urls []string
	if err := json.Unmarshal([]byte(jsonStr), &urls); err != nil {
		return []entity.Image{}
	}
	result := make([]entity.Image, len(urls))
	for i, u := range urls {
		result[i] = entity.Image{Src: u}
	}
	return result
}

func (r *scrapingRepository) unmarshalJSONField(jsonStr string, target interface{}) error {
	if jsonStr == "" {
		return fmt.Errorf("empty json string")
	}
	return json.Unmarshal([]byte(jsonStr), target)
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
