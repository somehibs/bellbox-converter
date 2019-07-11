package converter

// Translation rules with key checks win when picking one
type TranslationRule struct {
	Name string
	KeyExists []string
	Encoding string
	Map map[string]string
	Format map[string]string
}

type TranslationFile []TranslationRule

type ConvertRule struct {
	Target string
	Path string
	Ruleset string
	SenderName string
}

type TranslationConfig struct {
	Bellbox string
	Default ConvertRule
	Convert []ConvertRule
}
