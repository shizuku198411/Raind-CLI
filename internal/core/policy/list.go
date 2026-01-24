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
	fmt.Println("FLAG: \"*\" - Applied, \"+\" - Apply next commit, \"-\" - Remove next commit, N/A - Not applied\n")

	var chainName string
	switch param.ChainName {
	case "ew":
		chainName = "RAIND-EW"
		if err := s.requestGetList(chainName); err != nil {
			return err
		}

	case "ns-obs":
		chainName = "RAIND-NS-OBS"
		if err := s.requestGetList(chainName); err != nil {
			return err
		}

	case "ns-enf":
		chainName = "RAIND-NS-ENF"
		if err := s.requestGetList(chainName); err != nil {
			return err
		}

	default:
		if err := s.requestGetList("RAIND-EW"); err != nil {
			return err
		}
		fmt.Println("\n===================\n")
		if err := s.requestGetList("RAIND-NS-OBS"); err != nil {
			return err
		}
		fmt.Println("\n===================\n")
		if err := s.requestGetList("RAIND-NS-ENF"); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServicePolicyList) requestGetList(chainName string) error {
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

	s.printPolicyList(chainName, respModel.Data.Mode, respModel.Data.Policies)

	return nil
}

func (s *ServicePolicyList) printPolicyList(chainName string, mode string, PolicyList []PolicyModel) {
	w := tabwriter.NewWriter(
		os.Stdout,
		0,
		0,
		2,
		' ',
		tabwriter.DiscardEmptyColumns,
	)

	flagMap := map[string]string{
		"before_commit":      "[+]",
		"remove_next_commit": "[-]",
		"applied":            "[*]",
		"unresolved":         "[ ]",
	}

	fmt.Printf("POLICY TYPE: %s\n", chainName)
	fmt.Printf("MODE: %s\n", mode)

	if chainName == "RAIND-EW" {
		fmt.Fprintln(w, "\nFLAG\tPOLICY ID\tSRC CONTAINER\tDST CONTAINER\tPROTOCOL\tDST PORT\tCOMMENT\tREASON")
	} else {
		fmt.Fprintln(w, "\nFLAG\tPOLICY ID\tSRC CONTAINER\tDST ADDR\tPROTOCOL\tDST PORT\tCOMMENT\tREASON")
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
		var dst string
		if chainName == "RAIND-EW" {
			dst = p.Destination.ContainerName
		} else {
			dst = p.Destination.Address
		}
		protocol := p.Protocol
		dport := parseDport(p.DestPort)
		comment := p.Comment
		reason := p.Reason

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", flag, id, src, dst, protocol, dport, comment, reason)
	}

	w.Flush()
}
