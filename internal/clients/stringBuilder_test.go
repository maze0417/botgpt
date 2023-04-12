package clients

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringBuilder(t *testing.T) {

	var builder strings.Builder
	builder.WriteString("語音辨識中(Processing your audio) ... \n")

	var content = "可以寫一個C Sharp的Hello World的sample給我看看嗎?感謝"
	txt := fmt.Sprintf("語音轉文字(Transcriptions result): %s\n", content)
	builder.WriteString(txt)
	builder.WriteString("當然可以。以下是一個簡單的C# Hello World 程式碼示例：\n\n```csharp\nusing System;\n\nclass Program {\n    static void Main(string[] args) {\n        Console.WriteLine(\"Hello, World!\");\n    }\n}\n```\n\n這個程式碼的主要功能是使用 `Console.WriteLine` 方法在控制台列印出 \"Hello, World!\"。\n")

	result := builder.String()

	fmt.Println(result)
}
