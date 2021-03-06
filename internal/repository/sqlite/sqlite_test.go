package sqlite

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	r "github.com/ibeauregard/url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var m = &r.MappingModel{
	ID:      "1",
	Key:     "asd42",
	LongUrl: "http://foobar.com",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		log.Fatalf("unexpected error %q when opening a stub database connection", err)
	}
	return db, mock
}

func TestNewRepositoryOpenError(t *testing.T) {
	_, err := NewRepository("", func(_ string, _ string) (*sql.DB, error) {
		return nil, errors.New("")
	})
	assert.NotNil(t, err)
}

func TestNewRepositoryPingError(t *testing.T) {
	db, mock := NewMock()
	mock.ExpectPing().WillReturnError(errors.New(""))
	_, err := NewRepository("", func(_ string, _ string) (*sql.DB, error) {
		return db, nil
	})
	assert.NotNil(t, err)
}

func TestNewRepositorySuccess(t *testing.T) {
	db, mock := NewMock()
	mock.ExpectPing()
	_, err := NewRepository("", func(_ string, _ string) (*sql.DB, error) {
		return db, nil
	})
	assert.Nil(t, err)
}

func TestCloseError(t *testing.T) {
	db, mock := NewMock()
	mock.ExpectClose().WillReturnError(errors.New("error"))
	repo := &repository{db}
	assert.Error(t, repo.Close())
}

func TestCloseSuccess(t *testing.T) {
	db, mock := NewMock()
	mock.ExpectClose()
	repo := &repository{db}
	assert.Nil(t, repo.Close())
}

func TestFindByLongUrl(t *testing.T) {
	db, mock := NewMock()
	query := `SELECT id, key, long_url FROM mappings WHERE long_url=?`
	rows := sqlmock.NewRows([]string{"id", "key", "long_url"}).AddRow(m.ID, m.Key, m.LongUrl)
	mock.ExpectQuery(query).WithArgs(m.LongUrl).WillReturnRows(rows)
	mapping, err := (&repository{db}).FindByLongUrl(m.LongUrl)
	assert.NotNil(t, mapping)
	assert.NoError(t, err)
}

func TestFindByLongUrlError(t *testing.T) {
	db, mock := NewMock()
	query := `SELECT id, key, long_url FROM mappings WHERE long_url=?`
	rows := sqlmock.NewRows([]string{"id", "key", "long_url"})
	mock.ExpectQuery(query).WithArgs(m.LongUrl).WillReturnRows(rows)
	mapping, err := (&repository{db}).FindByLongUrl(m.LongUrl)
	assert.Empty(t, mapping)
	assert.Error(t, err)
}

func TestFindByKey(t *testing.T) {
	db, mock := NewMock()
	query := `SELECT id, key, long_url FROM mappings WHERE key=?`
	rows := sqlmock.NewRows([]string{"id", "key", "long_url"}).AddRow(m.ID, m.Key, m.LongUrl)
	mock.ExpectQuery(query).WithArgs(m.LongUrl).WillReturnRows(rows)
	mapping, err := (&repository{db}).FindByKey(m.LongUrl)
	assert.NotNil(t, mapping)
	assert.NoError(t, err)
}

func TestFindByKeyError(t *testing.T) {
	db, mock := NewMock()
	query := `SELECT id, key, long_url FROM mappings WHERE key=?`
	rows := sqlmock.NewRows([]string{"id", "key", "long_url"})
	mock.ExpectQuery(query).WithArgs(m.LongUrl).WillReturnRows(rows)
	mapping, err := (&repository{db}).FindByKey(m.LongUrl)
	assert.Empty(t, mapping)
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {
	db, mock := NewMock()
	query := "INSERT INTO mappings \\(key, long_url\\) VALUES\\(\\?, \\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(m.Key, m.LongUrl).WillReturnResult(sqlmock.NewResult(1, 1))
	err := (&repository{db}).Create(m)
	assert.NoError(t, err)
}

func TestCreateError(t *testing.T) {
	db, mock := NewMock()
	query := "INSERT INTO non_existent_table \\(key, long_url\\) VALUES\\(\\?, \\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(m.Key, m.LongUrl).WillReturnResult(sqlmock.NewResult(1, 0))
	err := (&repository{db}).Create(m)
	assert.Error(t, err)
}

func TestCreateStmtCloseError(t *testing.T) {
	db, mock := NewMock()
	query := "INSERT INTO mappings \\(key, long_url\\) VALUES\\(\\?, \\?\\)"
	prep := mock.ExpectPrepare(query).WillReturnCloseError(errors.New("error"))
	prep.ExpectExec().WithArgs(m.Key, m.LongUrl).WillReturnResult(sqlmock.NewResult(1, 1))
	err := (&repository{db}).Create(m)
	assert.NoError(t, err)
}

func TestGetLastIdReturn0(t *testing.T) {
	db, _ := NewMock()
	repo := &repository{db}
	lastId := repo.GetLastId()
	assert.EqualValues(t, 0, lastId)
}

func TestGetLastIdReturn42(t *testing.T) {
	db, mock := NewMock()
	query := "SELECT seq FROM sqlite_sequence WHERE name='mappings'"
	rows := sqlmock.NewRows([]string{"seq"}).AddRow("42")
	mock.ExpectQuery(query).WillReturnRows(rows)
	lastId := (&repository{db}).GetLastId()
	assert.EqualValues(t, 42, lastId)
}
