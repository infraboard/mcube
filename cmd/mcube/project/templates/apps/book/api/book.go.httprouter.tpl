package api

import (
	"net/http"

	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"

	"{{.PKG}}/apps/book"
)

func (h *handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	req := book.NewCreateBookRequest()

	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	set, err := h.service.CreateBook(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) QueryBook(w http.ResponseWriter, r *http.Request) {
	req := book.NewQueryBookRequestFromHTTP(r)
	set, err := h.service.QueryBook(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) DescribeBook(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := book.NewDescribeBookRequest(ctx.PS.ByName("id"))
	ins, err := h.service.DescribeBook(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

func (h *handler) PutBook(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := book.NewPutBookRequest(ctx.PS.ByName("id"))

	if err := request.GetDataFromRequest(r, req.Data); err != nil {
		response.Failed(w, err)
		return
	}

	set, err := h.service.UpdateBook(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) PatchBook(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := book.NewPatchBookRequest(ctx.PS.ByName("id"))

	if err := request.GetDataFromRequest(r, req.Data); err != nil {
		response.Failed(w, err)
		return
	}

	set, err := h.service.UpdateBook(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}

func (h *handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := book.NewDeleteBookRequestWithID(ctx.PS.ByName("id"))
	set, err := h.service.DeleteBook(r.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, set)
}
