package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	cmd *exec.Cmd
)

func init() {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		configDir = os.TempDir()
	}

	appDir := filepath.Join(configDir, "AMLL_coonector_for_QQMusic")
	var defaultConfig = map[string]any{
		"theme":                "auto",
		"lyrics_path":          filepath.Join(appDir, "lyrics"),
		"album_art_path":       filepath.Join(appDir, "album_art"),
		"auto_connect":         false,
		"auto_connect_address": "localhost",
		"auto_connect_port":    11444,
		"app_version":          APPVersion,
		"auto_change_volume":   true,
	}
	os.Mkdir(appDir, 0755)

	viper.SetConfigFile(filepath.Join(appDir, "config.json"))

	if err := viper.ReadInConfig(); err != nil {
		viper.WriteConfig()
	}

	// 创建文件夹
	if _, err := os.Stat(viper.GetString("lyrics_path")); os.IsNotExist(err) {
		os.Mkdir(viper.GetString("lyrics_path"), 0755)
	}
	if _, err := os.Stat(viper.GetString("album_art_path")); os.IsNotExist(err) {
		os.Mkdir(viper.GetString("album_art_path"), 0755)
	}

	//startOrRestartSubprocess()
	// 自动配置文件
	for key, value := range defaultConfig {
		// 检查字段是否存在
		if !viper.IsSet(key) {
			viper.Set(key, value)
		}
		if key == "app_version" {
			if viper.GetString(key) != APPVersion {
				viper.Set(key, APPVersion)
			}
		}
		viper.WriteConfig()
	}

}

func (g *GreetService) GetConfig() string {
	return getConfigString()
}

func (g *GreetService) SetConfig(key string, data any) string {
	viper.Set(key, data)
	viper.WriteConfig()
	viper.ReadInConfig()
	return getConfigString()
}

func (g *GreetService) GetConfigByKey(key string) any {
	return viper.Get(key)
}
func (g *GreetService) SetAllConfig(data string) string {
	us, _ := os.UserConfigDir()
	if err := os.WriteFile(filepath.Join(us, "AMLL_coonector_for_QQMusic", "config.json"), []byte(data), 0644); err != nil {
		return ""
	}
	viper.ReadInConfig()
	return getConfigString()
}
func getConfigString() string {
	us, _ := os.UserConfigDir()
	content, err := os.ReadFile(filepath.Join(us, "AMLL_coonector_for_QQMusic", "config.json"))
	if err != nil {
		return ""
	}
	return string(content)
}
