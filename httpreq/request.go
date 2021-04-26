package httpreq

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// ReqPOST jsonのPOSTリクエストを送ります
func ReqPOST(URL string, i interface{}) (*http.Response, error) {
	jsonT, _ := json.Marshal(i)
	res, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonT))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return res, nil
}
