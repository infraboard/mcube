package impl

const (
	insertBookSQL = `INSERT INTO books (
		id,create_at,create_by,update_at,update_by,name,author
	) VALUES (?,?,?,?,?,?,?);`

	updateBookSQL = `UPDATE books SET update_at=?,update_by=?,name=?,author=? WHERE id =?`

	queryBookSQL = `SELECT * FROM books`

	deleteBookSQL = `DELETE FROM books WHERE id = ?`
)