package cngn

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/sodiqscript111/cNGN-Go/envelope"
	cngnErr "github.com/sodiqscript111/cNGN-Go/error"
	"github.com/sodiqscript111/cNGN-Go/utils"
)

const BaseURL = "https://api.cngn.co/v1/api"

type Client struct {
	authToken     string
	baseURL       string
	http          *http.Client
	encryptionKey string
	privateKey    string
	hasSecurity   bool
}

func New(authToken string) *Client {
	return WithBaseURL(authToken, BaseURL)
}

func WithBaseURL(authToken, baseURL string) *Client {
	return &Client{
		authToken: authToken,
		baseURL:   strings.TrimRight(baseURL, "/"),
		http:      &http.Client{},
	}
}

func (c *Client) WithSecurity(encryptionKey, privateKey string) *Client {
	c.encryptionKey = encryptionKey
	c.privateKey = privateKey
	c.hasSecurity = true
	return c
}

func FromEnv() (*Client, *cngnErr.Error) {
	if err := godotenv.Load(); err != nil {
		return nil, cngnErr.NewConfigurationError("failed to load .env: " + err.Error())
	}

	authToken := os.Getenv("CNGN_KEY")
	if authToken == "" {
		return nil, cngnErr.NewConfigurationError("CNGN_KEY is not set")
	}

	encryptionKey := os.Getenv("CNGN_ENCRYPTION_KEY")
	if encryptionKey == "" {
		return nil, cngnErr.NewConfigurationError("CNGN_ENCRYPTION_KEY is not set")
	}

	privateKey := os.Getenv("CNGN_PRIVATE_KEY")
	if privateKey == "" {
		privateKey = os.Getenv("CNGN_SSH_PRIVATE_KEY")
	}
	if privateKey == "" {
		return nil, cngnErr.NewConfigurationError("CNGN_PRIVATE_KEY or CNGN_SSH_PRIVATE_KEY is not set")
	}

	return New(authToken).WithSecurity(encryptionKey, privateKey), nil
}

func (c *Client) Send(method, path string, query map[string]string, body any, result any) *cngnErr.Error {
	reqURL := c.baseURL + path

	reqBody, err := c.prepareRequestBody(body)
	if err != nil {
		return err
	}

	req, reqErr := http.NewRequest(method, reqURL, bytes.NewReader(reqBody))
	if reqErr != nil {
		return cngnErr.NewNetworkError(reqErr)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, httpErr := c.http.Do(req)
	if httpErr != nil {
		return cngnErr.NewNetworkError(httpErr)
	}
	defer resp.Body.Close()

	respBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return cngnErr.NewNetworkError(readErr)
	}

	return c.processResponse(resp.StatusCode, respBytes, result)
}

func (c *Client) prepareRequestBody(body any) ([]byte, *cngnErr.Error) {
	if body == nil {
		return nil, nil
	}

	reqBody, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		return nil, cngnErr.NewParseError(jsonErr)
	}

	if c.hasSecurity && c.encryptionKey != "" {
		encrypted, cryptErr := utils.AESEncrypt(string(reqBody), c.encryptionKey)
		if cryptErr != nil {
			return nil, cryptErr
		}
		reqBody, jsonErr = json.Marshal(encrypted)
		if jsonErr != nil {
			return nil, cngnErr.NewParseError(jsonErr)
		}
	}

	return reqBody, nil
}

func (c *Client) processResponse(statusCode int, respBytes []byte, result any) *cngnErr.Error {
	if statusCode < 200 || statusCode >= 300 {
		var errEnv envelope.ErrorEnvelope
		msg := string(respBytes)
		if json.Unmarshal(respBytes, &errEnv) == nil {
			msg = errEnv.Message
		}
		kind, field, apiMsg := cngnErr.ClassifyApiError(uint16(statusCode), msg)
		return cngnErr.NewApiError(uint16(statusCode), kind, field, apiMsg)
	}

	responseValue := make(map[string]any)
	if jsonErr := json.Unmarshal(respBytes, &responseValue); jsonErr != nil {
		return cngnErr.NewParseError(jsonErr)
	}

	if c.hasSecurity && c.privateKey != "" {
		if dataStr, ok := responseValue["data"].(string); ok {
			decrypted, decryptErr := utils.Ed25519DecryptWithPrivateKey(c.privateKey, dataStr)
			if decryptErr != nil {
				return decryptErr
			}
			var decryptedData map[string]any
			if json.Unmarshal([]byte(decrypted), &decryptedData) == nil {
				if innerData, ok := decryptedData["data"]; ok {
					responseValue["data"] = innerData
				} else {
					var rawData any
					if json.Unmarshal([]byte(decrypted), &rawData) == nil {
						responseValue["data"] = rawData
					}
				}
			}
		}
	}

	if result != nil {
		reJSON, jsonErr := json.Marshal(responseValue)
		if jsonErr != nil {
			return cngnErr.NewParseError(jsonErr)
		}
		if jsonErr := json.Unmarshal(reJSON, result); jsonErr != nil {
			return cngnErr.NewParseError(jsonErr)
		}
	}

	return nil
}

func SendRequest[T any](client *Client, method, path string, query map[string]string, body any) (*envelope.Response[T], *cngnErr.Error) {
	result := &envelope.Response[T]{}
	err := client.Send(method, path, query, body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func sendJSON[T any](client *Client, method, path string, query map[string]string, body any) (T, *cngnErr.Error) {
	var result T
	err := client.Send(method, path, query, body, &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

func fmtUint32(n uint32) string {
	return strconv.FormatUint(uint64(n), 10)
}
