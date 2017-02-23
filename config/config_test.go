package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultConfig(t *testing.T) {
	c := DefaultConfig()
	assert := assert.New(t)
	assert.Equal(c.AWSRegion, "eu-west-1")
	assert.Equal(c.Address, "0.0.0.0")
	assert.Equal(c.Port, 7654)
	assert.Equal(c.CrawlInterval, 600)
	assert.Equal(c.LogFormat, "text")
	assert.Equal(c.LogLevel, "debug")
	assert.Equal(c.Owner, os.Getenv("USER"))

}

func Test_ReadEnvironment(t *testing.T) {
	c := DefaultConfig()
	assert := assert.New(t)

	os.Setenv("MELKOR_ADDRESS", "10.0.0.0")
	os.Setenv("MELKOR_PORT", "9090")
	os.Setenv("MELKOR_LOGLEVEL", "error")
	os.Setenv("MELKOR_LOGFORMAT", "json")
	os.Setenv("MELKOR_CRAWLINTERVAL", "500")
	os.Setenv("MELKOR_AWSREGION", "eu-east-2")
	os.Setenv("MELKOR_OWNER", "the_boss")

	ReadEnvironment(c)

	os.Unsetenv("MELKOR_ADDRESS")
	os.Unsetenv("MELKOR_PORT")
	os.Unsetenv("MELKOR_LOGLEVEL")
	os.Unsetenv("MELKOR_LOGFORMAT")
	os.Unsetenv("MELKOR_CRAWLINTERVAL")
	os.Unsetenv("MELKOR_AWSREGION")
	os.Unsetenv("MELKOR_OWNER")

	assert.Equal(c.AWSRegion, "eu-east-2")
	assert.Equal(c.Address, "10.0.0.0")
	assert.Equal(c.Port, 9090)
	assert.Equal(c.CrawlInterval, 500)
	assert.Equal(c.LogFormat, "json")
	assert.Equal(c.LogLevel, "error")
	assert.Equal(c.Owner, "the_boss")
}

func Test_ReadConfigFile(t *testing.T) {
	c := DefaultConfig()
	wd, _ := os.Getwd()
	assert := assert.New(t)

	ReadConfigFile(c, fmt.Sprintf("%s/config_test.yml", wd))
	assert.Equal(c.AWSRegion, "us-east-1")
	assert.Equal(c.Address, "127.0.0.1")
	assert.Equal(c.Port, 8080)
	assert.Equal(c.CrawlInterval, 3600)
	assert.Equal(c.LogFormat, "json")
	assert.Equal(c.LogLevel, "info")
	assert.Equal(c.Owner, "the_team")
}

func Test_ReadConfigFile_Error(t *testing.T) {
	c := DefaultConfig()
	d := DefaultConfig()

	ReadConfigFile(c, getConfigFilePath())

	assert.Equal(t, c, d)
}

func Test_getConfigFilePath(t *testing.T) {
	fp := getConfigFilePath()
	assert.Empty(t, fp)
}

func Test_Initialize(t *testing.T) {
	c := Initialize()
	d := DefaultConfig()

	assert.Equal(t, c, d)
}
