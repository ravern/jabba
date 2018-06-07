package http

import (
	"net/http"
	"net/url"

	"github.com/ravernkoh/jabba/model"
)

// Root renders the index page.
func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse("https://google.com")
	executeTemplate(w, "layout.html", nil, "index.html", struct {
		Links []*model.Link
	}{
		Links: []*model.Link{
			{
				Slug:  "jdheif",
				URL:   url,
				Title: "Report on Wildlife - Google Docs",
			},
			{
				Slug:  "ugbrkq",
				URL:   url,
				Title: "Skeleton: Simple & Responsive CSS Boilerplate",
			},
		},
	})
}
