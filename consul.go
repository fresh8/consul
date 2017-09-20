package consul

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

var (
	// ErrNoCatalog is returned when there is no catalog in Consul client
	ErrNoCatalog = errors.New("no Consul catalog in client")

	// ErrNoService is returned when there are no services registered for that
	// service name and tag
	ErrNoService = errors.New("no service registered under this service name and tag")

	// Consul client
	Client        *api.Client
	ConsulCatalog Catalog

	serviceCache      = make(map[serviceTag]*consulService)
	serviceCacheMutex = sync.RWMutex{}
)

type serviceTag struct {
	Service, Tag string
}

type consulService struct {
	HostPort  string
	LastError *time.Time
}

func Setup() error {
	var err error
	if Client != nil && ConsulCatalog != nil {
		log.Println("Consul already initialised")
		return nil
	}
	Client, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}
	if Client.Catalog() == nil {
		return ErrNoCatalog
	}

	ConsulCatalog = Client.Catalog()

	return nil
}

// ServiceHostPort looks up a service by just service name or tag.service
// from local Consul agent
func ServiceHostPort(service string) (string, error) {
	if ConsulCatalog == nil {
		return "", ErrNoCatalog
	}

	serviceParts := strings.Split(service, ".")
	if len(serviceParts) == 2 {
		return TagServiceHostPort(serviceParts[1], serviceParts[0])
	} else {
		return TagServiceHostPort(service, "")
	}
}

// TagServiceHostPort looks up a service by service name and tag
// from local Consul agent
func TagServiceHostPort(service, tag string) (hostPort string, err error) {
	st := serviceTag{Service: service, Tag: tag}
	defer func() {
		if err != nil {
			serviceCacheMutex.Lock()
			defer serviceCacheMutex.Unlock()
			now := time.Now()

			// There has never been a succesful lookup
			// Set last error time and return error
			if serviceCache[st] == nil {
				serviceCache[st] = &consulService{LastError: &now}
				return
			}

			// It has been more than a minute since the last error
			// If LastError is nil the previous request was a success
			if serviceCache[st].LastError != nil && time.Since(*serviceCache[st].LastError) > 1*time.Minute {
				return
			}

			serviceCache[st].LastError = &now
			hostPort = serviceCache[st].HostPort
			err = nil
		}
	}()
	if ConsulCatalog == nil {
		err = ErrNoCatalog
		return
	}

	cservices, _, err := ConsulCatalog.Service(service, tag, nil)
	if err != nil {
		return
	}

	if len(cservices) == 0 {
		err = ErrNoService
		return
	}
	index := rand.Intn(len(cservices))

	cservice := cservices[index]
	addr := cservice.ServiceAddress
	if addr == "" {
		var node *api.CatalogNode
		node, _, err = ConsulCatalog.Node(cservice.Node, nil)
		if err != nil {
			err = errors.Wrap(err, "service registered on a node that does not exist")
			return
		}
		addr = node.Node.Address
	}

	hostPort = fmt.Sprintf("%s:%d", addr, cservice.ServicePort)

	// Add service to map if it doesn't exist or update host port
	serviceCacheMutex.RLock()
	firstLookup := serviceCache[st] == nil
	hostPortChanged := !firstLookup && serviceCache[st].HostPort != hostPort
	serviceCacheMutex.RUnlock()

	if firstLookup {
		serviceCache[st] = &consulService{HostPort: hostPort}
	} else if hostPortChanged {
		serviceCacheMutex.Lock()
		serviceCache[st].HostPort = hostPort
		serviceCache[st].LastError = nil
		serviceCacheMutex.Unlock()
	}
	return
}
