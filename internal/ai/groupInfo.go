package ai

type GroupSetting struct {
	SystemMessage string
	CommandMode   string
	Enable        bool
	LLMModel      string
}

const (
	DefaultMode            = "chat"
	LineGroupGpt3Turbo     = "Cede0b311a552b14f418db2e1acfd88de"
	TelegramGroupGpt3Turbo = "-910957863"
	PrivateTestBot         = "1066396636"
	PhTgProject            = "-728297760"
	LineCard               = "Cd893d1f8eb0e28de1a4b06a0237bdbb2"
	CsUpdated              = "-1001829478512"
	Wife                   = "Uceccedb9ead38da1c98e35a24325a372"

	GPTLatestModel  = "gpt-4-vision-preview"
	GPTCheaperModel = "gpt-3.5-turbo-1106"
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
		SystemMessage: "你是一個專業的聊天好朋友，會用可愛女生的口氣回答各種你知道的問題，當不知道問題的時候，會說: '不好意思 你的問題太難了 可以教教我嗎?鳩咪。'",
		CommandMode:   Chat,
		Enable:        true,
		LLMModel:      GPTLatestModel,
	},
	PhTgProject: {
		CommandMode: EnToTw,
		Enable:      false,
	},
	LineCard: {
		SystemMessage: "作為一個專業塔羅牌大師，具有廣泛的知識，包括所有78張塔羅牌的含義以及各種不同的塔羅牌展開方式。你可以幫助用戶進行塔羅牌抽牌以及讀牌，並解釋每張牌的象徵意義和可能的解讀。你的目標是提供專業，深入和有洞察力的塔羅牌解讀，並且溫馨提醒該注意的事情，以幫助用戶理解他們的問題和情況",
		CommandMode:   Chat,
		Enable:        true,
		LLMModel:      GPTLatestModel,
	},
	Wife: {
		SystemMessage: "你是一個專業的室內設計顧問，可以回答跟室內設計有關的專業問題，包含報價，合約，等製作問題",
		CommandMode:   Chat,
		Enable:        true,
		LLMModel:      GPTLatestModel,
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
