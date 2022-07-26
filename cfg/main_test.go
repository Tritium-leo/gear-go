package cfg

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

func TestLoadDirs(t *testing.T) {
	c, err := LoadDirs("./example/example.yaml", "./example/example.json", "./example/example.ini")
	assert.Equal(t, err, nil)
	fmt.Printf("%#v", c)
}

func TestConfig_Get(t *testing.T) {
	c, _ := LoadDirs("./example/example.yaml")
	v, err := c.Get("mysql.Address")
	assert.Equal(t, err, nil)
	assert.Equal(t, v.(string), "localhost")
}

func TestConfig_GetBool(t *testing.T) {
	c, _ := LoadDirs("./example/example.yaml")
	v, err := c.GetBool("mysql.LogMode")
	assert.Equal(t, err, nil)
	assert.Equal(t, v, true)

}

func TestConfig_GetInt(t *testing.T) {
	c, _ := LoadDirs("./example/example.yaml")
	v, err := c.GetInt("mysql.Port")
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 3306)
}

func TestConfig_GetFloat64(t *testing.T) {
	c, _ := LoadDirs("./example/example.yaml")
	v, err := c.GetFloat64("test.float_test")
	assert.Equal(t, err, nil)
	assert.Equal(t, v, 3.5)
}

func TestConfig_GetString(t *testing.T) {
	c, _ := LoadDirs("./example/example.yaml")
	v, err := c.GetString("mysql.Username")
	assert.Equal(t, err, nil)
	assert.Equal(t, v, "root")
}

func TestConfig_GetStruct(t *testing.T) {
	c, _ := LoadDirs("./example/example.yaml")
	type MysqlConfig struct {
		Username string
		Password string `yaml:"password"`
		Address  string
		Port     int
		LogMode  bool
	}
	cf := &MysqlConfig{}
	err := c.GetStruct("mysql", cf)
	assert.Equal(t, err, nil)
	assert.Equal(t, cf.Username, "root")
	assert.Equal(t, cf.Password, "root")
	assert.Equal(t, cf.Address, "localhost")
	assert.Equal(t, cf.Port, 3306)
	assert.Equal(t, cf.LogMode, true)
}
