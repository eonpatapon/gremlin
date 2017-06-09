package gremlin

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Request struct {
	client    *GremlinClient `json:"-"`
	RequestId string         `json:"requestId"`
	Op        string         `json:"op"`
	Processor string         `json:"processor"`
	Args      *RequestArgs   `json:"args"`
}

type RequestArgs struct {
	Gremlin           string            `json:"gremlin,omitempty"`
	Session           string            `json:"session,omitempty"`
	Bindings          Bind              `json:"bindings,omitempty"`
	Language          string            `json:"language,omitempty"`
	Rebindings        Bind              `json:"rebindings,omitempty"`
	Sasl              []byte            `json:"sasl,omitempty"`
	BatchSize         int               `json:"batchSize,omitempty"`
	ManageTransaction bool              `json:"manageTransaction,omitempty"`
	Aliases           map[string]string `json:"aliases,omitempty"`
}

type Bind map[string]interface{}

func (req *Request) Bindings(bindings Bind) *Request {
	req.Args.Bindings = bindings
	return req
}

func (req *Request) ManageTransaction(flag bool) *Request {
	req.Args.ManageTransaction = flag
	return req
}

func (req *Request) Aliases(aliases map[string]string) *Request {
	req.Args.Aliases = aliases
	return req
}

func (req *Request) Session(session string) *Request {
	req.Args.Session = session
	return req
}

func (req *Request) SetProcessor(processor string) *Request {
	req.Processor = processor
	return req
}

func (req *Request) Exec() (data []byte, err error) {
	// Prepare the Data
	message, err := json.Marshal(req)
	if err != nil {
		return
	}
	// Prepare the request message
	var requestMessage []byte
	mimeType := []byte("application/json")
	mimeTypeLen := byte(len(mimeType))
	requestMessage = append(requestMessage, mimeTypeLen)
	requestMessage = append(requestMessage, mimeType...)
	requestMessage = append(requestMessage, message...)

	if err = req.client.wsConn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return
	}
	if err = req.client.wsConn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return
	}
	if err = req.client.wsConn.WriteMessage(websocket.BinaryMessage, requestMessage); err != nil {
		return
	}

	return req.ReadResponse()
}
