package commands

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/theiskaa/nt/assets"
	"github.com/theiskaa/nt/lib/services"
	"github.com/theiskaa/nt/pkg"
)

var pushCommand = &cobra.Command{
	Use:   "push",
	Short: "Pushes all nodes from [X] service to [Y] service(in case, if nodes doesn't exists in [Y] service)",
	Run:   runPushCommand,
}

func initPushCommand() {
	appCommand.AddCommand(pushCommand)
}

func runPushCommand(cmd *cobra.Command, args []string) {
	determineService()
	loading.Start()

	availableServices := []string{}
	// Generate a list of available services
	// by not including current service.
	for _, s := range services.Services {
		if service.Type() == s {
			continue
		}

		availableServices = append(availableServices, s)
	}

	loading.Stop()

	// Ask for service selection.
	var selected string
	survey.AskOne(
		assets.ChooseRemotePrompt(availableServices),
		&selected,
	)
	if len(selected) == 0 {
		os.Exit(-1)
		return
	}

	selectedService := serviceFromType(selected, true)

	loading.Start()
	pushedNodes, errs := service.Push(selectedService)
	loading.Stop()

	if len(pushedNodes) == 0 && len(errs) == 0 {
		pkg.Print("Everything up-to-date", color.FgHiGreen)
		return
	}

	pkg.PrintErrors("push", errs)
	pkg.Alert(pkg.SuccessL, fmt.Sprintf("Pushed %v nodes", len(pushedNodes)))
}
