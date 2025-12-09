package authcmd

import (
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func NewAuthCmd(use string, short string) *cobra.Command{
	return &cobra.Command{
		Use: use,
		Short: short,
		RunE: Run,
	}
}

func Run(cmd *cobra.Command, args []string) error{
    browser.OpenURL("http://localhost:3000/test")
	return nil
}