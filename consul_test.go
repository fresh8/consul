package consul

import (
	"errors"
	"testing"
	"time"

	"github.com/fresh8/consul/mock_api"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestTagServiceHostPort(t *testing.T) {
	st := serviceTag{Service: "servicename", Tag: "tagname"}
	t.Run("no consul", func(t *testing.T) {
		_, err := TagServiceHostPort(st.Service, st.Tag)
		assert.Error(t, err, "lookup should fail as there is no catalog")
		assert.NotNil(t, serviceCache[st], "the cache should now have an entry for this service")
		assert.WithinDuration(t, time.Now(), *serviceCache[st].LastError, 1*time.Minute)
	})
	ctrl := gomock.NewController(t)
	mCatalog := mock_api.NewMockCatalog(ctrl)
	ConsulCatalog = mCatalog

	t.Run("no service", func(t *testing.T) {
		mCatalog.EXPECT().Service(st.Service, st.Tag, nil).Return([]*api.CatalogService{}, nil, nil).Times(1)
		_, err := TagServiceHostPort(st.Service, st.Tag)
		assert.Error(t, ErrNoService, err, "ErrNoService error should be returned if consul returns 0 services")
	})

	t.Run("no service address and cannot find node", func(t *testing.T) {
		mCatalog.EXPECT().Service("service1", "", nil).Return([]*api.CatalogService{
			&api.CatalogService{
				Node: "node1",
			},
		}, nil, nil).Times(1)
		mCatalog.EXPECT().Node("node1", nil).Return(nil, nil, errors.New("node get error")).Times(1)
		_, err := TagServiceHostPort("service1", "")
		assert.Error(t, err, "an error should be returned if there is no service address and node cannot be found")
	})
	t.Run("use node address instead of service address", func(t *testing.T) {
		mCatalog.EXPECT().Service("service2", "", nil).Return([]*api.CatalogService{
			&api.CatalogService{
				Node:        "node1",
				ServicePort: 80,
			},
		}, nil, nil).Times(1)
		mCatalog.EXPECT().Node("node1", nil).Return(&api.CatalogNode{
			Node: &api.Node{
				Address: "nodeAddress",
			},
		}, nil, nil).Times(1)
		hostPort, err := TagServiceHostPort("service2", "")
		assert.NoError(t, err)
		assert.Equal(t, "nodeAddress:80", hostPort, "the node address should be concatenated with service port")
	})
}
