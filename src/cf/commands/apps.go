package commands

import (
	"cf/api"
	"cf/requirements"
	"cf/terminal"
	"fmt"
	"github.com/codegangsta/cli"
	"strings"
)

type Apps struct {
	ui        terminal.UI
	spaceRepo api.SpaceRepository
}

func NewApps(ui terminal.UI, spaceRepo api.SpaceRepository) (a Apps) {
	a.ui = ui
	a.spaceRepo = spaceRepo
	return
}

func (a Apps) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		reqFactory.NewTargetedSpaceRequirement(),
	}
	return
}

func (a Apps) Run(c *cli.Context) {
	a.ui.Say("Getting applications in %s...", a.spaceRepo.GetCurrentSpace().Name)

	space, err := a.spaceRepo.GetSummary()

	if err != nil {
		a.ui.Failed(err.Error())
		return
	}

	apps := space.Applications

	a.ui.Ok()

	table := [][]string{
		[]string{"name", "status", "usage", "urls"},
	}

	for _, app := range apps {
		table = append(table, []string{
			app.Name,
			app.State,
			fmt.Sprintf("%d x %s", app.Instances, byteSize(app.Memory*MEGABYTE)),
			strings.Join(app.Urls, ", "),
		})
	}

	a.ui.DisplayTable(table, a.coloringFunc)
}

func (a Apps) coloringFunc(value string, row int, col int) string {
	if row > 0 && col == 1 {
		return coloredState(value)
	}

	return terminal.DefaultColoringFunc(value, row, col)
}
