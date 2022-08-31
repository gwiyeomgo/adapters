package config

import (
	"github.com/jinzhu/configor"
)

var Config = struct {
	Dooray struct {
		ApiKey  string
		Project struct {
			Url  string
			List struct {
				ErrorEvent struct {
					ProjectNo            string
					ProjectMemberGroupId string
				}
			}
		}
	}
}{}

func ConfigureEnvironment(path string) {
	configor.Load(&Config, path+"config/config.json")
}
