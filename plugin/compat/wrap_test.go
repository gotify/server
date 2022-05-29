//go:build linux || darwin
// +build linux darwin

package compat

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"plugin"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CompatSuite struct {
	suite.Suite

	p      Plugin
	tmpDir test.TmpDir
}

func (s *CompatSuite) SetupSuite() {
	s.tmpDir = test.NewTmpDir("gotify_compatsuite")

	test.WithWd(path.Join(test.GetProjectDir(), "./plugin/example/echo"), func(origWd string) {
		exec.Command("go", "get", "-d").Run()
		goBuildFlags := []string{"build", "-buildmode=plugin", "-o=" + s.tmpDir.Path("echo.so")}

		goBuildFlags = append(goBuildFlags, extraGoBuildFlags...)

		cmd := exec.Command("go", goBuildFlags...)
		cmd.Stderr = os.Stderr
		assert.Nil(s.T(), cmd.Run())
	})

	plugin, err := plugin.Open(s.tmpDir.Path("echo.so"))
	assert.Nil(s.T(), err)
	wrappedPlugin, err := Wrap(plugin)
	assert.Nil(s.T(), err)

	s.p = wrappedPlugin
}

func (s *CompatSuite) TearDownSuite() {
	assert.Nil(s.T(), s.tmpDir.Clean())
}

func (s *CompatSuite) TestGetPluginAPIVersion() {
	assert.Equal(s.T(), "v1", s.p.APIVersion())
}

func (s *CompatSuite) TestGetPluginInfo() {
	info := s.p.PluginInfo()

	assert.Equal(s.T(), examplePluginPath, info.ModulePath)
	assert.True(s.T(), info.String() != "")
}

func (s *CompatSuite) TestInstantiatePlugin() {
	inst := s.p.NewPluginInstance(UserContext{
		ID:   1,
		Name: "test",
	})

	assert.NotNil(s.T(), inst)
}

func (s *CompatSuite) TestGetCapabilities() {
	inst := s.p.NewPluginInstance(UserContext{
		ID:   2,
		Name: "test2",
	})

	c := inst.Supports()

	assert.Contains(s.T(), c, Webhooker)
	assert.Contains(s.T(), c.Strings(), string(Webhooker))
	assert.True(s.T(), HasSupport(inst, Webhooker))
	assert.False(s.T(), HasSupport(inst, "not_exist"))
}

func (s *CompatSuite) TestSetConfig() {
	inst := s.p.NewPluginInstance(UserContext{
		ID:   3,
		Name: "test3",
	})

	defaultConfig := inst.DefaultConfig()
	assert.Nil(s.T(), inst.ValidateAndSetConfig(defaultConfig))
}

func (s *CompatSuite) TestRegisterWebhook() {
	inst := s.p.NewPluginInstance(UserContext{
		ID:   4,
		Name: "test4",
	})

	e := gin.New()
	g := e.Group("/")
	assert.NotPanics(s.T(), func() {
		inst.RegisterWebhook("/plugin/4/custom/Pabcd/", g)
	})
}

func (s *CompatSuite) TestEnableDisable() {
	inst := s.p.NewPluginInstance(UserContext{
		ID:   5,
		Name: "test5",
	})
	assert.Nil(s.T(), inst.Enable())
	assert.Nil(s.T(), inst.Disable())
}

func (s *CompatSuite) TestGetDisplay() {
	inst := s.p.NewPluginInstance(UserContext{
		ID:   6,
		Name: "test6",
	})

	assert.NotEqual(s.T(), "", inst.GetDisplay(nil))
}

func TestCompatSuite(t *testing.T) {
	suite.Run(t, new(CompatSuite))
}

func TestWrapIncompatiblePlugins(t *testing.T) {
	tmpDir := test.NewTmpDir("gotify_testwrapincompatibleplugins")
	defer tmpDir.Clean()
	for i, modulePath := range []string{
		"github.com/gotify/server/v2/plugin/testing/broken/noinstance",
		"github.com/gotify/server/v2/plugin/testing/broken/nothing",
		"github.com/gotify/server/v2/plugin/testing/broken/unknowninfo",
		"github.com/gotify/server/v2/plugin/testing/broken/malformedconstructor",
	} {
		fName := tmpDir.Path(fmt.Sprintf("broken_%d.so", i))
		exec.Command("go", "get", "-d").Run()
		goBuildFlags := []string{"build", "-buildmode=plugin", "-o=" + fName}
		goBuildFlags = append(goBuildFlags, extraGoBuildFlags...)
		goBuildFlags = append(goBuildFlags, modulePath)

		cmd := exec.Command("go", goBuildFlags...)
		cmd.Stderr = os.Stderr
		assert.Nil(t, cmd.Run())

		plugin, err := plugin.Open(fName)
		assert.Nil(t, err)
		_, err = Wrap(plugin)
		assert.Error(t, err)
		os.Remove(fName)
	}
}
