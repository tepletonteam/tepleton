package rpc

import (
	"fmt"

	"github.com/spf13/cobra"
	wire "github.com/tepleton/go-wire"

	"github.com/tepleton/tepleton-sdk/client"
)

func statusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Query remote node for status",
		RunE:  checkStatus,
	}
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:46657", "Node to connect to")
	return cmd
}

func checkStatus(cmd *cobra.Command, args []string) error {
	// get the node
	node, err := client.GetNode()
	if err != nil {
		return err
	}
	res, err := node.Status()
	if err != nil {
		return err
	}

	output, err := wire.MarshalJSON(res)
	// output, err := json.MarshalIndent(res, "  ", "")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
