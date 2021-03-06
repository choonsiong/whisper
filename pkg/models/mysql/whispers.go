package mysql

import (
	"database/sql"
	"errors"
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
	stmt := `SELECT id, title, content, created, expires FROM whispers WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use the QueryRow() method on the connection pool to execute our SQL statement,
	// passing in the untrusted id variable as the value for the placeholder parameter.
	// This returns a pointer to a sql.Row object which holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed whisper struct.
	s := &models.Whisper{}

	// Use row.Scan() to copy the values from each field in sql.Row to the corresponding field in the
	// whisper struct. Notice that the arguments to row.Scan() are pointers to the place you want to
	// copy the data into, and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything OK, return the whisper object.
	return s, nil
}

// Return the 10 most recently created whispers.
func (m *WhisperModel) Latest() ([]*models.Whisper, error) {
	return nil, nil
}