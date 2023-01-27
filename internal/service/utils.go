package service

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func sendRequest(method, uri string, reader io.Reader, respStruct interface{}, headers *map[string]string) error {
	uri += "/client/quiz/withdraw-balance"
	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	// req.Header.Set("Content-Type")
	if headers != nil {
		for s, v := range *headers {
			req.Header.Set(s, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		switch respStruct.(type) {
		case io.Reader:
			respStruct = resp.Body
			return nil
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if respStruct != nil {
			err = json.Unmarshal(body, &respStruct)
			if err != nil {
				return err
			}
		}

	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	return nil

}
