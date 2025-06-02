package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	cmd    *exec.Cmd
	cancel context.CancelFunc
	mu     sync.Mutex
)

func init() {
	configDir, _ := os.UserConfigDir()
	if configDir == "" {
		configDir = os.TempDir()
	}

	appDir := filepath.Join(configDir, "AMLL_coonector_for_QQMusic")

	os.Mkdir(appDir, 0755)

	viper.SetConfigFile(filepath.Join(appDir, "config.json"))

	if err := viper.ReadInConfig(); err != nil {
		viper.Set("theme", "auto")
		viper.Set("lyrics_path", filepath.Join(appDir, "lyrics"))
		viper.Set("album_art_path", filepath.Join(appDir, "album_art"))
		viper.Set("auto_coonect", false)
		viper.Set("auto_coonect_address", "localhost")
		viper.Set("auto_coonect_port", 11444)

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
