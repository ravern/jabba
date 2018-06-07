package http

import (
	"net/http"

	"github.com/ravernkoh/jabba/model"
)

// Root renders the index page.
func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "layout.html", nil, "index.html", struct {
		Links []*model.Link
	}{
		Links: []*model.Link{
			{
				Slug:  "jdheif",
				Title: "Report on Wildlife - Google Docs",
			},
			{
				Slug:  "ugbrkq",
				Title: "Skeleton: Simple & Responsive CSS Boilerplate",
			},
		},
	})
}
