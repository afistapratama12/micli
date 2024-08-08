package repo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/afistapratama12/micli/src/model"
	"github.com/gorilla/websocket"
)

func Call(r model.ReqData, dest interface{}) error {
	var (
		req *http.Request
		err error
	)

	if r.Body != nil {
		req, err = http.NewRequest(r.Method, r.Url, r.Body)
	} else {
		req, err = http.NewRequest(r.Method, r.Url, nil)
	}

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if r.Header != nil || len(r.Header) > 0 {
		for key, val := range r.Header {
			req.Header.Set(key, val)
		}
	}

	// add req param
	if r.Params != nil || len(r.Params) > 0 {
		req.URL.RawQuery = r.Params.Encode()
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("err client do", err)
		return err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("err read body", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d, body %v", resp.StatusCode, string(b))
	}

	err = json.Unmarshal(b, &dest)
	if err != nil {
		log.Println("err unmarshal", err)
		return err
	}

	return nil
}

// "wss://stream-cloud.binanceru.net/ws/btcusdt@trade/ethusdt@trade/ethbtc@trade"

func Stream(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
