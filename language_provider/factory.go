package language_provider

type Provider interface {
	GetTemplate(input GetTemplateInput) (*GetTemplateOutput, error)
	RunTestcase(input RunTestcaseInput) (*RunTestcaseOutput, error)
}

func New(provider string) Provider {
	if provider == "Golang" {
		return NewGoProvider()
	}
	return nil
}
