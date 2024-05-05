package ybtools

//	   2024-05-05	DAT   Add debug statement to see name of file it is trying to open

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/metal3d/go-slugify"
	"gopkg.in/yaml.v2"
)

//
// Yapperbot Tools, the internal system bits for Yapperbot and co.
// Copyright (C) 2020 Naypta

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//

type configObject struct {
	APIEndpoint string
	BotUsername string
}

// acts like an interface for config files
// edit limits from tool configs are unloaded into here
type toolConfigWithEditLimit struct {
	EditLimit int64
}

const localConfigFilename string = "config.yml"
const globalConfigFilename string = "config-global.yml"
const botPasswordFilename string = "botpassword"

var botPassword string
var config configObject
var taskConfigFile []byte

// ParseTaskConfig takes in a pointer to a config object
// and populates the config object with the task configuration
func ParseTaskConfig(cobj interface{}) {
	yaml.Unmarshal(taskConfigFile, cobj)
}

func init() {
	botConfigFile, err := ioutil.ReadFile(findConfigFile(localConfigFilename, globalConfigFilename))
	if err != nil {
		PanicErr("Bot config file could not be read at detected path!")
	}
	err = yaml.UnmarshalStrict(botConfigFile, &config)
	if err != nil {
		PanicErr("Bot config file was invalid!")
	}

	botPasswordFile, err := ioutil.ReadFile(findConfigFile(botPasswordFilename, botPasswordFilename))
	if err != nil {
		PanicErr("Bot password file could not be read at detected path!")
	}
	botPassword = string(botPasswordFile)
}

func setupTaskConfigFile() {
	var err error
	var taskConfigForEditLimit toolConfigWithEditLimit

	taskConfigFile, err = ioutil.ReadFile("config-" + strings.ToLower(slugify.Marshal(settings.TaskName)) + ".yml")
	if err != nil {
// add for debug
		log.Println("config-" + strings.ToLower(slugify.Marshal(settings.TaskName)) + ".yml")
//
		log.Println("No task-specific config file found, ignoring")
	}

	// Immediately parse the file for an edit limit and only an edit limit
	yaml.Unmarshal(taskConfigFile, &taskConfigForEditLimit)
	if taskConfigForEditLimit.EditLimit > 0 {
		setupEditLimit(taskConfigForEditLimit.EditLimit)
	}
}

// findConfigFile takes a local filename as a string, and a global filename as a string
// It then finds the local filename in the current directory, or if there isn't one,
// the global filename in the parent directory, and returns the name of the file it found.
// If it doesn't find either, it fatally errors.
func findConfigFile(filename string, globalfilename string) string {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if _, err := os.Stat("../" + globalfilename); os.IsNotExist(err) {
			PanicErr("Couldn't find a config file for ", filename, " either in this directory or in the one above it!")
		}
		return "../" + globalfilename
	}
	return filename
}
