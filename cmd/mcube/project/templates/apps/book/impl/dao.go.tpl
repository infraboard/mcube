package impl

import (
	"context"
	"fmt"

	"{{.PKG}}/apps/book"
)

func (s *service) deleteBook(ctx context.Context, ins *book.Book) error {
	if ins == nil || ins.Id == "" {
		return fmt.Errorf("book is nil")
	}

	stmt, err := s.db.Prepare(deleteBookSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ins.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) updateBook(ctx context.Context, ins *book.Book) error {
	stmt, err := s.db.Prepare(updateBookSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		ins.UpdateAt, ins.UpdateBy, ins.Data.Name, ins.Data.Author, ins.Id,
	)
	if err != nil {
		return err
	}

	return nil
}
