package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StorageTest struct {
	suite.Suite
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageTest))
}

func (t *StorageTest) Test_MGClientPool() {
	clientPool, err := NewMGClientPool(1)
	t.Assert().NoError(err)

	client := clientPool.Get("test_token", "test_url")
	t.Assert().Equal("test_url", client.URL)

	clientPool.Remove("test_token")
	clientPool.Close()
}

func (t *StorageTest) Test_NegativeCapacity() {
	_, err := NewMGClientPool(-1)
	t.Assert().Equal(ErrNegativeCapacity.Error(), err.Error())
}
