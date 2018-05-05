package binance

type UserDataStreamResponse struct {
	ListenKey string `json:"listenKey"`
}

// GetUserDataStream makes the get request for a user data stream listen key.
func (c *RestClient) GetUserDataStream() (string, error) {
	httpResponse, err := c.Post("/api/v1/userDataStream", nil)
	if err != nil {
		return "", err
	}

	var response UserDataStreamResponse
	if _, err = c.decodeBody(httpResponse, &response); err != nil {
		return "", err
	}

	return response.ListenKey, nil
}
