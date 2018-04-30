// The MIT License (MIT)
//
// Copyright (c) 2018 Cranky Kernel
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use, copy,
// modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package binance

import (
	"github.com/gorilla/websocket"
	"strings"
	"fmt"
	"net/http"
)

const WS_STREAM_URL = "wss://stream.binance.com:9443"

type StreamClient struct {
	conn *websocket.Conn
}

func NewStreamClient() *StreamClient {
	client := &StreamClient{}
	return client
}

func (c *StreamClient) Connect(streams ... string) (err error) {
	url := fmt.Sprintf("%s/stream?streams=%s", WS_STREAM_URL,
		strings.Join(streams, "/"))
	var httpResponse *http.Response
	c.conn, httpResponse, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	if httpResponse.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("%s", httpResponse.Status)
	}
	return nil
}

func (c *StreamClient) ConnectSingle(stream string) (err error) {
	url := fmt.Sprintf("%s/ws/%s", WS_STREAM_URL, stream)
	var httpResponse *http.Response
	c.conn, httpResponse, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	if httpResponse.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("%s", httpResponse.Status)
	}
	return nil
}

func (c *StreamClient) Close() {
	c.conn.Close()
}

func (c *StreamClient) Next() (messageType int, body []byte, err error) {
	return c.conn.ReadMessage()
}

// Next reads the next message into a generic map.
func (c *StreamClient) NextJSON() (interface{}, error) {
	var message interface{}
	err := c.conn.ReadJSON(&message)
	return message, err
}
