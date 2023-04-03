package telegram_bot

import "fmt"

type TextTemplate string

var (
	tmplArbitrageNotify TextTemplate = `<strong>&#128060;&#128060;&#128060;  Notify &#128060;&#128060;&#128060;</strong>
<strong>=======================</strong>
<strong>Spread: </strong><u>%s</u>
<strong>Invested Amount: </strong><u>%s</u>
<strong>%s Buy: </strong><u>%s</u>
<strong>%s Sell: </strong><u>%s</u>
<strong>Arbitrage: </strong><u>%s</u>
<strong>Estimated Profit: </strong><u>%s</u>
<strong>Author: </strong><a href="tg://user?id=%s">%s</a>
`

	tmplErrorNotify TextTemplate = `<strong> Error Notification </strong>
<strong>=======================</strong>
<strong>Title: </strong><u>%s</u>
<strong>Error Message: </strong><u>%s</u>
<strong>Time: </strong><u>%s</u>
`
)

func (t TextTemplate) String() string {
	return string(t)
}

func (t TextTemplate) Format(args ...interface{}) string {
	return fmt.Sprintf(t.String(), args...)
}

func (t TextTemplate) Type() TemplateType {
	switch t {
	case tmplArbitrageNotify:
		return HTML
	case tmplErrorNotify:
		return HTML
	default:
		return PlainText
	}
}

type TemplateType string

var (
	PlainText  TemplateType = "PlainText"
	HTML       TemplateType = "HTML"
	MarkdownV2 TemplateType = "MarkdownV2"
)

func (t TemplateType) String() string {
	return string(t)
}
