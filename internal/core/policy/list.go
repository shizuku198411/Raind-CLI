package policy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	httpclient "raind/internal/core/client"
	"strconv"
	"text/tabwriter"
)

func NewServicePolicyList() *ServicePolicyList {
	return &ServicePolicyList{}
}

type ServicePolicyList struct{}

func (s *ServicePolicyList) List(param ListRequestModel) error {
	var chainName string
	switch param.ChainName {
	case "ew":
		chainName = "RAIND-EW"
	case "ns-obs":
		chainName = "RAIND-NS-OBS"
	case "ns-enf":
		chainName = "RAIND-NS-ENF"
	}

	// request body

	httpClient := httpclient.NewHttpClient()
	httpClient.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/v1/policies/%s", chainName),
		nil,
	)
	resp, err := httpClient.Client.Do(httpClient.Request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respModel ListResponseModel

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

	s.printPolicyList(param.ChainName, respModel.Data.Mode, respModel.Data.Policies)

	return nil
}

func (s *ServicePolicyList) printPolicyList(group string, mode string, PolicyList []PolicyModel) {
	w := tabwriter.NewWriter(
		os.Stdout,
		0,
		0,
		2,
		' ',
		tabwriter.DiscardEmptyColumns,
	)

	// header
	groupMap := map[string]string{
		"ew":     "EAST-WEST",
		"ns-obs": "NORTH-SOUTH",
		"ns-enf": "NORTH-SOUTH",
	}
	fmt.Printf("POLICY TYPE: %s\n", groupMap[group])
	fmt.Printf("MODE: %s\n\n", mode)
	fmt.Println(`FLAG: "*" - Applied, "+" - Apply next commit, "-" - Remove next commit, N/A - Not applied`)
	fmt.Fprintln(w, "\nFLAG\tPOLICY ID\tSRC CONTAINER\tDST CONTAINER\tPROTOCOL\tDST PORT\tCOMMENT\tREASON")

	flagMap := map[string]string{
		"before_commit":      "[+]",
		"remove_next_commit": "[-]",
		"applied":            "[*]",
		"unresolved":         "[ ]",
	}

	// helper
	parseDport := func(port int) string {
		if port == 0 {
			return "*"
		}
		return strconv.Itoa(port)
	}

	// data
	for _, p := range PolicyList {
		flag := flagMap[p.Status]
		id := p.Id
		src := p.Source.ContainerName
		dst := p.Destination.ContainerName
		protocol := p.Protocol
		dport := parseDport(p.DestPort)
		comment := p.Comment
		reason := p.Reason

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", flag, id, src, dst, protocol, dport, comment, reason)
	}

	w.Flush()
}
