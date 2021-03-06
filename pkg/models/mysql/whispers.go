package mysql

import (
	"database/sql"
	"github.com/choonsiong/whisper/pkg/models"
)

// Define a WhisperModel type which wraps a sql.DB connection pool.
type WhisperModel struct {
	DB *sql.DB
}

// Insert a new whisper into the database.
func (m *WhisperModel) Insert(title, content, expires string) (int, error) {
	// SQL statement we want to execute.
	stmt := `INSERT INTO whispers (title, content, created, expires) VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() on the embedded connection pool to execute the statement. This method
	// returns a sql.Result object, which contains some basic information about what happened when
	// the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}


	// Get the ID of newly inserted record in the whispers table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// Return a specific whisper based on its id.
func (m *WhisperModel) Get(id int) (*models.Whisper, error) {
	return nil, nil
}

// Return the 10 most recently created whispers.
func (m *WhisperModel) Latest() ([]*models.Whisper, error) {
	return nil, nil
}