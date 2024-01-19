package ui

import "github.com/egoodhall/go-links/pkg/config"

type Data struct {
	Query   string
	GoLinks []GoLink
}

type GoLink struct {
	Title       string
	Description string
	Aliases     []string
	Links       []Link
}

type Link struct {
	From string
	To   string
}

func NewGoLinks(frm config.Targets) []GoLink {
	to := make([]GoLink, len(frm))
	for i, tgt := range frm {
		to[i] = GoLink{
			Title:       tgt.Title,
			Description: tgt.Description,
			Aliases:     tgt.Aliases,
			Links:       make([]Link, len(tgt.Urls)*len(tgt.Aliases)),
		}
		targets := to[i].Links
		i := 0
		for _, url := range tgt.Urls {
			for _, goLink := range url.GetGoLinks(tgt.Aliases) {
				targets[i] = Link{goLink, string(url)}
				i++
			}
		}
	}
	return to
}

type ArgFormData struct {
	From string
	To   string
}
