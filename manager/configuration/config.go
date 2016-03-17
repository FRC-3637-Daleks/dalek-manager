/* Team 3637 Dalek Manager - A web base robot configuration manager
   Copyright (C) 2016  Team 3637

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package configuration
import (
	"log"
	"os"
	"github.com/FRC-3637-Daleks/dalek-manager/manager/model"
)

type Config struct {
	Manifest model.Manifest
	WebDir,
	Debug    bool
}

func (config *Config) Log(v ...interface{}) {
	log.Println(v)
}

func (config *Config) DebugLog(v ...interface{}) {
	if (config.Debug) {
		log.Println(v)
	}
}

func (config *Config) ErrorLog(e error) {
	if (e != nil) {
		log.Println(e)
	}
}

func (config *Config) Mkdir(path string, permissions os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.DebugLog("Makeing Directory: ", path)
		err = os.MkdirAll(path, permissions)
		if (err != nil) {
			config.DebugLog(err)
		}
	}
}

func (config *Config) Mkfile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config.DebugLog("Makeing File: ", path)
		_, err = os.Create(path)
		if (err != nil) {
			config.DebugLog(err)
		}
	}
}