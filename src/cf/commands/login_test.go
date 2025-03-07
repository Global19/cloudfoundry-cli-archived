package commands_test

import (
	"cf/api"
	. "cf/commands"
	"cf/configuration"
	"cf/terminal"
	"github.com/stretchr/testify/assert"
	"testhelpers"
	"testing"
)

func testSuccessfulLogin(t *testing.T, args []string, inputs []string) (ui *testhelpers.FakeUI) {
	configRepo := testhelpers.FakeConfigRepository{}
	configRepo.Delete()
	config, _ := configRepo.Get()

	ui = new(testhelpers.FakeUI)
	ui.Inputs = inputs
	auth := &testhelpers.FakeAuthenticator{
		AccessToken:  "my_access_token",
		RefreshToken: "my_refresh_token",
		ConfigRepo:   configRepo,
	}
	callLogin(
		args,
		ui,
		configRepo,
		&testhelpers.FakeOrgRepository{},
		&testhelpers.FakeSpaceRepository{},
		auth,
	)

	savedConfig := testhelpers.SavedConfiguration

	assert.Contains(t, ui.Outputs[0], config.Target)
	assert.Contains(t, ui.Outputs[2], "OK")

	assert.Equal(t, savedConfig.AccessToken, "my_access_token")
	assert.Equal(t, savedConfig.RefreshToken, "my_refresh_token")
	assert.Equal(t, auth.Email, "user@example.com")
	assert.Equal(t, auth.Password, "password")

	return
}

func TestSuccessfullyLoggingIn(t *testing.T) {
	ui := testSuccessfulLogin(t, []string{}, []string{"user@example.com", "password"})

	assert.Contains(t, ui.PasswordPrompts[0], "Password")
}

func TestSuccessfullyLoggingInWithUsernameAsArgument(t *testing.T) {
	ui := testSuccessfulLogin(t, []string{"user@example.com"}, []string{"password"})

	assert.Contains(t, ui.PasswordPrompts[0], "Password")
}

func TestSuccessfullyLoggingInWithUsernameAndPasswordAsArguments(t *testing.T) {
	testSuccessfulLogin(t, []string{"user@example.com", "password"}, []string{})
}

func TestUnsuccessfullyLoggingIn(t *testing.T) {
	configRepo := testhelpers.FakeConfigRepository{}
	configRepo.Delete()
	config, _ := configRepo.Get()

	ui := new(testhelpers.FakeUI)
	ui.Inputs = []string{
		"foo@example.com",
		"bar",
		"bar",
		"bar",
		"bar",
	}

	callLogin(
		[]string{},
		ui,
		configRepo,
		&testhelpers.FakeOrgRepository{},
		&testhelpers.FakeSpaceRepository{},
		&testhelpers.FakeAuthenticator{AuthError: true, ConfigRepo: configRepo},
	)

	assert.Contains(t, ui.Outputs[0], config.Target)
	assert.Equal(t, ui.Outputs[1], "Authenticating...")
	assert.Equal(t, ui.Outputs[2], "FAILED")
	assert.Equal(t, ui.Outputs[4], "Authenticating...")
	assert.Equal(t, ui.Outputs[5], "FAILED")
	assert.Equal(t, ui.Outputs[7], "Authenticating...")
	assert.Equal(t, ui.Outputs[8], "FAILED")
}

func TestUnsuccessfullyLoggingInWithoutInteractivity(t *testing.T) {
	configRepo := testhelpers.FakeConfigRepository{}
	configRepo.Delete()
	config, _ := configRepo.Get()

	ui := new(testhelpers.FakeUI)

	callLogin(
		[]string{
			"foo@example.com",
			"bar",
		},
		ui,
		configRepo,
		&testhelpers.FakeOrgRepository{},
		&testhelpers.FakeSpaceRepository{},
		&testhelpers.FakeAuthenticator{AuthError: true, ConfigRepo: configRepo},
	)

	assert.Contains(t, ui.Outputs[0], config.Target)
	assert.Equal(t, ui.Outputs[1], "Authenticating...")
	assert.Equal(t, ui.Outputs[2], "FAILED")
	assert.Contains(t, ui.Outputs[3], "Error authenticating")
	assert.Equal(t, len(ui.Outputs), 4)
}

func callLogin(args []string, ui terminal.UI, configRepo configuration.ConfigurationRepository, orgRepo api.OrganizationRepository, spaceRepo api.SpaceRepository, auth api.Authenticator) {
	l := NewLogin(ui, configRepo, orgRepo, spaceRepo, auth)
	l.Run(testhelpers.NewContext("login", args))
}
