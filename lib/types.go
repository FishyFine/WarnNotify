package lib

type Config struct {
	NotifyChannel int  `json:"notify_channel;options=1" yaml:"NotifyChannel;options=1"`
	NeedPing      bool `json:"need_ping" yaml:"NeedPing"`

	TelegramBotToken    string  `json:"telegram_bot_token,optional" yaml:"TelegramBotToken,optional"`
	TelegramNotifyUsers []int64 `json:"telegram_notify_users,optional" yaml:"TelegramNotifyUsers,optional"`
}
