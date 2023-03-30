package domain

type CommandHandler interface {
	Reply(id int64) string
}
