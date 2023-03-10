package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	filepath string
	conflist []map[string]map[string]string
}

func NewConfig(filepath string) *Config {
	c := new(Config)
	c.filepath = filepath
	return c
}

// To obtain corresponding value of the key values
func (c *Config) GetValue(section, name string) (string, bool) {
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				return value[name], true
			}
		}
	}
	return "", false
}

func (c *Config) GetValueNoErr(section, name string) string {
	res, _ := c.GetValue(section, name)
	return res
}

func (c *Config) GetNode(nodeKey string) (map[string]string, bool) {
	conf := c.ReadList()
	for _, v := range conf {
		if mp, ok := v[nodeKey]; ok {
			return mp, true
		}
	}
	return map[string]string{}, false
}

// Set the corresponding value of the key value, if not add, if there is a key change
func (c *Config) SetValue(section, key, value string) bool {
	c.ReadList()
	data := c.conflist
	var ok bool
	var index = make(map[int]bool)
	var conf = make(map[string]map[string]string)
	for i, v := range data {
		_, ok = v[section]
		index[i] = ok
	}

	i, ok := func(m map[int]bool) (i int, v bool) {
		for i, v := range m {
			if v {
				return i, true
			}
		}
		return 0, false
	}(index)

	if ok {
		c.conflist[i][section][key] = value
		return true
	} else {
		conf[section] = make(map[string]string)
		conf[section][key] = value
		c.conflist = append(c.conflist, conf)
		return true
	}
}

// Delete the corresponding key values
func (c *Config) DeleteValue(section, name string) bool {
	c.ReadList()
	data := c.conflist
	for i, v := range data {
		for key, _ := range v {
			if key == section {
				delete(c.conflist[i][key], name)
				return true
			}
		}
	}
	return false
}

// List all the configuration file
func (c *Config) ReadList() []map[string]map[string]string {
	if c.conflist != nil {
		return c.conflist
	}
	file, err := os.Open(c.filepath)
	if err != nil {
		CheckErr(err)
	}
	defer file.Close()
	var data map[string]map[string]string = make(map[string]map[string]string, 10)
	var section string
	buf := bufio.NewReader(file)
	isFirst := true
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				CheckErr(err)
			}
			if len(line) == 0 {
				break
			}
		}
		// ????????????????????????
		if isFirst {
			isFirst = false
			btLine := []byte(line)
			if len(btLine) >= 3 {
				if btLine[0] == 239 && btLine[1] == 187 && btLine[2] == 191 {
					nBtLine := btLine[3:]
					line = string(nBtLine)
				}
			}
		}
		switch {
		case len(line) == 0:
		case string(line[0]) == "#": //????????????????????????
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			if i > 0 {
				value := strings.TrimSpace(line[i+1:])
				valKey := strings.TrimSpace(line[0:i])
				if _, ok := data[section]; !ok {
					data[section] = make(map[string]string, 2)
				}
				data[section][valKey] = value
				if c.uniquappend(section) {
					c.conflist = append(c.conflist, data)
				}
			}
		}

	}

	return c.conflist
}

func CheckErr(err error) string {
	if err != nil {
		return fmt.Sprintf("Error is :'%s'", err.Error())
	}
	return "Notfound this error"
}

// Ban repeated appended to the slice method
func (c *Config) uniquappend(conf string) bool {
	for _, v := range c.conflist {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}
