package api

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"

	"{{.PKG}}/apps/book"
)

func (h *handler) CreateBook(c *gin.Context) {
	req := book.NewCreateBookRequest()

	if err := c.BindJSON(req); err != nil {
		response.Failed(c.Writer, err)
		return
	}

	set, err := h.service.CreateBook(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, set)
}

func (h *handler) QueryBook(c *gin.Context) {
	req := book.NewQueryBookRequestFromHTTP(c.Request)
	set, err := h.service.QueryBook(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}

func (h *handler) DescribeBook(c *gin.Context) {
	req := book.NewDescribeBookRequest(c.Param("id"))
	ins, err := h.service.DescribeBook(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}

	response.Success(c.Writer, ins)
}

func (h *handler) PutBook(c *gin.Context) {
	req := book.NewPutBookRequest(c.Param("id"))

	if err := c.BindJSON(req.Data); err != nil {
		response.Failed(c.Writer, err)
		return
	}

	set, err := h.service.UpdateBook(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}

func (h *handler) PatchBook(c *gin.Context) {
	req := book.NewPatchBookRequest(c.Param("id"))

	if err := c.BindJSON(req.Data); err != nil {
		response.Failed(c.Writer, err)
		return
	}

	set, err := h.service.UpdateBook(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}

func (h *handler) DeleteBook(c *gin.Context) {
	req := book.NewDeleteBookRequestWithID(c.Param("id"))
	set, err := h.service.DeleteBook(c.Request.Context(), req)
	if err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, set)
}
