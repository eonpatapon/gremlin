package gremlin

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type GremlinClient struct {
	Connected bool
	Server    *url.URL
	servers   []*url.URL
	conn      net.Conn
	wsConn    *websocket.Conn
}

func (c *GremlinClient) Connect() error {
	var err error
	for _, s := range c.servers {
		conn, err := net.DialTimeout("tcp", s.Host, 1*time.Second)
		if err != nil {
			continue
		}
		c.Connected = true
		c.conn = conn
		c.Server = s
		break
	}
	if !c.Connected {
		return errors.New("Could not establish connection. Please check your connection string and ensure at least one server is up.")
	}

	d := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
	}
	c.wsConn, _, err = d.Dial(c.Server.String(), http.Header{})
	if err != nil {
		return fmt.Errorf("Failed to connect with websocket: %s", err)
	}
	return nil
}

func (c *GremlinClient) Query(query string) *Request {
	args := &RequestArgs{
		Gremlin:  query,
		Language: "gremlin-groovy",
	}

	return &Request{
		client:    c,
		RequestId: uuid.NewV4().String(),
		Op:        "eval",
		Processor: "",
		Args:      args,
	}
}

func (c *GremlinClient) Close() (err error) {
	err = c.wsConn.Close()
	if err != nil {
		return
	}
	return c.conn.Close()
}

func NewGremlinClient(s ...string) (c GremlinClient, err error) {
	// If no arguments use environment variable
	if len(s) == 0 {
		connString := strings.TrimSpace(os.Getenv("GREMLIN_SERVERS"))
		if connString == "" {
			err = errors.New("No servers set. Configure servers to connect to using the GREMLIN_SERVERS environment variable.")
			return
		}
		c.servers, err = SplitServers(connString)
	} else {
		// Else use the supplied servers
		for _, v := range s {
			var u *url.URL
			if u, err = url.Parse(v); err != nil {
				return
			}
			c.servers = append(c.servers, u)
		}
	}
	return
}

func SplitServers(connString string) (servers []*url.URL, err error) {
	serverStrings := strings.Split(connString, ",")
	if len(serverStrings) < 1 {
		return nil, errors.New("Connection string is not in expected format. An example of the expected format is 'ws://server1:8182, ws://server2:8182'.")
	}
	for _, serverString := range serverStrings {
		var u *url.URL
		if u, err = url.Parse(strings.TrimSpace(serverString)); err != nil {
			return
		}
		servers = append(servers, u)
	}
	return
}
