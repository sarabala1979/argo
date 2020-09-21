package cron

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/argoproj/argo/cmd/argo/commands/client"
	cmdcommon "github.com/argoproj/argo/cmd/argo/commands/common"
	cronworkflowpkg "github.com/argoproj/argo/pkg/apiclient/cronworkflow"
	"github.com/argoproj/argo/workflow/validate"
)

func NewLintCommand() *cobra.Command {
	var (
		strict bool
	)
	var command = &cobra.Command{
		Use:   "lint FILE...",
		Short: "validate files or directories of cron workflow manifests",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, apiClient := cmdcommon.CreateNewAPIClientFunc()
			serviceClient := apiClient.NewCronWorkflowServiceClient()
			lint := func(file string) error {
				wfs, err := validate.ParseCronWorkflowsFromFile(file, strict)
				if err != nil {
					return err
				}
				if len(wfs) == 0 {
					return fmt.Errorf("there was nothing to validate")
				}
				for _, wf := range wfs {
					if wf.Namespace == "" {
						wf.Namespace = client.Namespace()
					}
					_, err := serviceClient.LintCronWorkflow(ctx, &cronworkflowpkg.LintCronWorkflowRequest{Namespace: wf.Namespace, CronWorkflow: &wf})
					if err != nil {
						return err
					}
				}
				fmt.Printf("%s is valid\n", file)
				return nil
			}

			invalid := false
			for _, file := range args {
				stat, err := os.Stat(file)
				if err != nil {
					return err
				}
				if stat.IsDir() {
					err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
						fileExt := filepath.Ext(info.Name())
						switch fileExt {
						case ".yaml", ".yml", ".json":
						default:
							return nil
						}
						if err := lint(path); err != nil {
							log.Errorf("Error in file %s: %s", path, err)
							invalid = true
						}
						return nil
					})
					if err != nil {
						return err
					}
				} else {
					if err := lint(file); err != nil {
						log.Errorf("Error in file %s: %s", file, err)
						invalid = true
					}
				}
			}
			if invalid {
				return fmt.Errorf("Errors encountered in validation")
			}
			fmt.Printf("Cron workflow manifests validated\n")
			return nil
		},
	}
	command.Flags().BoolVar(&strict, "strict", true, "perform strict workflow validation")
	return command
}
