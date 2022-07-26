package cfg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
	"gopkg.in/yaml.v2"
)

type Config struct {
	data map[interface{}]interface{}
}

func LoadDirs(dirs ...string) (c *Config, err error) {
	c = &Config{}
	for _, dir := range dirs {
		fileType := c.guessFileType(dir)
		switch fileType {
		case "yaml":
			err = c.loadFromYaml(dir)
		case "ini":
			err = c.loadFromIni(dir)
		case "json":
			err = c.loadFromJson(dir)
		case "unknown":
			err = errors.New("UNKNOWN File Type " + dir)
		}
	}
	return c, err
}

func (c *Config) guessFileType(path string) string {
	s := strings.Split(path, ".")
	ext := s[len(s)-1]
	switch ext {
	case "yaml", "yml":
		return "yaml"
	case "ini":
		return "ini"
	case "json":
		return "json"
	default:
		return "unknown"
	}
}

// load function
func (c *Config) loadFromYaml(path string) (err error) {
	yamlS, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlS, &c.data)
	if err != nil {
		return errors.New("con not parse " + path + " config")
	}
	return nil
}

func (c *Config) loadFromJson(path string) error {
	jsonS, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonS, &c.data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) loadFromIni(path string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	secs := cfg.Sections()
	for _, sec := range secs {
		c.data[sec.Name()] = map[interface{}]interface{}{}
		for k, v := range sec.KeysHash() {
			c.data[sec.Name()][k.(interface{})] = v
		}
	}
	//err = cfg.MapTo(&c.data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Get(path string) (value interface{}, err error) {
	paths := strings.Split(path, ".")
	pre := c.data
	for idx, p := range paths {
		v, ok := pre[p]
		if !ok {
			return nil, errors.New("Path doesn't exist" + p)
		}
		if (idx + 1) == len(paths) {
			return v.(interface{}), nil
		}
		pre = v.(map[interface{}]interface{})
	}
	return pre, nil
}

func (c *Config) GetString(path string) (value string, err error) {
	var v interface{}
	v, err = c.Get(path)
	if err != nil {
		return "", err
	}

	switch value := v.(type) {
	case string:
		return value, nil
	case bool, int64, int, float64, int32, float32:
		return fmt.Sprint(value), nil
	default:
		return "", errors.New("can't get this value as string")
	}
	return
}

func (c *Config) GetInt(path string) (value int, err error) {
	var v interface{}
	v, err = c.Get(path)
	if err != nil {
		return 0, err
	}
	switch value := v.(type) {
	case string:
		i, err := strconv.Atoi(value)
		return i, err
	case int:
		return value, nil
	case bool:
		if value {
			return 1, nil
		}
		return 0, nil
	case float64:
		return int(value), nil
	default:
		return 0, nil
	}
}

func (c *Config) GetBool(path string) (value bool, err error) {
	var v interface{}
	v, err = c.Get(path)
	if err != nil {
		return false, err
	}
	switch value := v.(type) {
	case bool:
		return value, nil
	case int:
		if value != 0 {
			return true, nil
		}
		return false, nil
	case float64:
		if value != 0.0 {
			return true, nil
		}
		return false, nil
	default:
		return false, nil
	}
}

func (c *Config) GetFloat64(path string) (value float64, err error) {
	var v interface{}
	v, err = c.Get(path)
	if err != nil {
		return 0.0, err
	}
	switch value := v.(type) {
	case float64:
		return value, nil
	case int:
		return float64(value), nil
	case bool:
		if value {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0.0, nil
	}
}

func (c *Config) GetStruct(path string, s interface{}) error {
	d, err := c.Get(path)
	if err != nil {
		return err
	}
	switch d.(type) {
	case string:
		err = c.SetField(s, path, d)
		if err != nil {
			return err
		}
	case map[interface{}]interface{}:
		c.mapToStruct(d.(map[interface{}]interface{}), s)
	}
	return nil
}

func (c *Config) mapToStruct(m map[interface{}]interface{}, s interface{}) {
	for k, v := range m {
		switch k.(type) {
		case string:
			_ = c.SetField(s, k.(string), v)
		}
	}
}

func (c *Config) SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.Indirect(reflect.ValueOf(obj))
	structFieldValue := structValue.FieldByName(name)
	fmt.Print(structValue.Field(0))
	// isValid
	if !structFieldValue.IsValid() {
		return fmt.Errorf("No Such Field: %s in obj", name)
	}
	// canset
	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	if structFieldType.Kind() == reflect.Struct && val.Kind() == reflect.Map {
		vint := val.Interface()
		switch vint.(type) {
		case map[interface{}]interface{}:
			for k, v := range vint.(map[interface{}]interface{}) {
				_ = c.SetField(structFieldValue.Addr().Interface(), k.(string), v)
			}
		case map[string]interface{}:
			for k, v := range vint.(map[string]interface{}) {
				_ = c.SetField(structFieldValue.Addr().Interface(), k, v)
			}

		}
	} else {
		if structFieldType != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}
		structFieldValue.Set(val)
	}
	return nil
}
