package impl

import (
	"context"
{{ if $.EnableMySQL -}}
	"database/sql"
{{- end }}

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/pb/request"
{{ if $.EnableMySQL -}}
	"github.com/infraboard/mcube/sqlbuilder"
{{- end }}

	"{{.PKG}}/apps/book"
)

{{ if $.EnableMySQL -}}
func (s *service) CreateBook(ctx context.Context, req *book.CreateBookRequest) (
	*book.Book, error) {
	ins, err := book.NewBook(req)
	if err != nil {
		return nil, exception.NewBadRequest("validate create book error, %s", err)
	}

	stmt, err := s.db.Prepare(insertBookSQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		ins.Id, ins.CreateAt, ins.Data.CreateBy, ins.UpdateAt, ins.UpdateBy,
		ins.Data.Name, ins.Data.Author,
	)
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (s *service) QueryBook(ctx context.Context, req *book.QueryBookRequest) (
	*book.BookSet, error) {
	query := sqlbuilder.NewQuery(queryBookSQL)

	// 支持关键字参数
	if req.Keywords != "" {
		query.Where("name LIKE ? OR author = ?",
			"%"+req.Keywords+"%",
			req.Keywords,
		)
	}

	querySQL, args := query.Order("create_at").Desc().Limit(req.Page.ComputeOffset(), uint(req.Page.PageSize)).BuildQuery()
	s.Debug().Msgf("sql: %s, args: %v", querySQL, args)

	queryStmt, err := s.db.Prepare(querySQL)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare query book error, %s", err.Error())
	}
	defer queryStmt.Close()

	rows, err := queryStmt.Query(args...)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	set := book.NewBookSet()
	for rows.Next() {
		ins := book.NewDefaultBook()
		err := rows.Scan(
			&ins.Id, &ins.CreateAt, &ins.Data.CreateBy, &ins.UpdateAt, &ins.UpdateBy,
			&ins.Data.Name, &ins.Data.Author,
		)
		if err != nil {
			return nil, exception.NewInternalServerError("query book error, %s", err.Error())
		}
		set.Add(ins)
	}

	// 获取total SELECT COUNT(*) FROMT t Where ....
	countSQL, args := query.BuildCount()
	countStmt, err := s.db.Prepare(countSQL)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	defer countStmt.Close()
	err = countStmt.QueryRow(args...).Scan(&set.Total)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	return set, nil
}

func (s *service) DescribeBook(ctx context.Context, req *book.DescribeBookRequest) (
	*book.Book, error) {
	query := sqlbuilder.NewQuery(queryBookSQL)
	querySQL, args := query.Where("id = ?", req.Id).BuildQuery()
	s.Debug().Msgf("sql: %s", querySQL)

	queryStmt, err := s.db.Prepare(querySQL)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare query book error, %s", err.Error())
	}
	defer queryStmt.Close()

	ins := book.NewDefaultBook()
	err = queryStmt.QueryRow(args...).Scan(
		&ins.Id, &ins.CreateAt, &ins.Data.CreateBy, &ins.UpdateAt, &ins.UpdateBy,
		&ins.Data.Name, &ins.Data.Author,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("%s not found", req.Id)
		}
		return nil, exception.NewInternalServerError("describe book error, %s", err.Error())
	}

	return ins, nil
}

func (s *service) UpdateBook(ctx context.Context, req *book.UpdateBookRequest) (
	*book.Book, error) {
	ins, err := s.DescribeBook(ctx, book.NewDescribeBookRequest(req.Id))
	if err != nil {
		return nil, err
	}

	switch req.UpdateMode {
	case request.UpdateMode_PUT:
		ins.Update(req)
	case request.UpdateMode_PATCH:
		err := ins.Patch(req)
		if err != nil {
			return nil, err
		}
	}

	// 校验更新后数据合法性
	if err := ins.Data.Validate(); err != nil {
		return nil, err
	}

	if err := s.updateBook(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}

func (s *service) DeleteBook(ctx context.Context, req *book.DeleteBookRequest) (
	*book.Book, error) {
	ins, err := s.DescribeBook(ctx, book.NewDescribeBookRequest(req.Id))
	if err != nil {
		return nil, err
	}

	if err := s.deleteBook(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}
{{- end }}

{{ if $.EnableMongoDB -}}
func (s *service) CreateBook(ctx context.Context, req *book.CreateBookRequest) (
	*book.Book, error) {
	ins, err := book.NewBook(req)
	if err != nil {
		return nil, exception.NewBadRequest("validate create book error, %s", err)
	}

	if err := s.save(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}

func (s *service) DescribeBook(ctx context.Context, req *book.DescribeBookRequest) (
	*book.Book, error) {
	return s.get(ctx, req.Id)
}

func (s *service) QueryBook(ctx context.Context, req *book.QueryBookRequest) (
	*book.BookSet, error) {
	query := newQueryBookRequest(req)
	return s.query(ctx, query)
}

func (s *service) UpdateBook(ctx context.Context, req *book.UpdateBookRequest) (
	*book.Book, error) {
	ins, err := s.DescribeBook(ctx, book.NewDescribeBookRequest(req.Id))
	if err != nil {
		return nil, err
	}

	switch req.UpdateMode {
	case request.UpdateMode_PUT:
		ins.Update(req)
	case request.UpdateMode_PATCH:
		err := ins.Patch(req)
		if err != nil {
			return nil, err
		}
	}

	// 校验更新后数据合法性
	if err := ins.Data.Validate(); err != nil {
		return nil, err
	}

	if err := s.update(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}

func (s *service) DeleteBook(ctx context.Context, req *book.DeleteBookRequest) (
	*book.Book, error) {
	ins, err := s.DescribeBook(ctx, book.NewDescribeBookRequest(req.Id))
	if err != nil {
		return nil, err
	}

	if err := s.deleteBook(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}
{{- end }}
