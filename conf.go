/**
 * Read the configuration file
 *
 * @copyright           (C) 2014  widuu
 * @lastmodify          2014-2-22
 * @website		http://www.widuu.com
 *
 */

package goini

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Config struct {
	r        io.Reader                      //your ini file path directory+file
	conflist []map[string]map[string]string //configuration information slice
}

//Create an empty configuration file
func SetConfig(r io.Reader) *Config {
	c := new(Config)
	c.Reload(r)
	return c
}
func (c *Config) Reload(r io.Reader) {
	c.r = r
	c.ReadList()
}

//To obtain corresponding value of the key values
func (c *Config) GetValue(section, name string) string {
	//c.ReadList()
	//conf := c.ReadList()
	for _, v := range c.conflist {
		for key, value := range v {
			if key == section {
				return value[name]
			}
		}
	}
	return "no value"
}

//Set the corresponding value of the key value, if not add, if there is a key change
func (c *Config) SetValue(section, key, value string) bool {
	//c.ReadList()
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
			if v == true {
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

	return false
}

//Delete the corresponding key values
func (c *Config) DeleteValue(section, name string) bool {
	//c.ReadList()
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

//List all the configuration file
func (c *Config) ReadList() []map[string]map[string]string {
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(c.r)
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
		switch {
		case len(line) == 0:
		case string(line[0]) == "#": //增加配置文件备注
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			if i == -1 {
				continue
			}
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if c.uniquappend(section) == true {
				c.conflist = append(c.conflist, data)
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

//Ban repeated appended to the slice method
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
