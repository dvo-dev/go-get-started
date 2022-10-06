package subcommands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dvo-dev/go-get-started/utils/requests"
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
		if len(args) < 1 {
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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("failed to read response: %v\n", err)
			return
		}

		var JSON map[string]any
		err = json.Unmarshal([]byte(body), &JSON)
		if err != nil {
			fmt.Printf("failed to read response: %v\n", err)
			return
		}
		fmt.Printf("server response:\n\t%+v\n", JSON)
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Command to upload data (string only at this time)",
	Long:  "This is a data subcommand to upload data (string only at this time) to the webapp's datastorage with a given name",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("this command requires 2 arguments:\n\t1. Name of data\n\t2. Data value (string)")
			return
		}

		params := map[string]string{"name": string(args[0])}
		data := map[string][]byte{"data": []byte(args[1])}
		resp, err := requests.PostRequest(
			"http://0.0.0.0:8080/datastorage",
			"multipart/form-data",
			&params,
			&data,
			nil,
		)

		if err != nil {
			fmt.Printf("failed to POST to /datastorage: %v\n", err)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("failed to read response: %v\n", err)
			return
		}

		var JSON map[string]any
		err = json.Unmarshal([]byte(body), &JSON)
		if err != nil {
			fmt.Printf("failed to read response: %v\n", err)
			return
		}
		fmt.Printf("server response:\n\t%+v\n", JSON)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Command to delete datastorage data with a given name",
	Long:  "This is a data subcommand to delete data in webapp's datastorage with a given name",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("this command requires the name of the data to delete")
			return
		}

		params := map[string]string{"name": string(args[0])}
		resp, err := requests.CustomRequest(
			"http://0.0.0.0:8080/datastorage",
			http.MethodDelete,
			&params,
			nil,
		)
		if err != nil {
			fmt.Printf("failed to DELETE from /datastorage: %v\n", err)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("failed to read response: %v\n", err)
			return
		}

		var JSON map[string]any
		err = json.Unmarshal([]byte(body), &JSON)
		if err != nil {
			fmt.Printf("failed to read response: %v\n", err)
			return
		}
		fmt.Printf("server response:\n\t%+v\n", JSON)
	},
}

func init() {
	dataCmd.AddCommand(retrieveCmd)
	dataCmd.AddCommand(uploadCmd)
	dataCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(dataCmd)
}
