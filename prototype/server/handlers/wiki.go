package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
)

type WikiHandler struct {
	Store *db.Store
}

func (h *WikiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/wiki")
	switch {
	case r.Method == "GET" && path == "/search":
		h.search(w, r)
	case r.Method == "GET" && path == "/tree":
		h.tree(w, r)
	case r.Method == "GET" && path == "/pages":
		h.list(w, r)
	case r.Method == "POST" && path == "/pages":
		h.create(w, r)
	case r.Method == "GET" && strings.HasPrefix(path, "/pages/"):
		h.get(w, r, strings.TrimPrefix(path, "/pages/"))
	case r.Method == "PUT" && strings.HasPrefix(path, "/pages/"):
		h.update(w, r, strings.TrimPrefix(path, "/pages/"))
	case r.Method == "DELETE" && strings.HasPrefix(path, "/pages/"):
		h.delete(w, r, strings.TrimPrefix(path, "/pages/"))
	default:
		http.NotFound(w, r)
	}
}

func (h *WikiHandler) list(w http.ResponseWriter, r *http.Request) {
	pages := h.Store.ListWikiPages()
	if pages == nil {
		pages = []db.WikiPage{}
	}
	writeJSON(w, http.StatusOK, pages)
}

func (h *WikiHandler) create(w http.ResponseWriter, r *http.Request) {
	var p db.WikiPage
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	p.ID = genID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if p.Slug == "" {
		p.Slug = slugify(p.Title)
	}
	if err := h.Store.AddWikiPage(p); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (h *WikiHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	p := h.Store.GetWikiPage(id)
	if p == nil {
		writeErr(w, http.StatusNotFound, "wiki page not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *WikiHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	err := h.Store.UpdateWikiPage(id, func(p *db.WikiPage) {
		if v, ok := updates["title"].(string); ok {
			p.Title = v
		}
		if v, ok := updates["content"].(string); ok {
			p.Content = v
		}
		if v, ok := updates["parent_id"].(string); ok {
			p.ParentID = v
		}
		if v, ok := updates["slug"].(string); ok {
			p.Slug = v
		}
		if v, ok := updates["tags"].([]interface{}); ok {
			tags := make([]string, 0, len(v))
			for _, t := range v {
				if s, ok := t.(string); ok {
					tags = append(tags, s)
				}
			}
			p.Tags = tags
		}
		p.UpdatedAt = time.Now()
	})
	if err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, h.Store.GetWikiPage(id))
}

func (h *WikiHandler) delete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Store.DeleteWikiPage(id); err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *WikiHandler) search(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(r.URL.Query().Get("q"))
	if q == "" {
		writeJSON(w, http.StatusOK, []db.WikiPage{})
		return
	}
	pages := h.Store.SearchWikiPages(q)
	if pages == nil {
		pages = []db.WikiPage{}
	}
	writeJSON(w, http.StatusOK, pages)
}

type wikiTreeNode struct {
	db.WikiPage
	Children []wikiTreeNode `json:"children"`
}

func (h *WikiHandler) tree(w http.ResponseWriter, r *http.Request) {
	pages := h.Store.ListWikiPages()
	byParent := make(map[string][]db.WikiPage)
	for _, p := range pages {
		byParent[p.ParentID] = append(byParent[p.ParentID], p)
	}
	var buildTree func(parentID string) []wikiTreeNode
	buildTree = func(parentID string) []wikiTreeNode {
		children := byParent[parentID]
		nodes := make([]wikiTreeNode, 0, len(children))
		for _, c := range children {
			nodes = append(nodes, wikiTreeNode{WikiPage: c, Children: buildTree(c.ID)})
		}
		return nodes
	}
	writeJSON(w, http.StatusOK, buildTree(""))
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		if r == ' ' || r == '-' || r == '_' {
			return '-'
		}
		return -1
	}, s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}
