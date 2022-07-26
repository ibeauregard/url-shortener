package handling

import (
	"errors"
	"fmt"
	repo "github.com/ibeauregard/url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepo struct {
	Repository
	getLastIdOutput uint
	outputMapping   *repo.MappingModel
	outputError     error
}

func TestNewRepoAdapter(t *testing.T) {
	assert.NotNil(t, NewRepoAdapter(&mockRepo{}))
}

func (m *mockRepo) Close() error {
	return m.outputError
}

func TestCloseError(t *testing.T) {
	adapter := &repoAdapter{&mockRepo{outputError: errors.New("error")}}
	err := adapter.Close()
	assert.NotNil(t, err)
}

func TestCloseSuccess(t *testing.T) {
	adapter := &repoAdapter{&mockRepo{outputError: nil}}
	err := adapter.Close()
	assert.Nil(t, err)
}

func (m *mockRepo) FindByLongUrl(_ string) (*repo.MappingModel, error) {
	return m.outputMapping, m.outputError
}

func TestGetKeyNotFound(t *testing.T) {
	adapter := &repoAdapter{&mockRepo{outputMapping: nil, outputError: errors.New("error")}}
	_, found := adapter.getKey("key_not_found")
	assert.False(t, found)
}

func TestGetKeyFound(t *testing.T) {
	modelKey := "my_key"
	adapter := &repoAdapter{&mockRepo{outputMapping: &repo.MappingModel{
		Key: modelKey,
	}, outputError: nil}}
	longUrl, found := adapter.getKey("http://foobar.com")
	assert.EqualValues(t, modelKey, longUrl)
	assert.True(t, found)
}

func (m *mockRepo) FindByKey(_ string) (*repo.MappingModel, error) {
	return m.FindByLongUrl("my_key")
}

func TestGetLongUrlNotFound(t *testing.T) {
	adapter := &repoAdapter{&mockRepo{outputMapping: nil, outputError: errors.New("error")}}
	_, found := adapter.getLongUrl("key_not_found")
	assert.False(t, found)
}

func TestGetLongUrlFound(t *testing.T) {
	modelLongUrl := "http://foobar.com"
	adapter := &repoAdapter{&mockRepo{outputMapping: &repo.MappingModel{
		LongUrl: modelLongUrl,
	}, outputError: nil}}
	longUrl, found := adapter.getLongUrl("key_found")
	assert.EqualValues(t, modelLongUrl, longUrl)
	assert.True(t, found)
}

func (m *mockRepo) Create(_ *repo.MappingModel) error {
	return m.outputError
}

func TestAddMappingError(t *testing.T) {
	adapter := &repoAdapter{&mockRepo{outputError: errors.New("error")}}
	_, err := adapter.addMapping("http://foobar.com")
	assert.NotNil(t, err)
}

func TestAddMappingSuccess(t *testing.T) {
	adapter := &repoAdapter{&mockRepo{outputError: nil, getLastIdOutput: 42}}
	key, err := adapter.addMapping("http://foobar.com")
	assert.EqualValues(t, "Pwt", key)
	assert.Nil(t, err)
}

func (m *mockRepo) GetLastId() uint {
	return m.getLastIdOutput
}

func TestGetNextDatabaseId(t *testing.T) {
	lastIds := []uint{0, 1, 10}
	for _, last := range lastIds {
		t.Run(fmt.Sprintf("last=%v", last), func(t *testing.T) {
			assert.EqualValues(t, last+1, (&repoAdapter{&mockRepo{
				getLastIdOutput: last,
			}}).getNextDatabaseId())
		})
	}
}
