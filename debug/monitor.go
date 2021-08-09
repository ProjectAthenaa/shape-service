package debug

import (
	"github.com/ProjectAthenaa/shape"
	"github.com/ProjectAthenaa/shape/deobfuscation"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type site struct {
	Mod          string
	Source       string
	GlobalHolder *deobfuscation.GlobalHolder
}

var (
	NewBalance = site{}
	Target     = site{}
)

func GetShapeVersions() {
	for {
		target()
		newBalance()
	}
}

func target() {
	urlRe := regexp.MustCompile(`"http([^"]+)"`)
	c := http.Client{}

	resp, _ := c.Get("http://assets.targetimg1.com/ssx/ssx.mod.js?async")
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	Target.Mod = string(body)


	shapeURL := strings.ReplaceAll(urlRe.FindString(string(body)), "\"", "")

	resp, _ = c.Get(shapeURL)
	body, _ = ioutil.ReadAll(resp.Body)
	Target.Source = string(body)

	Target.GlobalHolder = shape.CreateDeobfuscator(Target.Source, Target.Mod)
}

func newBalance() {

}
