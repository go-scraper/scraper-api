# Scraper API

## How to run using Docker

* Run `docker-compose up --build`

## How to run tests

* Run `go test -coverprofile=coverage.out ./...`

## API Documentation

#### Request
1. Srape a URL

> * Request type: `GET`
> * URL: `http://localhost:8080/scrape?url=<URL to scrape>`
> * Parameters:
>    * `URL` - URL to scrape

#### Response

1. Success response

```json
{
    "request_id": "20250102001144-oTtblaYW",
    "pagination": {
        "page_size": 10,
        "current_page": 1,
        "total_pages": 5,
        "next_page": "/scrape/20250102001144-oTtblaYW/2"
    },
    "scraped": {
        "html_version": "HTML 5",
        "title": "Facebook – log in or sign up",
        "headings": {
            "h2": 1
        },
        "contains_login_form": true,
        "total_urls": 48,
        "internal_urls": 24,
        "external_urls": 24,
        "paginated": {
            "inaccessible_urls": 0,
            "urls": [
                {
                    "url": "https://facebook.com",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://www.facebook.com/recover/initiate/?privacy_mutation_token=eyJ0eXBlIjowLCJjcmVhdGlvbl90aW1lIjoxNzM1NzU2OTA0LCJjYWxsc2l0ZV9pZCI6MzgxMjI5MDc5NTc1OTQ2fQ%3D%3D&ars=facebook_login&next",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://facebook.com/r.php?entry_point=login",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://facebook.com/pages/create/?ref_type=registration_form",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://www.facebook.com/",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://si-lk.facebook.com/",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://ta-in.facebook.com/",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://es-la.facebook.com/",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://de-de.facebook.com/",
                    "http_status": 200,
                    "error": null
                },
                {
                    "url": "https://it-it.facebook.com/",
                    "http_status": 200,
                    "error": null
                }
            ]
        }
    }
}
```

2. Error response

```json
{
    "error": "Error message"
}
```
