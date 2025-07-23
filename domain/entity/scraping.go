package entity

import "time"

type ScrapingResult struct {
	ID          int64     `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Keywords    string    `json:"keywords"`
	Author      string    `json:"author"`
	Language    string    `json:"language"`
	Favicon     string    `json:"favicon"`
	ImageURL    string    `json:"image_url"`
	SiteName    string    `json:"site_name"`
	Links       []string  `json:"links"`
	Images      []string  `json:"images"`
	Headers     []Header  `json:"headers"`
	StatusCode  int       `json:"status_code"`
	ContentType string    `json:"content_type"`
	WordCount   int       `json:"word_count"`
	LoadTime    int64     `json:"load_time_ms"`
	CreatedAt   time.Time `json:"created_at"`
}

type Header struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}
