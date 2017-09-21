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
	t.Run("standard lookup", func(t *testing.T) {
		mCatalog.EXPECT().Service("service3", "", nil).Return([]*api.CatalogService{
			&api.CatalogService{
				ServiceAddress: "serviceAddress",
				ServicePort:    80,
			},
		}, nil, nil).Times(1)
		hostPort, err := TagServiceHostPort("service3", "")
		assert.NoError(t, err)
		assert.Equal(t, "serviceAddress:80", hostPort, "service address should be used with service port")
	})
	t.Run("success followed by failure", func(t *testing.T) {
		mCatalog.EXPECT().Service("service4", "", nil).Return([]*api.CatalogService{
			&api.CatalogService{
				ServiceAddress: "serviceAddress",
				ServicePort:    80,
			},
		}, nil, nil).Times(1)
		hostPort, err := TagServiceHostPort("service4", "")
		assert.NoError(t, err)
		assert.Equal(t, "serviceAddress:80", hostPort, "service address should be used with service port")

		mCatalog.EXPECT().Service("service4", "", nil).Return(nil, nil, errors.New("random error")).Times(1)

		hostPort, err = TagServiceHostPort("service4", "")
		assert.NoError(t, err, "there should not be an error as there was a recent successful lookup")
		assert.Equal(t, "serviceAddress:80", hostPort, "service address should be used with service port")
	})
	t.Run("host port changes after first lookup", func(t *testing.T) {
		mCatalog.EXPECT().Service("service5", "", nil).Return([]*api.CatalogService{
			&api.CatalogService{
				ServiceAddress: "serviceAddress",
				ServicePort:    80,
			},
		}, nil, nil).Times(1)
		hostPort, err := TagServiceHostPort("service5", "")
		assert.NoError(t, err)
		assert.Equal(t, "serviceAddress:80", hostPort, "service address should be used with service port")
		mCatalog.EXPECT().Service("service5", "", nil).Return([]*api.CatalogService{
			&api.CatalogService{
				ServiceAddress: "serviceAddress2",
				ServicePort:    80,
			},
		}, nil, nil).Times(1)
		hostPort, err = TagServiceHostPort("service5", "")

		assert.NoError(t, err, "there should not be an error as there was a recent successful lookup")
		assert.Equal(t, "serviceAddress2:80", hostPort, "service address should be used with service port")
	})

	t.Run("multi host port", func(t *testing.T) {
		cservices := []*api.CatalogService{
			&api.CatalogService{
				ServiceAddress: "serviceAddress1",
				ServicePort:    80,
			},
			&api.CatalogService{
				ServiceAddress: "serviceAddress2",
				ServicePort:    80,
			},
		}
		expectedHostPorts := []string{
			"serviceAddress1:80",
			"serviceAddress2:80",
		}
		mCatalog.EXPECT().Service("service10", "", nil).Return(cservices, nil, nil).Times(1)

		hostPorts, err := TagServiceHostPortMulti("service10", "")
		assert.NoError(t, err)
		assert.Equal(t, expectedHostPorts, hostPorts)
	})
}
