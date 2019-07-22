package config

import (
	"errors"
	"github.com/larspensjo/config"
	. "github.com/redochen/tools/string"
)

type Config struct {
	FilePath string
	Config   *config.Config
}

//是否有效
func (c *Config) IsValid() bool {
	if c.Config != nil {
		return true
	} else {
		return false
	}
}

//是否为默认节点
func (c *Config) IsDefaultSection(section string) bool {
	if "" == section || len(section) <= 0 {
		return true
	}

	if defaultSection == section {
		return true
	}

	return false
}

//检查是否存在特定的配置节（配置节总是存在）
func (c *Config) HasSection(section string) bool {
	if nil == c.Config {
		return false
	}

	return c.Config.HasSection(section)
}

//获取所有节点列表
func (c *Config) GetSections() []string {
	if nil == c.Config {
		return nil
	}

	return c.Config.Sections()
}

//检查是否存在特定的配置项
func (c *Config) HasOption(section, option string) bool {
	if nil == c.Config {
		return false
	}

	return c.Config.HasOption(section, option)
}

//获取某个节点下的所有配置项
func (c *Config) GetOptions(section string) ([]string, error) {
	if nil == c.Config {
		return nil, errors.New("section not exist")
	}

	return c.Config.Options(section)
}

//读取默认配置节下的特定配置项
func (c *Config) DefaultString(option string) (string, error) {
	if nil == c.Config {
		return "", errors.New("config not initialized")
	}

	return c.Config.RawStringDefault(option)
}

//读取默认配置节下的特定配置项（不存或失败时返回默认值）
func (c *Config) DefaultStringEx(option, defaultValue string) string {
	if nil == c.Config {
		return defaultValue
	}

	value, _ := c.Config.RawStringDefault(option)

	return CcStr.FirstValid(value, defaultValue)
}

//读取特定配置节下的特定配置项
func (c *Config) String(section, option string) (string, error) {
	if nil == c.Config {
		return "", errors.New("config not initialized")
	}

	return c.Config.String(section, option)
}

//读取特定配置节下的特定配置项（不存或失败时返回默认值）
func (c *Config) StringEx(section, option, defaultValue string) string {
	if nil == c.Config {
		return defaultValue
	}

	if !c.HasOption(section, option) {
		return defaultValue
	}

	value, _ := c.Config.String(section, option)

	return CcStr.FirstValid(value, defaultValue)
}

//读取默认配置节下的特定配置项
func (c *Config) DefaultInt(option string) (int, error) {
	str, err := c.DefaultString(option)
	if err != nil {
		return 0, err
	}

	return CcStr.ParseInt(str), nil
}

//读取默认配置节下的特定配置项（不存或失败时返回默认值）
func (c *Config) DefaultIntEx(option string, defaultValue int) int {
	str, _ := c.DefaultString(option)
	if str == "" {
		return defaultValue
	}

	return CcStr.ParseInt(str)
}

//读取特定配置节下的特定配置项
func (c *Config) Int(section, option string) (int, error) {
	str, err := c.String(section, option)
	if err != nil {
		return 0, err
	}

	return CcStr.ParseInt(str), nil
}

//读取特定配置节下的特定配置项（不存或失败时返回默认值）
func (c *Config) IntEx(section, option string, defaultValue int) int {
	str, _ := c.String(section, option)
	if str == "" {
		return defaultValue
	}

	return CcStr.ParseInt(str)
}

//读取默认配置节下的特定配置项
func (c *Config) DefaultBool(option string) (bool, error) {
	str, err := c.DefaultString(option)
	if err != nil {
		return false, err
	}

	return CcStr.ParseBool(str), nil
}

//读取默认配置节下的特定配置项（不存或失败时返回默认值）
func (c *Config) DefaultBoolEx(option string, defaultValue bool) bool {
	str, _ := c.DefaultString(option)
	if str == "" {
		return defaultValue
	}

	return CcStr.ParseBool(str)
}

//读取特定配置节下的特定配置项
func (c *Config) Bool(section, option string) (bool, error) {
	str, err := c.String(section, option)
	if err != nil {
		return false, err
	}

	return CcStr.ParseBool(str), nil
}

//读取特定配置节下的特定配置项（不存或失败时返回默认值）
func (c *Config) BoolEx(section, option string, defaultValue bool) bool {
	str, _ := c.String(section, option)
	if str == "" {
		return defaultValue
	}

	return CcStr.ParseBool(str)
}

//读取默认配置节下的特定配置项
func (c *Config) DefaultFloat(option string) (float32, error) {
	str, err := c.DefaultString(option)
	if err != nil {
		return 0, err
	}

	return CcStr.ParseFloat(str), nil
}

//读取默认配置节下的特定配置项（不存或失败时返回默认值）
func (c *Config) DefaultFloatEx(option string, defaultValue float32) float32 {
	str, _ := c.DefaultString(option)
	if str == "" {
		return defaultValue
	}

	return CcStr.ParseFloat(str)
}

//读取特定配置节下的特定配置项
func (c *Config) Float(section, option string) (float32, error) {
	str, err := c.String(section, option)
	if err != nil {
		return 0, err
	}

	return CcStr.ParseFloat(str), nil
}

//读取特定配置节下的特定配置项（不存或失败时返回默认值）
func (c *Config) FloatEx(section, option string, defaultValue float32) float32 {
	str, _ := c.String(section, option)
	if str == "" {
		return defaultValue
	}

	return CcStr.ParseFloat(str)
}
