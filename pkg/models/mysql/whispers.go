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
	stmt := `SELECT id, title, content, created, expires FROM whispers WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// This returns a sql.Rows resultset containing the result of our query.
	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is always properly closed before the Latest() method
	// returns. This defer statement should come after you check for an error from the Query() method. Otherwise,
	// if Query() returns an error, you'll get a panic trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the models.Whispers objects.
	whispers := []*models.Whisper{}

	// Use rows.Next() to iterate through the rows in the resultset. This prepares the first (and then each
	// subsequent) row to be acted on by the rows.Scan() method. If iteration over all the rows completes then
	// the resultset automatically closes itself and frees-up the underlying database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed whisper struct.
		s := &models.Whisper{}

		// Use rows.Scan() to copy the values from each field in the row to the new whisper object
		// that we created. Again, the arguments to row.Scan() must be pointers to the place you want
		// to copy the data into, and the number of arguments must be exactly the same as the number
		// of columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}

		// Append it to the slice of whispers.
		whispers = append(whispers, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any error
	// that was encountered during the iteration. It's important to call this - don't assume
	// that a successful iteration was completed over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went ok then returns the whisper slice.
	return whispers, nil
}