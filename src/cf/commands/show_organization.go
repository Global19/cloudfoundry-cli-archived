package commands

import (
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/codegangsta/cli"
	"strings"
)

type ShowOrganization struct {
	ui     terminal.UI
	orgReq requirements.OrganizationRequirement
}

func NewShowOrganization(ui terminal.UI) (cmd *ShowOrganization) {
	cmd = new(ShowOrganization)
	cmd.ui = ui
	return
}

func (cmd *ShowOrganization) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "org")
		return
	}

	cmd.orgReq = reqFactory.NewOrganizationRequirement(c.Args()[0])
	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		cmd.orgReq,
	}

	return
}

func (cmd *ShowOrganization) Run(c *cli.Context) {
	org := cmd.orgReq.GetOrganization()
	cmd.ui.Say("Getting info for org %s", org.Name)
	cmd.ui.Ok()
	cmd.ui.Say("%s:", terminal.EntityNameColor(org.Name))

	domains := []string{}
	for _, domain := range org.Domains {
		domains = append(domains, domain.Name)
	}

	spaces := []string{}
	for _, space := range org.Spaces {
		spaces = append(spaces, space.Name)
	}

	cmd.ui.Say("  domains: %s", terminal.EntityNameColor(strings.Join(domains, ", ")))
	cmd.ui.Say("  spaces: %s", terminal.EntityNameColor(strings.Join(spaces, ", ")))
}
