package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"git.circuitco.de/self/bellbox"
	"github.com/gin-gonic/gin"
)

func New() error {
	// Configure the converter and map all the endpoints.
	config, err := LoadConfig()
	if err != nil {
		return err
	}
	translations := LoadTranslations()
	router := gin.Default()
	return Route(router, config, translations)
}

type ConvertRoute struct {
	rule    ConvertRule
	ruleset Translation
	sender  bellbox.SenderAuth
}

func (cr *ConvertRoute) Handle(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{})
	}
	msg := cr.ruleset.Handle(b)
	if msg.Title != "" {
		cr.sender.Send(cr.rule.Target, msg.Title, msg.Message)
	}
	fmt.Printf("Integration returned: %+v\n", msg)
	c.JSON(200, gin.H{})
}

func Route(router *gin.Engine, config TranslationConfig, tlrules map[string]Translation) error {
	//d := config.Default
	for _, conversion := range config.Convert {
		if tlrules[conversion.Ruleset] == nil {
			fmt.Printf("Refusing to add conversion rule: %s path: %s (missing ruleset)", conversion.Ruleset, conversion.Path)
			panic("")
			continue
		}
		if conversion.SenderName == "" {
			conversion.SenderName = conversion.Ruleset
		}
		sender := bellbox.StartSender(conversion.SenderName, config.Bellbox)
		if conversion.Target == "" {
			conversion.Target = config.Default.Target
		}
		if conversion.Target == "" {
			panic("")
		}
		_, err := sender.SingleTarget(conversion.Target)
		if err != nil {
			panic(err.Error())
		}
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

func LoadTranslations() map[string]Translation {
	ret := map[string]Translation{}
	ret["prometheus"] = Prometheus{}
	return ret
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
