package handling

import (
	"errors"
	"fmt"
	repo "github.com/ibeauregard/url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockRepo struct {
	repo.Repository
	getLastIdOutput uint
	outputMapping   *repo.MappingModel
	outputError     error
}

func TestNewRepoProxy(t *testing.T) {
	assert.NotNil(t, NewRepoProxy(&mockRepo{}))
}

func (m *mockRepo) Close() error {
	return m.outputError
}

func TestCloseError(t *testing.T) {
	proxy := &concreteRepoProxy{&mockRepo{outputError: errors.New("error")}}
	err := proxy.Close()
	assert.NotNil(t, err)
}

func TestCloseSuccess(t *testing.T) {
	proxy := &concreteRepoProxy{&mockRepo{outputError: nil}}
	err := proxy.Close()
	assert.Nil(t, err)
}

func (m *mockRepo) FindByLongUrl(_ string) (*repo.MappingModel, error) {
	return m.outputMapping, m.outputError
}

func TestGetKeyNotFound(t *testing.T) {
	proxy := &concreteRepoProxy{&mockRepo{outputMapping: nil, outputError: errors.New("error")}}
	_, found := proxy.getKey("key_not_found")
	assert.False(t, found)
}

func TestGetKeyFound(t *testing.T) {
	modelKey := "my_key"
	proxy := &concreteRepoProxy{&mockRepo{outputMapping: &repo.MappingModel{
		Key: modelKey,
	}, outputError: nil}}
	longUrl, found := proxy.getKey("http://foobar.com")
	assert.EqualValues(t, modelKey, longUrl)
	assert.True(t, found)
}

func (m *mockRepo) FindByKey(_ string) (*repo.MappingModel, error) {
	return m.FindByLongUrl("my_key")
}

func TestGetLongUrlNotFound(t *testing.T) {
	proxy := &concreteRepoProxy{&mockRepo{outputMapping: nil, outputError: errors.New("error")}}
	_, found := proxy.getLongUrl("key_not_found")
	assert.False(t, found)
}

func TestGetLongUrlFound(t *testing.T) {
	modelLongUrl := "http://foobar.com"
	proxy := &concreteRepoProxy{&mockRepo{outputMapping: &repo.MappingModel{
		LongUrl: modelLongUrl,
	}, outputError: nil}}
	longUrl, found := proxy.getLongUrl("key_found")
	assert.EqualValues(t, modelLongUrl, longUrl)
	assert.True(t, found)
}

func (m *mockRepo) Create(_ *repo.MappingModel) error {
	return m.outputError
}

func TestAddMappingError(t *testing.T) {
	proxy := &concreteRepoProxy{&mockRepo{outputError: errors.New("error")}}
	_, err := proxy.addMapping("http://foobar.com")
	assert.NotNil(t, err)
}

func TestAddMappingSuccess(t *testing.T) {
	proxy := &concreteRepoProxy{&mockRepo{outputError: nil, getLastIdOutput: 42}}
	key, err := proxy.addMapping("http://foobar.com")
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
			assert.EqualValues(t, last+1, (&concreteRepoProxy{&mockRepo{
				getLastIdOutput: last,
			}}).getNextDatabaseId())
		})
	}
}
