// +build linux darwin
// +build !race

package plugin

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/gotify/server/auth"
	"github.com/gotify/server/model"
	"github.com/gotify/server/plugin/compat"
	"github.com/gotify/server/plugin/testing/mock"
	"github.com/gotify/server/test"
	"github.com/gotify/server/test/testdb"

	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
)

const examplePluginPath = "github.com/gotify/server/plugin/example/echo"
const mockPluginPath = mock.ModulePath
const danglingPluginPath = "github.com/gotify/server/plugin/testing/removed"

type ManagerSuite struct {
	suite.Suite
	db          *testdb.Database
	manager     *Manager
	e           *gin.Engine
	g           *gin.RouterGroup
	msgReceiver chan MessageWithUserID

	tmpDir test.TmpDir
}

func (s *ManagerSuite) Notify(uid uint, message *model.MessageExternal) {
	s.msgReceiver <- MessageWithUserID{
		Message: *message,
		UserID:  uid,
	}
}

func (s *ManagerSuite) SetupSuite() {
	s.tmpDir = test.NewTmpDir("gotify_managersuite")

	test.WithWd(path.Join(test.GetProjectDir(), "./plugin/example/echo"), func(origWd string) {
		exec.Command("go", "get", "-d").Run()
		goBuildFlags := []string{"build", "-buildmode=plugin", "-o=" + s.tmpDir.Path("echo.so")}

		for _, extraFlag := range extraGoBuildFlags {
			goBuildFlags = append(goBuildFlags, extraFlag)
		}

		cmd := exec.Command("go", goBuildFlags...)
		cmd.Stderr = os.Stderr
		assert.Nil(s.T(), cmd.Run())
	})

	s.db = testdb.NewDBWithDefaultUser(s.T())
	s.makeDanglingPluginConf(1)

	e := gin.New()
	manager, err := NewManager(s.db.GormDatabase, s.tmpDir.Path(), e.Group("/plugin/:id/custom/"), s)
	s.e = e
	assert.Nil(s.T(), err)

	p := new(mock.Plugin)
	assert.Nil(s.T(), manager.LoadPlugin(p))
	assert.Nil(s.T(), manager.initializeSingleUserPlugin(compat.UserContext{
		ID:    1,
		Admin: true,
	}, p))

	s.manager = manager
	s.msgReceiver = make(chan MessageWithUserID)

	assert.Contains(s.T(), s.manager.plugins, examplePluginPath)
	if pluginConf, err := s.db.GetPluginConfByUserAndPath(1, examplePluginPath); assert.NoError(s.T(), err) {
		assert.NotNil(s.T(), pluginConf)
	}
}

func (s *ManagerSuite) TearDownSuite() {
	assert.Nil(s.T(), s.tmpDir.Clean())
}

func (s *ManagerSuite) getConfForExamplePlugin(uid uint) *model.PluginConf {
	pluginConf, err := s.db.GetPluginConfByUserAndPath(uid, examplePluginPath)
	assert.NoError(s.T(), err)
	return pluginConf

}

func (s *ManagerSuite) getConfForMockPlugin(uid uint) *model.PluginConf {
	pluginConf, err := s.db.GetPluginConfByUserAndPath(uid, mockPluginPath)
	assert.NoError(s.T(), err)
	return pluginConf
}

func (s *ManagerSuite) getMockPluginInstance(uid uint) *mock.PluginInstance {
	pid := s.getConfForMockPlugin(uid).ID
	return s.manager.instances[pid].(*mock.PluginInstance)
}

func (s *ManagerSuite) makeDanglingPluginConf(uid uint) *model.PluginConf {
	conf := &model.PluginConf{
		UserID:     uid,
		ModulePath: danglingPluginPath,
		Token:      auth.GeneratePluginToken(),
		Enabled:    true,
	}
	s.db.CreatePluginConf(conf)
	return conf
}

func (s *ManagerSuite) TestWebhook_blockedIfDisabled() {
	conf := s.getConfForExamplePlugin(1)
	t := httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/custom/%s/echo", conf.ID, conf.Token), nil)

	r := httptest.NewRecorder()
	s.e.ServeHTTP(r, t)

	assert.Equal(s.T(), 400, r.Code)
}

func (s *ManagerSuite) TestWebhook_successIfEnabled() {
	conf := s.getConfForExamplePlugin(1)

	assert.Nil(s.T(), s.manager.SetPluginEnabled(conf.ID, true))
	defer func() { assert.Nil(s.T(), s.manager.SetPluginEnabled(conf.ID, false)) }()
	assert.True(s.T(), s.getConfForExamplePlugin(1).Enabled)

	t := httptest.NewRequest("GET", fmt.Sprintf("/plugin/%d/custom/%s/echo", conf.ID, conf.Token), nil)

	r := httptest.NewRecorder()
	s.e.ServeHTTP(r, t)

	assert.Equal(s.T(), 200, r.Code)
}

func (s *ManagerSuite) TestInitializePlugin_noOpIfEmpty() {
	assert.Nil(s.T(), s.manager.loadPlugins(""))
}
func (s *ManagerSuite) TestInitializePlugin_directoryInvalid_expectError() {
	assert.Error(s.T(), s.manager.loadPlugins("<<"))
}

func (s *ManagerSuite) TestInitializePlugin_invalidPlugin_expectError() {
	assert.Error(s.T(), s.manager.loadPlugins(test.GetProjectDir()))
}

func (s *ManagerSuite) TestInitializePlugin_brokenPlugin_expectError() {
	tmpDir := test.NewTmpDir("gotify_testbrokenplugin")
	defer tmpDir.Clean()
	test.WithWd(path.Join(test.GetProjectDir(), "./plugin/testing/broken/nothing"), func(origWd string) {
		exec.Command("go", "get", "-d").Run()
		goBuildFlags := []string{"build", "-buildmode=plugin", "-o=" + tmpDir.Path("empty.so")}

		for _, extraFlag := range extraGoBuildFlags {
			goBuildFlags = append(goBuildFlags, extraFlag)
		}

		cmd := exec.Command("go", goBuildFlags...)
		cmd.Stderr = os.Stderr
		assert.Nil(s.T(), cmd.Run())
	})
	assert.Error(s.T(), s.manager.loadPlugins(tmpDir.Path()))
}

func (s *ManagerSuite) TestInitializePlugin_alreadyLoaded_expectError() {
	assert.Error(s.T(), s.manager.loadPlugins(s.tmpDir.Path()))
}

func (s *ManagerSuite) TestInitializePlugin_alreadyEnabledInConf_expectAutoEnable() {
	s.db.User(2)
	s.db.CreatePluginConf(&model.PluginConf{
		UserID:     2,
		ModulePath: mockPluginPath,
		Token:      "P1234",
		Enabled:    true,
	})

	assert.Nil(s.T(), s.manager.InitializeForUserID(2))
	inst := s.getMockPluginInstance(2)
	assert.True(s.T(), inst.Enabled)

}

func (s *ManagerSuite) TestInitializePlugin_alreadyEnabledInConf_failedToLoadConfig_disableAutomatically() {
	s.db.User(3)
	s.db.CreatePluginConf(&model.PluginConf{
		UserID:     3,
		ModulePath: mockPluginPath,
		Token:      "Ptttt",
		Enabled:    true,
		Config:     []byte(`invalid: """`),
	})

	assert.Nil(s.T(), s.manager.InitializeForUserID(3))
	inst := s.getMockPluginInstance(3)
	assert.False(s.T(), inst.Enabled)

}

func (s *ManagerSuite) TestInitializePlugin_alreadyEnabled_cannotEnable_disabledAutomatically() {
	s.db.NewUserWithName(4, "enable_fail_2")
	mock.ReturnErrorOnEnableForUser(4, errors.New("test error"))
	s.db.CreatePluginConf(&model.PluginConf{
		UserID:     4,
		ModulePath: mockPluginPath,
		Token:      "P5478",
		Enabled:    true,
	})

	assert.Nil(s.T(), s.manager.InitializeForUserID(4))
	inst := s.getMockPluginInstance(4)
	assert.False(s.T(), inst.Enabled)
	assert.False(s.T(), s.getConfForMockPlugin(4).Enabled)
}

func (s *ManagerSuite) TestInitializePlugin_userIDNotExist_expectError() {
	assert.Error(s.T(), s.manager.InitializeForUserID(99))
}

func (s *ManagerSuite) TestSetPluginEnabled() {
	pid := s.getConfForMockPlugin(1).ID
	assert.Nil(s.T(), s.manager.SetPluginEnabled(pid, true))
	assert.Error(s.T(), s.manager.SetPluginEnabled(pid, true))
	assert.Nil(s.T(), s.manager.SetPluginEnabled(pid, false))
}

func (s *ManagerSuite) TestSetPluginEnabled_EnableReturnsError_cannotEnable() {
	s.db.NewUserWithName(5, "enable_fail")
	errExpected := errors.New("test error")
	mock.ReturnErrorOnEnableForUser(5, errExpected)

	assert.Nil(s.T(), s.manager.InitializeForUserID(5))

	pid := s.getConfForMockPlugin(5).ID
	assert.Error(s.T(), s.manager.SetPluginEnabled(pid, false))
	assert.EqualError(s.T(), s.manager.SetPluginEnabled(pid, true), errExpected.Error())

	assert.False(s.T(), s.getConfForMockPlugin(5).Enabled)
}

func (s *ManagerSuite) TestSetPluginEnabled_DisableReturnsError_cannotDisable() {
	s.db.NewUserWithName(6, "disable_fail")
	errExpected := errors.New("test error")
	mock.ReturnErrorOnDisableForUser(6, errExpected)

	assert.Nil(s.T(), s.manager.InitializeForUserID(6))

	pid := s.getConfForMockPlugin(6).ID
	assert.Nil(s.T(), s.manager.SetPluginEnabled(pid, true))
	assert.EqualError(s.T(), s.manager.SetPluginEnabled(pid, false), errExpected.Error())

	assert.True(s.T(), s.getConfForMockPlugin(6).Enabled)
}

func (s *ManagerSuite) TestAddRemoveNewUser() {
	s.db.User(7)
	s.makeDanglingPluginConf(7)

	assert.Nil(s.T(), s.manager.InitializeForUserID(7))
	pid := s.getConfForExamplePlugin(7).ID
	assert.True(s.T(), s.manager.HasInstance(pid))

	assert.Nil(s.T(), s.manager.SetPluginEnabled(s.getConfForMockPlugin(7).ID, true))

	assert.Nil(s.T(), s.manager.RemoveUser(7))
	assert.False(s.T(), s.manager.HasInstance(pid))
}

func (s *ManagerSuite) TestRemoveUser_DisableFail_cannotRemove() {
	s.manager.initializeForUser(*s.db.NewUserWithName(8, "disable_fail_2"))
	errExpected := errors.New("test error")
	mock.ReturnErrorOnDisableForUser(8, errExpected)
	s.manager.SetPluginEnabled(s.getConfForMockPlugin(8).ID, true)

	assert.EqualError(s.T(), s.manager.RemoveUser(8), errExpected.Error())
}

func (s *ManagerSuite) TestRemoveUser_danglingConf_expectSuccess() {
	// make a dangling conf for this instance
	s.db.User(9)
	s.db.CreatePluginConf(&model.PluginConf{
		ModulePath: mockPluginPath,
		Enabled:    true,
		UserID:     9,
		Token:      auth.GenerateNotExistingToken(auth.GeneratePluginToken, s.manager.pluginConfExists),
	})
	s.db.CreatePluginConf(&model.PluginConf{
		ModulePath: examplePluginPath,
		Enabled:    true,
		UserID:     9,
		Token:      auth.GenerateNotExistingToken(auth.GeneratePluginToken, s.manager.pluginConfExists),
	})
	assert.Nil(s.T(), s.manager.RemoveUser(9))
}

func (s *ManagerSuite) TestTriggerMessage() {
	inst := s.getMockPluginInstance(1)
	inst.TriggerMessage()
	select {
	case msg := <-s.msgReceiver:
		assert.Equal(s.T(), uint(1), msg.UserID)
		assert.NotEmpty(s.T(), msg.Message.Extras)
	case <-time.After(1 * time.Second):
		assert.Fail(s.T(), "read message time out")
	}
}

func (s *ManagerSuite) TestStorage() {
	inst := s.getMockPluginInstance(1)

	assert.Nil(s.T(), inst.SetStorage([]byte("test")))
	storage, err := inst.GetStorage()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test", string(storage))
}

func (s *ManagerSuite) TestGetPluginInfo() {
	assert.Equal(s.T(), mock.Name, s.manager.PluginInfo(mock.ModulePath).Name)
}

func (s *ManagerSuite) TestGetPluginInfo_notFound_doNotPanic() {
	assert.NotPanics(s.T(), func() {
		s.manager.PluginInfo("not/exist")
	})
}

func (s *ManagerSuite) TestSetPluginEnabled_expectNotFound() {
	assert.Error(s.T(), s.manager.SetPluginEnabled(99, true))
}

func TestManagerSuite(t *testing.T) {
	suite.Run(t, new(ManagerSuite))
}

func TestNewManager_CannotLoadDirectory_expectError(t *testing.T) {
	_, err := NewManager(nil, "<>", nil, nil)
	assert.Error(t, err)
}

func TestNewManager_NonPluginFile_expectError(t *testing.T) {
	_, err := NewManager(nil, path.Join(test.GetProjectDir(), "test/assets/"), nil, nil)
	assert.Error(t, err)
}

func TestNewManager_FaultyDB_expectError(t *testing.T) {
	tmpDir := test.NewTmpDir("gotify_testnewmanager_faultydb")
	defer tmpDir.Clean()
	for _, suite := range []struct {
		pkg         string
		faultyTable string
		name        string
	}{{"plugin/example/minimal/", "plugin_confs", "minimal"}, {"plugin/example/clock/", "applications", "clock"}} {
		test.WithWd(path.Join(test.GetProjectDir(), suite.pkg), func(origWd string) {
			exec.Command("go", "get", "-d").Run()
			goBuildFlags := []string{"build", "-buildmode=plugin", "-o=" + tmpDir.Path(fmt.Sprintf("%s.so", suite.name))}

			for _, extraFlag := range extraGoBuildFlags {
				goBuildFlags = append(goBuildFlags, extraFlag)
			}

			cmd := exec.Command("go", goBuildFlags...)
			cmd.Stderr = os.Stderr
			assert.Nil(t, cmd.Run())
		})
		db := testdb.NewDBWithDefaultUser(t)
		db.GormDatabase.DB.Callback().Create().Register("no_create", func(s *gorm.Scope) {
			if s.TableName() == suite.faultyTable {
				s.Err(errors.New("database failed"))
			}
		})
		_, err := NewManager(db, tmpDir.Path(), nil, nil)
		assert.Error(t, err)
		os.Remove(tmpDir.Path(fmt.Sprintf("%s.so", suite.name)))
	}
}

func TestNewManager_InternalApplicationManagement(t *testing.T) {
	db := testdb.NewDBWithDefaultUser(t)

	{
		// Application exist, no plugin conf
		db.CreateApplication(&model.Application{
			Token:    "Ainternal_obsolete",
			Internal: true,
			Name:     "obsolete plugin application",
			UserID:   1,
		})

		if app, err := db.GetApplicationByToken("Ainternal_obsolete"); assert.NoError(t, err) {
			assert.True(t, app.Internal)
		}
		_, err := NewManager(db, "", nil, nil)
		assert.Nil(t, err)
		if app, err := db.GetApplicationByToken("Ainternal_obsolete"); assert.NoError(t, err) {
			assert.False(t, app.Internal)
		}
	}
	{
		// Application exist, conf exist, no compat
		assert.NoError(t, db.CreateApplication(&model.Application{
			Token:    "Ainternal_not_loaded",
			Internal: true,
			Name:     "not loaded plugin application",
			UserID:   1,
		}))
		if app, err := db.GetApplicationByToken("Ainternal_not_loaded"); assert.NoError(t, err) {
			assert.NoError(t, db.CreatePluginConf(&model.PluginConf{
				ApplicationID: app.ID,
				UserID:        1,
				Enabled:       true,
				Token:         auth.GeneratePluginToken(),
			}))
		}

		if app, err := db.GetApplicationByToken("Ainternal_not_loaded"); assert.NoError(t, err) {
			assert.True(t, app.Internal)
		}
		_, err := NewManager(db, "", nil, nil)
		assert.Nil(t, err)
		if app, err := db.GetApplicationByToken("Ainternal_not_loaded"); assert.NoError(t, err) {
			assert.False(t, app.Internal)
		}
	}
	{
		// Application exist, conf exist, has compat
		assert.NoError(t, db.CreateApplication(&model.Application{
			Token:    "Ainternal_loaded",
			Internal: false,
			Name:     "not loaded plugin application",
			UserID:   1,
		}))
		if app, err := db.GetApplicationByToken("Ainternal_loaded"); assert.NoError(t, err) {
			assert.NoError(t, db.CreatePluginConf(&model.PluginConf{
				ApplicationID: app.ID,
				UserID:        1,
				Enabled:       true,
				ModulePath:    mock.ModulePath,
				Token:         auth.GeneratePluginToken(),
			}))
		}

		if app, err := db.GetApplicationByToken("Ainternal_loaded"); assert.NoError(t, err) {
			assert.False(t, app.Internal)
		}
		manager, err := NewManager(db, "", nil, nil)
		assert.Nil(t, err)
		assert.Nil(t, manager.LoadPlugin(new(mock.Plugin)))
		assert.Nil(t, manager.InitializeForUserID(1))
		if app, err := db.GetApplicationByToken("Ainternal_loaded"); assert.NoError(t, err) {
			assert.True(t, app.Internal)
		}
	}
}

func TestPluginFileLoadError(t *testing.T) {
	err := pluginFileLoadError{Filename: "test.so", UnderlyingError: errors.New("test error")}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test.so")
	assert.Contains(t, err.Error(), "test error")
}
