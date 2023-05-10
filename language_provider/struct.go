package language_provider

type GetTemplateInput struct {
	MethodName string
	Params     []Param
	Output     Output
}

type GetTemplateOutput struct {
	Content string
}

type Param struct {
	Name string
	Type string
}

type Output struct {
	Name string
	Type string
}

type Testcase struct {
	Input  []string
	Output string
}

type RunTestcaseInput struct {
	MethodName string
	Params     []Param
	Output     Output
	Content    string
	Testcase   []Testcase
}

type RunTestcaseOutput struct {
}
