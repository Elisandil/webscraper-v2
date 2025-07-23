package usecase

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"webscraper/config"
	"webscraper/domain/entity"
	"webscraper/domain/repository"

	"golang.org/x/net/html"
)

type ScrapingUseCase struct {
	repo   repository.ScrapingRepository
	config *config.Config
	client *http.Client
}

func NewScrapingUseCase(repo repository.ScrapingRepository, cfg *config.Config) *ScrapingUseCase {
	return &ScrapingUseCase{
		repo:   repo,
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.Scraping.Timeout) * time.Second,
		},
	}
}

func (uc *ScrapingUseCase) ScrapeURL(targetURL string) (*entity.ScrapingResult, error) {
	startTime := time.Now()

	// Crear la solicitud HTTP
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", uc.config.Scraping.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// addRealisticHeaders(req)

	resp, err := uc.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	loadTime := time.Since(startTime).Milliseconds()

	doc, err := html.Parse(strings.NewReader(string(body)))

	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	// Extraer datos
	result := &entity.ScrapingResult{
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
	uc.extractFavicon(targetURL, result)
	uc.calculateWordCount(string(body), result)

	// Guardar en la base de datos, PERSISTENCIA
	if err := uc.repo.Save(result); err != nil {
		return nil, fmt.Errorf("error saving result: %w", err)
	}

	return result, nil
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
	var name, property, content string

	for _, attr := range node.Attr { // for in range, en la cual el índice no es necesario (_, attr)
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

func (uc *ScrapingUseCase) extractFavicon(targetURL string, result *entity.ScrapingResult) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return
	}

	// Intentar encontrar favicon en ubicaciones comunes
	faviconURLs := []string{
		fmt.Sprintf("%s://%s/favicon.ico", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/favicon.png", parsedURL.Scheme, parsedURL.Host),
		fmt.Sprintf("%s://%s/apple-touch-icon.png", parsedURL.Scheme, parsedURL.Host),
	}

	for _, faviconURL := range faviconURLs {
		if uc.checkURLExists(faviconURL) {
			result.Favicon = faviconURL
			break
		}
	}
}

func (uc *ScrapingUseCase) checkURLExists(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (uc *ScrapingUseCase) calculateWordCount(content string, result *entity.ScrapingResult) {
	// Remover etiquetas HTML
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(content, " ")

	// Contar palabras
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

func (uc *ScrapingUseCase) GetAllResults() ([]*entity.ScrapingResult, error) {
	return uc.repo.FindAll()
}

func (uc *ScrapingUseCase) GetResult(id int64) (*entity.ScrapingResult, error) {
	return uc.repo.FindByID(id)
}

func (uc *ScrapingUseCase) DeleteResult(id int64) error {
	return uc.repo.Delete(id)
}

//func addRealisticHeaders(req *http.Request) {
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "+
//		"(KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
//	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
//	req.Header.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8")
//	req.Header.Set("Connection", "keep-alive")
//	req.Header.Set("Upgrade-Insecure-Requests", "1")
//	req.Header.Set("Cache-Control", "no-cache")
//}
