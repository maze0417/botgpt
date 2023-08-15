package ai

type GroupSetting struct {
	SystemMessage string
	CommandMode   string
	Enable        bool
}

const (
	DefaultMode            = "chat"
	LineGroupGpt3Turbo     = "Cede0b311a552b14f418db2e1acfd88de"
	TelegramGroupGpt3Turbo = "-910957863"
	PrivateTestBot         = "1066396636"
	PhTgProject            = "-728297760"
	LineCard               = "Cd893d1f8eb0e28de1a4b06a0237bdbb2"
	CsUpdated              = "-859160990"
)

var GroupMapping = map[string]GroupSetting{
	PrivateTestBot: {
		CommandMode: ChatWithoutTag,
		Enable:      false,
	},
	LineGroupGpt3Turbo: {
		CommandMode: Chat,
		Enable:      false,
	},
	TelegramGroupGpt3Turbo: {
		CommandMode: Chat,
		Enable:      false,
	},
	PhTgProject: {
		CommandMode: EnToTw,
		Enable:      false,
	},
	LineCard: {
		SystemMessage: "作為一個專業塔羅牌大師，具有廣泛的知識，包括所有78張塔羅牌的含義以及各種不同的塔羅牌展開方式。你可以幫助用戶進行塔羅牌抽牌以及讀牌，並解釋每張牌的象徵意義和可能的解讀。你的目標是提供專業，深入和有洞察力的塔羅牌解讀，並且溫馨提醒該注意的事情，以幫助用戶理解他們的問題和情況",
		CommandMode:   Chat,
		Enable:        true,
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
		v = GroupSetting{
			CommandMode: cmd,
			Enable:      true,
		}
		GroupMapping[groupID] = v
		return &v
	}
	cmdInfo := GetCommandBy(cmd)
	if cmdInfo == nil {
		return nil
	}

	v.CommandMode = cmdInfo.Cmd
	if v.Enable {
		v.Enable = false
	} else {
		v.Enable = true
	}

	GroupMapping[groupID] = v

	return &v
}
