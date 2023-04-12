package ai

type GroupSetting struct {
	DefaultCommandMode string
	Enable             bool
}

const (
	Default                = "你是一個有用AI助手，會嘗試回覆各種問題"
	LineGroupGpt3Turbo     = "Cede0b311a552b14f418db2e1acfd88de"
	TelegramGroupGpt3Turbo = "-910957863"
	PrivateTestBot         = "1066396636"
	PhTgProject            = "-728297760"
)

var GroupMapping = map[string]GroupSetting{
	PrivateTestBot: {
		DefaultCommandMode: Private,
		Enable:             false,
	},
	LineGroupGpt3Turbo: {
		DefaultCommandMode: LineBot,
		Enable:             false,
	},
	TelegramGroupGpt3Turbo: {
		DefaultCommandMode: TgBot,
		Enable:             false,
	},
	PhTgProject: {
		DefaultCommandMode: EnToTw,
		Enable:             true,
	},
}

type GptMsgModel struct {
	IsReply     bool
	ReplyString string
	Message     string
}

func GetGroupMode(groupID string) *GroupSetting {

	v, ok := GroupMapping[groupID]
	if !ok {
		return nil
	}
	if !v.Enable {
		return nil
	}
	return &v
}

func SetGroupMode(groupID string, cmd string) *GroupSetting {

	v, ok := GroupMapping[groupID]
	if !ok {
		return nil
	}
	cmdInfo := GetCommandBy(cmd)
	if cmdInfo == nil {
		return nil
	}

	v.DefaultCommandMode = cmdInfo.Cmd
	if v.Enable {
		v.Enable = false
	}
	if !v.Enable {
		v.Enable = true
	}

	GroupMapping[groupID] = v

	return &v
}
