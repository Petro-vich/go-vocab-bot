package domain

type Translator interface {
	Translate(text string, targetLang string) (string, error)
}
