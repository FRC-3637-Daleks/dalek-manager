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
package model

type Manifest struct {
	Templates struct {
			  Configs struct {
					  Ports     string `json:"ports"`
					  Controls  string `json:"controls"`
					  Settings  string `json:"settings"`
					  LogConfig string `json:"log-config"`
				  } `json:"configs"`
		  } `json:"templates"`
	Runtime   struct {
			  Configs    struct {
					     Ports    string `json:"ports"`
					     Controls string `json:"controls"`
					     Settings string `json:"settings"`
					     Logs     string `json:"logs"`
				     } `json:"configs"`
			  Autonomous string `json:"autonomous"`
			  Binary     string `json:"binary"`
		  } `json:"runtime"`
	Server    struct {
				  WebRoot  string `json:"web-root"`
				  Port     int     `json:"port"`
				  MqttHost string `json:"mqttHost"`
				  MqttPort int `json:"mqttPort"`
			  } `json:"server"`
}