package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	// Obtain token -> https://pushe.co/docs/api/#api_get_token
	const token = "YOUR_TOKEN"

	// Android doc -> https://pushe.co/docs/api/

	reqData := map[string]interface{}{
		"app_ids":  []string{"YOUR_APP_ID"},
		"platform": 2,
		"data": map[string]interface{}{
			"title":   "This is a simple push",
			"content": "All of your users will see me",

            // actions -> https://pushe.co/docs/api/#api_action_type_table3
            "action": {
                "action_type": "U",
                "url": "https://pushe.co"
            },

            "buttons": [
                {
                    'btn_action': {
                        'action_data': 'MyActivityName',
                        'action_type': 'T',
                        'market_package_name': '',
                        'url': ''
                    },
                    'btn_content': 'Pushe',
                     // More icons -> https: //pushe.co/docs/api/#api_icon_notificaiton_table2
                    'btn_icon': 'open_in_browser',
                    'btn_order': 0,
                }
            ]
		},
		// additional keywords -> https://pushe.co/docs/api/#api_send_advance_notification
	}

	// Marshal returns the JSON encoding of reqData.
	reqJSON, err := json.Marshal(reqData)

	// check encoded json
	if err != nil {
		fmt.Println("json:", err)
		return
	}

	// create request obj
	request, err := http.NewRequest(
		http.MethodPost,
		"https://api.pushe.co/v2/messaging/notifications/",
		bytes.NewBuffer(reqJSON),
	)

	// check request
	if err != nil {
		fmt.Println("Req error:", err)
		return
	}

	// set header
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Token "+token)

	// send request and get response
	client := http.Client{}
	response, err := client.Do(request)

	// check response
	if err != nil {
		fmt.Println("Resp error:", err)
		return
	}

	defer response.Body.Close()

	// check status_code and response
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(response.Body)
	respContent := buf.String()

	fmt.Println("status code =>", response.StatusCode)
	fmt.Println("response =>", respContent)
	fmt.Println("==========")

	if response.StatusCode == http.StatusCreated {
		fmt.Println("success!")

		var respData map[string]interface{}
		_ = json.Unmarshal(buf.Bytes(), &respData)

		var reportURL string

		// hashed_id just generated for Non-Free plan
		if respData["hashed_id"] != nil {
			reportURL = "https://pushe.co/report?id=" + respData["hashed_id"].(string)
		} else {
			reportURL = "no report url for your plan"
		}

		fmt.Println("report_url:", reportURL)
		fmt.Println("notification id:", int(respData["wrapper_id"].(float64)))
	} else {
		fmt.Println("failed!")
	}
}
