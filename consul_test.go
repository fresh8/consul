package consul

import (
	"testing"
	"time"

	"github.com/fresh8/consul/mock_api"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestTagServiceHostPort(t *testing.T) {
	st := serviceTag{Service: "servicename", Tag: "tagname"}
	{
		_, err := TagServiceHostPort(st.Service, st.Tag)
		assert.Error(t, err, "lookup should fail as there is no catalog")
		assert.NotNil(t, serviceCache[st], "the cache should now have an entry for this service")
		assert.WithinDuration(t, time.Now(), *serviceCache[st].LastError, 1*time.Minute)
	}
	ctrl := gomock.NewController(t)
	mCatalog := mock_api.NewMockCatalog(ctrl)
	ConsulCatalog = mCatalog
	{
		mCatalog.EXPECT().Service(st.Service, st.Tag, nil).Return([]*api.CatalogService{}, nil, nil).Times(1)
		_, err := TagServiceHostPort(st.Service, st.Tag)
		assert.Error(t, ErrNoService, err, "ErrNoService error should be returned if consul returns 0 services")
	}
}
