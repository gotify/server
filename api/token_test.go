package api

import (
	"math/rand"
	"net/http/httptest"
	"testing"

	"strings"

	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gotify/server/mode"
	"github.com/gotify/server/model"
	"github.com/gotify/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	firstApplicationToken  = "APorrUa5b1IIK3y"
	secondApplicationToken = "AKo_Pp6ww_9vZal"
	firstClientToken       = "CPorrUa5b1IIK3y"
	secondClientToken      = "CKo_Pp6ww_9vZal"
)

func TestTokenSuite(t *testing.T) {
	suite.Run(t, new(TokenSuite))
}

type TokenSuite struct {
	suite.Suite
	db       *test.Database
	a        *TokenAPI
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
	notified bool
}

func (s *TokenSuite) BeforeTest(suiteName, testName string) {
	mode.Set(mode.TestDev)
	rand.Seed(50)
	s.recorder = httptest.NewRecorder()
	s.db = test.NewDB(s.T())
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	withURL(s.ctx, "http", "example.com")
	s.notified = false
	s.a = &TokenAPI{DB: s.db, NotifyDeleted: s.notify}
}

func (s *TokenSuite) notify(uint, string) {
	s.notified = true
}

func (s *TokenSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

// test application api

func (s *TokenSuite) Test_CreateApplication_mapAllParameters() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name&description=description_text")
	s.a.CreateApplication(s.ctx)

	expected := &model.Application{ID: 1, Token: firstApplicationToken, UserID: 5, Name: "custom_name", Description: "description_text"}
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Equal(s.T(), expected, s.db.GetApplicationByID(1))
}

func (s *TokenSuite) Test_CreateApplication_expectBadRequestOnEmptyName() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=&description=description_text")
	s.a.CreateApplication(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.Empty(s.T(), s.db.GetApplicationsByUser(5))
}

func (s *TokenSuite) Test_DeleteApplication_expectNotFoundOnCurrentUserIsNotOwner() {
	s.db.User(2)
	s.db.User(5).App(5)

	test.WithUser(s.ctx, 2)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/5", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "5"}}

	s.a.DeleteApplication(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
	s.db.AssertAppExist(5)
}

func (s *TokenSuite) Test_CreateApplication_onlyRequiredParameters() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")
	s.a.CreateApplication(s.ctx)

	expected := &model.Application{ID: 1, Token: firstApplicationToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.db.GetApplicationsByUser(5), expected)
}
func (s *TokenSuite) Test_ensureApplicationHasCorrectJsonRepresentation() {
	actual := &model.Application{ID: 1, UserID: 2, Token: "Aasdasfgeeg", Name: "myapp", Description: "mydesc", Image: "asd"}
	test.JSONEquals(s.T(), actual, `{"id":1,"token":"Aasdasfgeeg","name":"myapp","description":"mydesc", "image": "asd"}`)
}

func (s *TokenSuite) Test_CreateApplication_returnsApplicationWithID() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateApplication(s.ctx)

	expected := &model.Application{ID: 1, Token: firstApplicationToken, Name: "custom_name", Image: "http://example.com/static/defaultapp.png", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), expected, s.recorder)
}

func (s *TokenSuite) Test_CreateApplication_withExistingToken() {
	s.db.User(5)
	s.db.User(6).AppWithToken(1, firstApplicationToken)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateApplication(s.ctx)

	expected := &model.Application{ID: 2, Token: secondApplicationToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.db.GetApplicationsByUser(5), expected)
}

func (s *TokenSuite) Test_GetApplications() {
	userBuilder := s.db.User(5)
	first := userBuilder.NewAppWithToken(1, "perfper")
	second := userBuilder.NewAppWithToken(2, "asdasd")

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.a.GetApplications(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	first.Image = "http://example.com/static/defaultapp.png"
	second.Image = "http://example.com/static/defaultapp.png"
	test.BodyEquals(s.T(), []*model.Application{first, second}, s.recorder)
}

func (s *TokenSuite) Test_GetApplications_WithImage() {
	userBuilder := s.db.User(5)
	first := userBuilder.NewAppWithToken(1, "perfper")
	second := userBuilder.NewAppWithToken(2, "asdasd")
	first.Image = "abcd.jpg"
	s.db.UpdateApplication(first)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.a.GetApplications(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	first.Image = "http://example.com/image/abcd.jpg"
	second.Image = "http://example.com/static/defaultapp.png"
	test.BodyEquals(s.T(), []*model.Application{first, second}, s.recorder)
}

func (s *TokenSuite) Test_DeleteApplication_expectNotFound() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstApplicationToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "4"}}

	s.a.DeleteApplication(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *TokenSuite) Test_DeleteApplication() {
	s.db.User(5).App(1)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstApplicationToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.a.DeleteApplication(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertAppNotExist(1)
}

// test client api

func (s *TokenSuite) Test_ensureClientHasCorrectJsonRepresentation() {
	actual := &model.Client{ID: 1, UserID: 2, Token: "Casdasfgeeg", Name: "myclient"}
	test.JSONEquals(s.T(), actual, `{"id":1,"token":"Casdasfgeeg","name":"myclient"}`)
}

func (s *TokenSuite) Test_CreateClient_mapAllParameters() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name&description=description_text")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 1, Token: firstClientToken, UserID: 5, Name: "custom_name"}
	assert.Equal(s.T(), 200, s.recorder.Code)
	assert.Contains(s.T(), s.db.GetClientsByUser(5), expected)
}

func (s *TokenSuite) Test_CreateClient_expectBadRequestOnEmptyName() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=&description=description_text")

	s.a.CreateClient(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.Empty(s.T(), s.db.GetClientsByUser(5))
}

func (s *TokenSuite) Test_DeleteClient_expectNotFoundOnCurrentUserIsNotOwner() {
	s.db.User(5).Client(7)
	s.db.User(2)

	test.WithUser(s.ctx, 2)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/7", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "7"}}

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
	s.db.AssertClientExist(7)
}

func (s *TokenSuite) Test_CreateClient_returnsClientWithID() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 1, Token: firstClientToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), expected, s.recorder)
}

func (s *TokenSuite) Test_CreateClient_withExistingToken() {
	s.db.User(5).ClientWithToken(1, firstClientToken)

	test.WithUser(s.ctx, 5)
	s.withFormData("name=custom_name")

	s.a.CreateClient(s.ctx)

	expected := &model.Client{ID: 2, Token: secondClientToken, Name: "custom_name", UserID: 5}
	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), expected, s.recorder)
}

func (s *TokenSuite) Test_GetClients() {
	userBuilder := s.db.User(5)
	first := userBuilder.NewClientWithToken(1, "perfper")
	second := userBuilder.NewClientWithToken(2, "asdasd")

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("GET", "/tokens", nil)

	s.a.GetClients(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	test.BodyEquals(s.T(), []*model.Client{first, second}, s.recorder)
}

func (s *TokenSuite) Test_DeleteClient_expectNotFound() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstClientToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "8"}}

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

//
func (s *TokenSuite) Test_DeleteClient() {
	s.db.User(5).Client(8)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("DELETE", "/token/"+firstClientToken, nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "8"}}

	assert.False(s.T(), s.notified)

	s.a.DeleteClient(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	s.db.AssertClientNotExist(8)
	assert.True(s.T(), s.notified)
}

func (s *TokenSuite) Test_UploadAppImage_NoImageProvided_expectBadRequest() {
	s.db.User(5).App(1)
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.Close()
	s.ctx.Request = httptest.NewRequest("POST", "/irrelevant", &b)
	s.ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())

	test.WithUser(s.ctx, 5)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.a.UploadApplicationImage(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.Equal(s.T(), s.ctx.Errors[0].Err, errors.New("file with key 'file' must be present"))
}

func (s *TokenSuite) Test_UploadAppImage_OtherErrors_expectServerError() {
	s.db.User(5).App(1)
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	defer writer.Close()
	s.ctx.Request = httptest.NewRequest("POST", "/irrelevant", &b)
	s.ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())

	test.WithUser(s.ctx, 5)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.a.UploadApplicationImage(s.ctx)

	assert.Equal(s.T(), 500, s.recorder.Code)
	assert.Equal(s.T(), s.ctx.Errors[0].Err, errors.New("multipart: NextPart: EOF"))
}

func (s *TokenSuite) Test_UploadAppImage_WithImageFile_expectSuccess() {
	s.db.User(5).App(1)

	cType, buffer, err := upload(map[string]*os.File{"file": mustOpen("../test/assets/image.png")})
	assert.Nil(s.T(), err)
	s.ctx.Request = httptest.NewRequest("POST", "/irrelevant", &buffer)
	s.ctx.Request.Header.Set("Content-Type", cType)
	test.WithUser(s.ctx, 5)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.a.UploadApplicationImage(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	_, err = os.Stat("PorrUa5b1IIK3yKo_Pp6ww_9v.png")
	assert.Nil(s.T(), err)

	s.a.DeleteApplication(s.ctx)

	_, err = os.Stat("PorrUa5b1IIK3yKo_Pp6ww_9v.png")
	assert.True(s.T(), os.IsNotExist(err))
}

func (s *TokenSuite) Test_UploadAppImage_WithImageFile_DeleteExstingImageAndGenerateNewName() {
	s.db.User(5)
	s.db.CreateApplication(&model.Application{UserID: 5, ID: 1, Image: "PorrUa5b1IIK3yKo_Pp6ww_9v.png"})

	cType, buffer, err := upload(map[string]*os.File{"file": mustOpen("../test/assets/image.png")})
	assert.Nil(s.T(), err)
	s.ctx.Request = httptest.NewRequest("POST", "/irrelevant", &buffer)
	s.ctx.Request.Header.Set("Content-Type", cType)
	test.WithUser(s.ctx, 5)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}
	fakeImage(s.T(), "PorrUa5b1IIK3yKo_Pp6ww_9v.png")

	s.a.UploadApplicationImage(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	_, err = os.Stat("PorrUa5b1IIK3yKo_Pp6ww_9v.png")
	assert.True(s.T(), os.IsNotExist(err))
	_, err = os.Stat("Zal6-ySIuL-T3EMLCcFtityHn.png")
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), os.Remove("Zal6-ySIuL-T3EMLCcFtityHn.png"))
}

func (s *TokenSuite) Test_UploadAppImage_WithImageFile_DeleteExistingImage() {
	s.db.User(5)
	s.db.CreateApplication(&model.Application{UserID: 5, ID: 1, Image: "existing.png"})

	fakeImage(s.T(), "existing.png")
	cType, buffer, err := upload(map[string]*os.File{"file": mustOpen("../test/assets/image.png")})
	assert.Nil(s.T(), err)
	s.ctx.Request = httptest.NewRequest("POST", "/irrelevant", &buffer)
	s.ctx.Request.Header.Set("Content-Type", cType)
	test.WithUser(s.ctx, 5)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.a.UploadApplicationImage(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	_, err = os.Stat("existing.png")
	assert.True(s.T(), os.IsNotExist(err))

	os.Remove("PorrUa5b1IIK3yKo_Pp6ww_9v.png")
}

func (s *TokenSuite) Test_UploadAppImage_WithTextFile_expectBadRequest() {
	s.db.User(5).App(1)

	cType, buffer, err := upload(map[string]*os.File{"file": mustOpen("../test/assets/text.txt")})
	assert.Nil(s.T(), err)
	s.ctx.Request = httptest.NewRequest("POST", "/irrelevant", &buffer)
	s.ctx.Request.Header.Set("Content-Type", cType)
	test.WithUser(s.ctx, 5)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.a.UploadApplicationImage(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
	assert.Equal(s.T(), s.ctx.Errors[0].Err, errors.New("file must be an image"))
}

func (s *TokenSuite) Test_UploadAppImage_expectNotFound() {
	s.db.User(5)

	test.WithUser(s.ctx, 5)
	s.ctx.Request = httptest.NewRequest("POST", "/irrelevant", nil)
	s.ctx.Params = gin.Params{{Key: "id", Value: "4"}}

	s.a.UploadApplicationImage(s.ctx)

	assert.Equal(s.T(), 404, s.recorder.Code)
}

func (s *TokenSuite) Test_UploadAppImage_WithSaveError_expectServerError() {
	s.db.User(5).App(1)

	cType, buffer, err := upload(map[string]*os.File{"file": mustOpen("../test/assets/image.png")})
	assert.Nil(s.T(), err)
	s.ctx.Request = httptest.NewRequest("POST", "/irrelevant/", &buffer)
	s.a.ImageDir = "asdasd/asdasda/asdasd"
	s.ctx.Request.Header.Set("Content-Type", cType)
	test.WithUser(s.ctx, 5)
	s.ctx.Params = gin.Params{{Key: "id", Value: "1"}}

	s.a.UploadApplicationImage(s.ctx)

	assert.Equal(s.T(), 500, s.recorder.Code)
}

func (s *TokenSuite) withFormData(formData string) {
	s.ctx.Request = httptest.NewRequest("POST", "/token", strings.NewReader(formData))
	s.ctx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}

func withURL(ctx *gin.Context, scheme, host string) {
	ctx.Set("location", &url.URL{Scheme: scheme, Host: host})
}

// A modified version of https://stackoverflow.com/a/20397167/4244993 from Attila O.
func upload(values map[string]*os.File) (contentType string, buffer bytes.Buffer, err error) {
	w := multipart.NewWriter(&buffer)
	for key, r := range values {
		var fw io.Writer
		if fw, err = w.CreateFormFile(key, r.Name()); err != nil {
			return
		}

		if _, err = io.Copy(fw, r); err != nil {
			return
		}
	}
	contentType = w.FormDataContentType()
	w.Close()
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func fakeImage(t *testing.T, path string) {
	data, err := ioutil.ReadFile("../test/assets/image.png")
	assert.Nil(t, err)
	// Write data to dst
	err = ioutil.WriteFile(path, data, 0644)
	assert.Nil(t, err)
}
