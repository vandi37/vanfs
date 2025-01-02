package directory_test

import (
	"fmt"
	"testing"

	"github.com/vandi37/vanfs/pkg/directory"
)

func createComplexDirectory() *directory.Directory {
	root := directory.NewRoot("")

	root.AddFile("README.md")
	root.AddFile("config.yaml")
	root.AddFile("main.go")
	root.AddFile(".gitignore")

	root.AddFile("src/utils.go")
	root.AddFile("src/data.json")

	root.AddFile("src/api/api.go")
	root.AddFile("src/api/handlers.go")

	root.AddFile("src/api/auth/auth.go")
	root.AddFile("src/api/auth/token.go")

	root.AddFile("src/api/user/user.go")
	root.AddFile("src/api/user/models.go")

	root.AddFile("tests/main_test.go")
	root.AddFile("tests/integration_test.go")

	root.AddFile("tests/unit/unit_test.go")
	root.AddFile("tests/unit/helpers.go")

	root.AddFile("tests/e2e/e2e_test.go")
	root.AddFile("tests/e2e/setup.go")

	root.AddFile("nested/level1/level2/level3/level4/deep_file.txt")

	root.AddFile("tools/cli.go")

	root.AddFile("docs/api_doc.md")
	root.AddFile("docs/user_guide.md")

	root.AddFile("mixed/config.ini")
	root.AddFile("mixed/README.txt")
	root.AddDir("mixed/sub1")
	root.AddDir("mixed/sub2")

	return root
}

func TestDirectory(t *testing.T) {
	dir := createComplexDirectory()
	fmt.Println(dir)
	res, _ := dir.ToJsonDir().ToJson()
	fmt.Printf("%s\n", res)
	fmt.Println(dir.RemoveDir(""))
	fmt.Println(dir)
}
