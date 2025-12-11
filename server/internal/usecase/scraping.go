package usecase

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/domain/repository"
	"webscraper-v2/internal/infrastructure/config"
	pkgerrors "webscraper-v2/pkg/errors"
	"webscraper-v2/pkg/validator"

	"golang.org/x/net/html"
)

type ScrapingUseCase struct {
	repo      repository.ScrapingRepository
	config    *config.Config
	client    *http.Client
	validator *validator.Validator
}

func NewScrapingUseCase(repo repository.ScrapingRepository, cfg *config.Config) *ScrapingUseCase {
	return &ScrapingUseCase{
		repo:      repo,
		config:    cfg,
		validator: validator.NewValidator(),
		client: &http.Client{
			Timeout: time.Duration(cfg.Scraping.Timeout) * time.Second,
		},
	}
}

func (uc *ScrapingUseCase) ScrapeURL(ctx context.Context, targetURL string, userID int64) (*entity.ScrapingResult, error) {

	if err := uc.validator.ValidateURL(targetURL); err != nil {
		return nil, pkgerrors.ValidationError(err.Error())
	}

	startTime := time.Now()

	// Usar contexto para permitir cancelación
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to create request", err)
	}

	req.Header.Set("User-Agent", uc.config.Scraping.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := uc.client.Do(req)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to fetch URL", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to read response body", err)
	}

	loadTime := time.Since(startTime).Milliseconds()
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, pkgerrors.InternalError("failed to parse HTML", err)
	}

	result := &entity.ScrapingResult{
		UserID:      userID,
		URL:         targetURL,
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		LoadTime:    loadTime,
		CreatedAt:   time.Now(),
	}

	uc.extractMetadata(doc, result)
	uc.extractLinks(doc, result, targetURL)
	uc.extractImages(doc, result, targetURL)
	uc.extractHeaders(doc, result)
	uc.extractFavicon(ctx, targetURL, result)
	uc.calculateWordCount(string(body), result)
	if err := uc.repo.Save(result); err != nil {
		return nil, pkgerrors.DatabaseError("save scraping result", err)
	}

	return result, nil
}

func (uc *ScrapingUseCase) GetAllResults(userID int64) ([]*entity.ScrapingResult, error) {
	results, err := uc.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, pkgerrors.DatabaseError("get all results", err)
	}
	return results, nil
}

func (uc *ScrapingUseCase) GetResult(id int64, userID int64) (*entity.ScrapingResult, error) {
	result, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, pkgerrors.DatabaseError("get result", err)
	}

	if result == nil {
		return nil, pkgerrors.NotFoundError("result")
	}

	if result.UserID != userID {
		return nil, pkgerrors.New(
			pkgerrors.CodeAuthorization,
			"unauthorized: user does not own this result",
			pkgerrors.ErrUnauthorized,
		)
	}

	return result, nil
}

func (uc *ScrapingUseCase) DeleteResult(id int64, userID int64) error {
	result, err := uc.repo.FindByID(id)
	if err != nil {
		return pkgerrors.DatabaseError("get result for deletion", err)
	}
	if result == nil {
		return pkgerrors.NotFoundError("result")
	}

	if result.UserID != userID {
		return pkgerrors.New(
			pkgerrors.CodeAuthorization,
			"unauthorized: user does not own this result",
			pkgerrors.ErrUnauthorized,
		)
	}

	if err := uc.repo.Delete(id); err != nil {
		return pkgerrors.DatabaseError("delete result", err)
	}

	return nil

}

func (uc *ScrapingUseCase) GetAllResultsPaginated(userID int64, page, perPage int) (*entity.PaginatedScrapingResults, error) {
	paginationReq := entity.NewPaginationRequest(page, perPage)

	results, totalCount, err := uc.repo.FindAllByUserIDPaginated(userID, paginationReq)
	if err != nil {
		return nil, pkgerrors.DatabaseError("get paginated results", err)
	}

	paginationResp := entity.NewPaginationResponse(paginationReq.Page, paginationReq.PerPage, totalCount)

	return &entity.PaginatedScrapingResults{
		Data:       results,
		Pagination: paginationResp,
	}, nil
}

func (uc *ScrapingUseCase) extractMetadata(n *html.Node, result *entity.ScrapingResult) {
	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "title":
				if child := node.FirstChild; child != nil && child.Type == html.TextNode {
					result.Title = strings.TrimSpace(child.Data)
				}
			case "meta":
				uc.extractMetaTag(node, result)
			}
		}
	})
}

func (uc *ScrapingUseCase) extractMetaTag(node *html.Node, result *entity.ScrapingResult) {

	if node == nil || node.Type != html.ElementNode || node.Data != "meta" {
		return
	}

	var name, property, content string

	for _, attr := range node.Attr {
		switch attr.Key {
		case "name":
			name = strings.ToLower(attr.Val)
		case "property":
			property = strings.ToLower(attr.Val)
		case "content":
			content = attr.Val
		}
	}

	switch {
	case name == "description" || property == "og:description":
		if result.Description == "" {
			result.Description = content
		}
	case name == "keywords":
		result.Keywords = content
	case name == "author":
		result.Author = content
	case name == "language" || property == "og:locale":
		result.Language = content
	case property == "og:image":
		result.ImageURL = content
	case property == "og:site_name":
		result.SiteName = content
	}
}

func (uc *ScrapingUseCase) extractLinks(n *html.Node, result *entity.ScrapingResult, baseURL string) {
	linkMap := make(map[string]bool)

	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					href := strings.TrimSpace(attr.Val)
					if href != "" && !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "javascript:") {
						absoluteURL := uc.resolveURL(baseURL, href)
						if absoluteURL != "" && !linkMap[absoluteURL] {
							linkMap[absoluteURL] = true
							result.Links = append(result.Links, absoluteURL)
						}
					}
					break
				}
			}
		}
	})
}

func (uc *ScrapingUseCase) extractImages(n *html.Node, result *entity.ScrapingResult, baseURL string) {
	imageMap := make(map[string]bool)

	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "img" {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					src := strings.TrimSpace(attr.Val)
					if src != "" {
						absoluteURL := uc.resolveURL(baseURL, src)
						if absoluteURL != "" && !imageMap[absoluteURL] {
							imageMap[absoluteURL] = true
							result.Images = append(result.Images, absoluteURL)
						}
					}
					break
				}
			}
		}
	})
}

func (uc *ScrapingUseCase) extractHeaders(n *html.Node, result *entity.ScrapingResult) {
	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				if child := node.FirstChild; child != nil && child.Type == html.TextNode {
					level := int(node.Data[1] - '0')
					text := strings.TrimSpace(child.Data)
					if text != "" {
						result.Headers = append(result.Headers, entity.Header{
							Level: level,
							Text:  text,
						})
					}
				}
			}
		}
	})
}

func (uc *ScrapingUseCase) extractFavicon(ctx context.Context, targetURL string, result *entity.ScrapingResult) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return
	}

	faviconURLs := []string{
		fmt.Sprintf("%s://%s/favicon.ico", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/favicon.png", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/apple-touch-icon.png", parsedURL.Scheme, parsedURL.Host),
	}

	for _, faviconURL := range faviconURLs {
		if uc.checkURLExists(ctx, faviconURL) {
			result.Favicon = faviconURL
			break
		}
	}
}

func (uc *ScrapingUseCase) checkURLExists(ctx context.Context, url string) bool {
	// Crear contexto con timeout específico para favicon check
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return false
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (uc *ScrapingUseCase) calculateWordCount(content string, result *entity.ScrapingResult) {
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(content, " ")
	words := strings.Fields(text)
	result.WordCount = len(words)
}

func (uc *ScrapingUseCase) resolveURL(base, href string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}

	hrefURL, err := url.Parse(href)
	if err != nil {
		return ""
	}

	return baseURL.ResolveReference(hrefURL).String()
}

func (uc *ScrapingUseCase) traverseNode(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		uc.traverseNode(c, fn)
	}
}
