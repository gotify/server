package mode

import "github.com/gin-gonic/gin"

const (
	// Dev for development mode.
	Dev = "dev"
	// Prod for production mode.
	Prod = "prod"
	// TestDev used for tests.
	TestDev = "testdev"
)

var mode = Dev

// Set sets the new mode.
func Set(newMode string) {
	mode = newMode
	updateGinMode()
}

// Get returns the current mode.
func Get() string {
	return mode
}

// IsDev returns true if the current mode is dev mode.
func IsDev() bool {
	return Get() == Dev || Get() == TestDev
}

func updateGinMode() {
	switch Get() {
	case Dev:
		gin.SetMode(gin.DebugMode)
	case TestDev:
		gin.SetMode(gin.TestMode)
	case Prod:
		gin.SetMode(gin.ReleaseMode)
	default:
		panic("unknown mode")
	}
}
