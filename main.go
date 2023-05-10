package main

import (
	"log"
	"test/language_provider"
)

func sum(a int, b int) int {
	return a + b
}

//func main() {
//	// Read the log file
//	data, err := ioutil.ReadFile("output.txt")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// Convert the file data to a string
//	log := string(data)
//
//	// Compile the regex pattern
//	re := regexp.MustCompile(`--- FAIL: (\S+)`)
//
//	matches := re.FindAllStringSubmatch(log, -1)
//	failedTests := []string{}
//
//	for _, match := range matches {
//		failedTests = append(failedTests, match[1])
//	}
//
//	fmt.Printf("Failed tests: %v\n", failedTests)
//
//	re = regexp.MustCompile(`--- PASS: (\S+)`)
//	matches = re.FindAllStringSubmatch(log, -1)
//	passedTests := []string{}
//
//	for _, match := range matches {
//		passedTests = append(passedTests, match[1])
//	}
//
//	fmt.Printf("Passed tests: %v\n", passedTests)
//}

//func main() {
//	// Read the test output from a file
//	content, err := ioutil.ReadFile("output.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Convert the content to a string
//	output := string(content)
//
//	// Find the lines containing test failures
//	pattern := `=== RUN   (.+)
//--- (FAIL|PASS): \1 \((.+)\)
//    (?:.+\n)+?
//    \s+Error Trace:.*\n
//    \s+Error: +Not equal: +\n
//    \s+expected: (\d+)\n
//    \s+actual  : (\d+)\n
//    \s+Test: +\1\n
//    (?:.+\n)*
//    \s+Messages: +(.+)\n`
//
//	re := regexp.MustCompile(pattern)
//	matches := re.FindAllStringSubmatch(output, -1)
//
//	for _, match := range matches {
//		if match[2] == "FAIL" {
//			fmt.Printf("Test case %s failed: expected %s, got %s\n", match[1], match[3], match[4])
//			fmt.Printf("Error message: %s\n", match[5])
//		}
//	}
//}

func main() {
	provider := language_provider.New("Golang")
	output, _ := provider.GetTemplate(language_provider.GetTemplateInput{
		MethodName: "sum",
		//Params:     []language_provider.Param{{Name: "a", Type: "int"}, {Name: "b", Type: "int"}},
		//Output:     language_provider.Output{Name: "result", Type: "int"},
		Params: []language_provider.Param{{Name: "array", Type: "[]int"}},
		Output: language_provider.Output{Name: "result", Type: "int"},
	})
	log.Println(output.Content)

	output2, _ := provider.RunTestcase(language_provider.RunTestcaseInput{
		MethodName: "sum",
		Content:    output.Content,
		Params:     []language_provider.Param{{Name: "array", Type: "[]int"}},
		Output:     language_provider.Output{Name: "result", Type: "int"},
		Testcase: []language_provider.Testcase{
			{Input: []string{"[]int{1,2,3}"}, Output: "6"},
			{Input: []string{"[]int{1,2,3,4}"}, Output: "10"},
			{Input: []string{"[]int{1,2,3,4,5}"}, Output: "15"},
		},
	})

	log.Println(output2)
}
