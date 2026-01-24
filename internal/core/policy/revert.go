package policy

import (
	"encoding/json"
	"fmt"
	"net/http"
	httpclient "raind/internal/core/client"
)

func NewServicePolicyRevert() *ServicePolicyRevert {
	return &ServicePolicyRevert{}
}

type ServicePolicyRevert struct{}

func (s *ServicePolicyRevert) Revert() error {
	httpClient := httpclient.NewHttpClient()
	httpClient.NewRequest(
		http.MethodPost,
		"/v1/policies/revert",
		nil,
	)
	resp, err := httpClient.Client.Do(httpClient.Request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respModel RevertResponseModel

	if !httpClient.IsStatusOk(resp) {
		decodeErr := json.NewDecoder(resp.Body).Decode(&respModel)
		if decodeErr != nil {
			return fmt.Errorf("decode response: %w", decodeErr)
		}
		return fmt.Errorf("unexpected status: %s: %s", resp.Status, respModel.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(&respModel); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	fmt.Println("policy revert success")

	return nil
}
