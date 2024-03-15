package config

import (
	"github.com/jinzhu/configor"
	"log"
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
	Mail struct {
		Password   string
		ServerHost string
		ServerPort string
		Username   string
		SenderAddr string
	}
	AwsS3 struct {
		SecretAccessKey string
		Bucket          string
		Region          string
		AccessKeyId     string
		HttpEndPoint    string
	}
}{}

func ConfigureEnvironment(path string) {
	err := configor.Load(&Config, path+"config/config.json")
	if err != nil {
		log.Fatalln(err)
	}
}
