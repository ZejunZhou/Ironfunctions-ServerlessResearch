package datastore

import (
	"testing"

	"github.com/ZejunZhou/Ironfunctions-ServerlessResearch/api/datastore/internal/datastoretest"
)

func TestDatastore(t *testing.T) {
	datastoretest.Test(t, NewMock())
}
