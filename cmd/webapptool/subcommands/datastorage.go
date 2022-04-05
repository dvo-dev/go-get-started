package subcommands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dvo-dev/go-get-started/pkg/utils/requests"
	"github.com/spf13/cobra"
)

var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Root cmd for datastorage operations",
	Long:  "This is the root cmd for datastorage operations:\n\tretrieve\n\tupload\n\tdelete",
}

var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Command to retrieve datastorage data with a given name",
	Long:  "This is a data subcommand to access data in webapp's datastorage with a given name",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("this command requires the name of the data to access")
			return
		}

		params := map[string]string{"name": string(args[0])}
		resp, err := requests.GetRequest("http://0.0.0.0:8080/datastorage", &params, nil)
		if err != nil {
			fmt.Printf("failed to GET /datastorage: %v\n", err)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("failed to read response: %v\n", err)
			return
		}

		var JSON map[string]any
		err = json.Unmarshal([]byte(body), &JSON)
		fmt.Printf("server response:\n\t%+v\n", JSON)
	},
}

func init() {
	dataCmd.AddCommand(retrieveCmd)
	RootCmd.AddCommand(dataCmd)
}
