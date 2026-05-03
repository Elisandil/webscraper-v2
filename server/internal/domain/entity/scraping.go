package entity

import "time"

type Link struct {
	URL        string `json:"url"`
	AnchorText string `json:"anchor_text"`
	Rel        string `json:"rel"`
	IsInternal bool   `json:"is_internal"`
}

type Image struct {
	Src   string `json:"src"`
	Alt   string `json:"alt"`
	Title string `json:"title"`
}

type OGData struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Image       string `json:"image"`
	Description string `json:"description"`
	SiteName    string `json:"site_name"`
	Locale      string `json:"locale"`
}

type TwitterCard struct {
	Card        string `json:"card"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Site        string `json:"site"`
}

type ScrapingResult struct {
	ID              int64       `json:"id"`
	UserID          int64       `json:"user_id"`
	URL             string      `json:"url"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Keywords        string      `json:"keywords"`
	Author          string      `json:"author"`
	Language        string      `json:"language"`
	Favicon         string      `json:"favicon"`
	ImageURL        string      `json:"image_url"`
	SiteName        string      `json:"site_name"`
	Links           []Link      `json:"links"`
	Images          []Image     `json:"images"`
	Headers         []Header    `json:"headers"`
	StatusCode      int         `json:"status_code"`
	ContentType     string      `json:"content_type"`
	WordCount       int         `json:"word_count"`
	LoadTime        int64       `json:"load_time_ms"`
	CanonicalURL    string      `json:"canonical_url"`
	RobotsDirective string      `json:"robots_directive"`
	XRobotsTag      string      `json:"x_robots_tag"`
	Viewport        string      `json:"viewport"`
	OGData          OGData      `json:"og_data"`
	TwitterCard     TwitterCard `json:"twitter_card"`
	SchemaOrg       []string    `json:"schema_org"`
	RedirectChain   []string    `json:"redirect_chain"`
	FinalURL        string      `json:"final_url"`
	H1Count         int         `json:"h1_count"`
	HasMultipleH1   bool        `json:"has_multiple_h1"`
	SEOScore        int         `json:"seo_score"`
	CreatedAt       time.Time   `json:"created_at"`
}

type Header struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}
