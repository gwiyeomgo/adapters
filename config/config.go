package config

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
				ErrorReport struct {
					ProjectNo            string
					ProjectMemberGroupId string
				}
				ErrorNotFoundEvent struct {
					ProjectNo            string
					ProjectMemberGroupId string
				}
				EtcErrorNotFoundEvent struct {
					ProjectNo            string
					ProjectMemberGroupId string
				}
			}
		}
	}
}{}

