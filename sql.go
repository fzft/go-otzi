package go_otzi

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
)

type SQLRequest struct {
	ID     int
	Action string
	DB     string
	Table  string
	Cols   []string
}

type Result struct {
	Columns map[string][]string
	// ... other fields
}

type SQLService struct {
	Trie  *Trie
	Queue []SQLRequest
}

func NewSQLService() *SQLService {
	return &SQLService{
		Trie:  NewTrie(),
		Queue: []SQLRequest{},
	}
}

func (s *SQLService) AddRequest(req SQLRequest) {
	s.Queue = append(s.Queue, req)
	s.Trie.Insert(req)
}

func (s *SQLService) ExecuteMergedQuery() map[int]Result {
	//mergedQuery := s.Trie.Merge()
	//return FanOutWithMock(mergedQuery, s.Queue)
	return nil
}

func (s *SQLService) constructQueryFromRequest(req SQLRequest) string {
	// Construct the SQL query from the given request
	// This is a simplified version and should be enhanced
	return "SELECT ... FROM ..."
}

func CreateMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	return db, mock, nil
}

func FanOutWithMock(mergedQuery string, reqs []SQLRequest) map[int]Result {
	db, mock, err := CreateMock()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Mock the expected query and its result
	columns := make([]string, 0)
	for _, req := range reqs {
		columns = append(columns, req.Cols...)
	}

	// Set up expected result from mock DB
	rows := sqlmock.NewRows(columns)
	// rows.AddRow("value1", "value2", ...)

	mock.ExpectQuery(mergedQuery).WillReturnRows(rows)

	// Query the mock database
	result, _ := QueryDB(db, mergedQuery)

	// For simplicity, let's just map every request ID to the same result
	responseMap := make(map[int]Result)
	for i := range reqs {
		responseMap[i] = result
	}

	return responseMap
}

func QueryDB(db *sql.DB, query string) (Result, error) {
	// ...
	return Result{}, nil
}
