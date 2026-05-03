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

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

// ResultNotifier is notified after each scraping result is persisted.
// Implemented by the SSE hub in the presentation layer.
type ResultNotifier interface {
	Notify(userID int64)
}

type ScrapingUseCase struct {
	repo      repository.ScrapingRepository
	config    *config.Config
	validator *validator.Validator
	notifier  ResultNotifier
}

func NewScrapingUseCase(repo repository.ScrapingRepository, cfg *config.Config) *ScrapingUseCase {
	return &ScrapingUseCase{
		repo:      repo,
		config:    cfg,
		validator: validator.NewValidator(),
	}
}

func (uc *ScrapingUseCase) SetNotifier(n ResultNotifier) {
	uc.notifier = n
}

func (uc *ScrapingUseCase) ScrapeURL(ctx context.Context, targetURL string, userID int64) (*entity.ScrapingResult, error) {
	if err := uc.validator.ValidateURL(targetURL); err != nil {
		return nil, pkgerrors.ValidationError(err.Error())
	}

	startTime := time.Now()
	var redirectChain []string

	client := &http.Client{
		Timeout: time.Duration(uc.config.Scraping.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			redirectChain = append(redirectChain, req.URL.String())
			return nil
		},
	}

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to create request", err)
	}

	req.Header.Set("User-Agent", uc.config.Scraping.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
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
		UserID:        userID,
		URL:           targetURL,
		StatusCode:    resp.StatusCode,
		ContentType:   resp.Header.Get("Content-Type"),
		XRobotsTag:    resp.Header.Get("X-Robots-Tag"),
		LoadTime:      loadTime,
		RedirectChain: redirectChain,
		FinalURL:      resp.Request.URL.String(),
		CreatedAt:     time.Now(),
	}

	uc.extractMetadata(doc, result)
	uc.extractCanonical(doc, result)
	uc.extractSchemaOrg(doc, result)
	uc.extractLinks(doc, result, targetURL)
	uc.extractImages(doc, result, targetURL)
	uc.extractHeaders(doc, result)
	uc.validateHeadings(result)
	uc.extractFavicon(ctx, targetURL, result)
	uc.calculateWordCount(string(body), result)
	uc.calculateSEOScore(result)

	if err := uc.repo.Save(result); err != nil {
		return nil, pkgerrors.DatabaseError("save scraping result", err)
	}

	if uc.notifier != nil && userID != 0 {
		uc.notifier.Notify(userID)
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

// — Extraction —

func (uc *ScrapingUseCase) extractMetadata(n *html.Node, result *entity.ScrapingResult) {
	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "title":
				if result.Title == "" {
					result.Title = strings.TrimSpace(uc.getTextContent(node))
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
	case name == "description":
		if result.Description == "" {
			result.Description = content
		}
	case name == "keywords":
		result.Keywords = content
	case name == "author":
		result.Author = content
	case name == "language":
		if result.Language == "" {
			result.Language = content
		}
	case name == "robots":
		result.RobotsDirective = content
	case name == "viewport":
		result.Viewport = content

	// Open Graph
	case property == "og:title":
		result.OGData.Title = content
	case property == "og:url":
		result.OGData.URL = content
	case property == "og:type":
		result.OGData.Type = content
	case property == "og:image":
		result.OGData.Image = content
		result.ImageURL = content
	case property == "og:description":
		result.OGData.Description = content
		if result.Description == "" {
			result.Description = content
		}
	case property == "og:site_name":
		result.OGData.SiteName = content
		result.SiteName = content
	case property == "og:locale":
		result.OGData.Locale = content
		if result.Language == "" {
			result.Language = content
		}
	}

	// Twitter Card (puede venir como name= o property=)
	twitterKey := ""
	if strings.HasPrefix(name, "twitter:") {
		twitterKey = strings.TrimPrefix(name, "twitter:")
	} else if strings.HasPrefix(property, "twitter:") {
		twitterKey = strings.TrimPrefix(property, "twitter:")
	}
	if twitterKey != "" {
		switch twitterKey {
		case "card":
			result.TwitterCard.Card = content
		case "title":
			result.TwitterCard.Title = content
		case "description":
			result.TwitterCard.Description = content
		case "image":
			result.TwitterCard.Image = content
		case "site":
			result.TwitterCard.Site = content
		}
	}
}

func (uc *ScrapingUseCase) extractCanonical(n *html.Node, result *entity.ScrapingResult) {
	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "link" {
			var rel, href string
			for _, attr := range node.Attr {
				switch attr.Key {
				case "rel":
					rel = strings.ToLower(strings.TrimSpace(attr.Val))
				case "href":
					href = strings.TrimSpace(attr.Val)
				}
			}
			if rel == "canonical" && href != "" {
				result.CanonicalURL = href
			}
		}
	})
}

func (uc *ScrapingUseCase) extractSchemaOrg(n *html.Node, result *entity.ScrapingResult) {
	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "script" {
			for _, attr := range node.Attr {
				if attr.Key == "type" && attr.Val == "application/ld+json" {
					// getTextContent maneja whitespace-nodes y múltiples nodos de texto
					raw := uc.getTextContent(node)
					if raw != "" {
						result.SchemaOrg = append(result.SchemaOrg, raw)
					}
					break
				}
			}
		}
	})
}

func (uc *ScrapingUseCase) extractLinks(n *html.Node, result *entity.ScrapingResult, baseURL string) {
	linkMap := make(map[string]bool)
	baseParsed, _ := url.Parse(baseURL)

	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			var href, rel string
			for _, attr := range node.Attr {
				switch attr.Key {
				case "href":
					href = strings.TrimSpace(attr.Val)
				case "rel":
					rel = attr.Val
				}
			}
			if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") {
				return
			}
			absoluteURL := uc.resolveURL(baseURL, href)
			if absoluteURL == "" || linkMap[absoluteURL] {
				return
			}
			linkMap[absoluteURL] = true

			isInternal := false
			if parsed, err := url.Parse(absoluteURL); err == nil && baseParsed != nil {
				isInternal = parsed.Host == baseParsed.Host
			}

			result.Links = append(result.Links, entity.Link{
				URL:        absoluteURL,
				AnchorText: uc.getTextContent(node),
				Rel:        rel,
				IsInternal: isInternal,
			})
		}
	})
}

func (uc *ScrapingUseCase) extractImages(n *html.Node, result *entity.ScrapingResult, baseURL string) {
	imageMap := make(map[string]bool)

	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "img" {
			var src, alt, title string
			for _, attr := range node.Attr {
				switch attr.Key {
				case "src":
					src = strings.TrimSpace(attr.Val)
				case "alt":
					alt = attr.Val
				case "title":
					title = attr.Val
				}
			}
			if src == "" {
				return
			}
			absoluteURL := uc.resolveURL(baseURL, src)
			if absoluteURL == "" || imageMap[absoluteURL] {
				return
			}
			imageMap[absoluteURL] = true
			result.Images = append(result.Images, entity.Image{
				Src:   absoluteURL,
				Alt:   alt,
				Title: title,
			})
		}
	})
}

func (uc *ScrapingUseCase) extractHeaders(n *html.Node, result *entity.ScrapingResult) {
	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				level := int(node.Data[1] - '0')
				text := uc.getTextContent(node)
				if text != "" {
					result.Headers = append(result.Headers, entity.Header{
						Level: level,
						Text:  text,
					})
				}
			}
		}
	})
}

func (uc *ScrapingUseCase) validateHeadings(result *entity.ScrapingResult) {
	for _, h := range result.Headers {
		if h.Level == 1 {
			result.H1Count++
		}
	}
	result.HasMultipleH1 = result.H1Count > 1
}

func (uc *ScrapingUseCase) calculateSEOScore(result *entity.ScrapingResult) {
	score := 0

	// Title: ideal 50-60 chars (+20)
	titleLen := len(result.Title)
	switch {
	case titleLen >= 50 && titleLen <= 60:
		score += 20
	case titleLen >= 30 && titleLen <= 70:
		score += 10
	case titleLen > 0:
		score += 5
	}

	// Description: ideal 150-160 chars (+15)
	descLen := len(result.Description)
	switch {
	case descLen >= 150 && descLen <= 160:
		score += 15
	case descLen >= 100 && descLen <= 180:
		score += 8
	case descLen > 0:
		score += 3
	}

	// H1: exactly one (+15)
	switch {
	case result.H1Count == 1:
		score += 15
	case result.H1Count > 1:
		score += 5
	}

	// Canonical URL (+10)
	if result.CanonicalURL != "" {
		score += 10
	}

	// Not noindex (+10)
	if !strings.Contains(strings.ToLower(result.RobotsDirective), "noindex") &&
		!strings.Contains(strings.ToLower(result.XRobotsTag), "noindex") {
		score += 10
	}

	// JSON-LD present (+10)
	if len(result.SchemaOrg) > 0 {
		score += 10
	}

	// Images with alt text (+10)
	if len(result.Images) == 0 {
		score += 10
	} else {
		withAlt := 0
		for _, img := range result.Images {
			if img.Alt != "" {
				withAlt++
			}
		}
		score += int(float64(withAlt) / float64(len(result.Images)) * 10)
	}

	// No redirect chain (+10)
	if len(result.RedirectChain) == 0 {
		score += 10
	}

	result.SEOScore = score
}

func (uc *ScrapingUseCase) extractFavicon(ctx context.Context, targetURL string, result *entity.ScrapingResult) {
	// Usar el host final (post-redirect) para la búsqueda del favicon
	probeURL := targetURL
	if result.FinalURL != "" {
		probeURL = result.FinalURL
	}
	parsedURL, err := url.Parse(probeURL)
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

func (uc *ScrapingUseCase) checkURLExists(ctx context.Context, targetURL string) bool {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", targetURL, nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (uc *ScrapingUseCase) calculateWordCount(content string, result *entity.ScrapingResult) {
	text := htmlTagRe.ReplaceAllString(content, " ")
	result.WordCount = len(strings.Fields(text))
}

// — Helpers —

func (uc *ScrapingUseCase) getTextContent(n *html.Node) string {
	var sb strings.Builder
	uc.traverseNode(n, func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
	})
	return strings.TrimSpace(sb.String())
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
