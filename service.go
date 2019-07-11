package converter

import (
	"os"
	"fmt"
	"io"
	"strings"
	"io/ioutil"
	"path/filepath"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"git.circuitco.de/self/bellbox"
)

func New() error {
	// Configure the converter and map all the endpoints.
	config, err := LoadConfig()
	if err != nil {
		return err
	}
	translations, err := LoadTranslations()
	if err != nil {
		return err
	}
	router := gin.Default()
	return Route(router, config, translations)
}

type ConvertRoute struct {
	rule ConvertRule
	ruleset TranslationFile
	sender bellbox.SenderAuth
}

func (cr *ConvertRoute) Handle(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func Route(router *gin.Engine, config TranslationConfig, tlrules map[string]TranslationFile) error {
	//d := config.Default
	for _, conversion := range config.Convert {
		if tlrules[conversion.Ruleset] == nil {
			fmt.Printf("Refusing to add conversion rule: %s path: %s (missing ruleset)", conversion.Ruleset, conversion.Path)
			panic("")
			continue
		}
		sender := bellbox.StartSender(conversion.SenderName, config.Bellbox)
		convertRoute := ConvertRoute{conversion, tlrules[conversion.Ruleset], sender}
		if conversion.SenderName == "" {
			conversion.SenderName = conversion.Ruleset
		}
		router.POST(conversion.Path, convertRoute.Handle)
	}
	return router.Run("localhost:9009")
}

func LoadConfig() (TranslationConfig, error) {
	config := TranslationConfig{}
	return config, LoadJson("config.json", &config)
}

func LoadTranslations() (map[string]TranslationFile, error) {
	ret := map[string]TranslationFile{}
	err := filepath.Walk("translations/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != "translations/" {
			// don't descend into directories
			fmt.Println("Skipping unknown dir: " + path)
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.Contains(path, ".json") {
			tf := TranslationFile{}
			err = LoadJson(path, &tf)
			if err != nil {
				return err
			}
			simpleName := strings.Replace(info.Name(), ".json", "", 1)
			//fmt.Printf("File: %s dat: %+v\n", simpleName, tf)
			ret[simpleName] = tf
		}
		return nil
	})
	return ret, err
}

func LoadJson(file string, item interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return LoadJsonImpl(f, item)
}

func LoadJsonImpl(r io.Reader, item interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, item)
}
