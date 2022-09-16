package subcommands

import (
	"fmt"
	"net/http"

	"github.com/dvo-dev/go-get-started/pkg/utils/requests"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Checks the health status of the webapp",
	Long:  "This command makes a GET request to the webapp's `/health` endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: parse URL from args or config
		resp, err := requests.GetRequest("http://0.0.0.0:8080/health", nil, nil)
		if err != nil {
			fmt.Printf("failed to GET /health: %v\n", err)
			return
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Println("webapp server is healthy")
		} else {
			fmt.Printf("webapp server is unhealthy :'( %v\n", resp.StatusCode)
		}
	},
}

func init() {
	RootCmd.AddCommand(healthCmd)
}
