package internal

import (
	"fmt"
	"github.com/spf13/cobra"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/exp/slices"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var WireCmd = &cobra.Command{
	Use:   "wire",
	Short: "自识别 internal 目录下需要注入的方法，写入到 provider_set.go 文件中",
	Long:  `自识别 internal 目录下需要注入的方法，写入到 provider_set.go 文件中`,
	Run: func(cmd *cobra.Command, args []string) {
		wire := Wire{}
		wire.Run()
	},
}

var wireExcludePath = []string{"internal/data/ent"}
var wireExclude = ""
var wireSuccessCnt = 0

func init() {
	WireCmd.Flags().StringVarP(&wireExclude, "wireExclude", "e", "", "手动指定要排除的目录，多个用,号隔开，也可以将此项配置到项目根目录的.tools.yaml中")
}

type Wire struct {
}

func (w *Wire) Run() {
	//解析配置
	if wireExclude != "" {
		wireExcludePath = append(wireExcludePath, strings.Split(wireExclude, ",")...)
	}
	if len(Config.Wire.Exclude) > 0 {
		wireExcludePath = append(wireExcludePath, Config.Wire.Exclude...)
	}
	// 递归遍历目录读取文件
	w.walkDir("./internal")
	fmt.Printf("已注入 %d 个方法\n", wireSuccessCnt)
}

// 递归遍历目录
func (w *Wire) walkDir(dir string) {
	if slices.Contains(wireExcludePath, dir) {
		return
	}

	// 递归遍历目录读取文件
	haveProviderSet := false
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, info := range files {
		if info.IsDir() {
			subPath := dir + string(os.PathSeparator) + info.Name()
			w.walkDir(subPath)
		}
		if info.Name() == "provider_set.go" {
			haveProviderSet = true
		}
	}
	if haveProviderSet {
		w.writeProvider(dir)
	}
}

func (w *Wire) writeProvider(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() || path.Ext(file.Name()) != ".go" {
			continue
		}
		filePath := dir + string(os.PathSeparator) + file.Name()
		if dir == "./internal/data" && file.Name() == "transaction_repo_impl.go" {
			//这个文件在其它地方已经注入了
			continue
		}

		fs := token.NewFileSet()
		//获取此路径文件中的所有方法名
		f, pErr := parser.ParseFile(fs, filePath, nil, parser.ParseComments)
		if pErr != nil {
			panic(pErr)
		}
		for _, decl := range f.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				// decl.Name.Name 以New 开头的函数
				if funcDecl.Recv != nil || len(funcDecl.Name.Name) <= 3 || funcDecl.Name.Name[:3] != "New" {
					continue
				}
				w.writeCodeToProvider(dir, funcDecl.Name.Name)
			}
		}
	}
}

func (w *Wire) writeCodeToProvider(dir string, funcName string) {
	providerPath := dir + string(os.PathSeparator) + "provider_set.go"
	content, err := os.ReadFile(providerPath)
	if err != nil {
		panic(err)
	}
	code := string(content)
	if strings.Contains(code, funcName) {
		// 已经存在了
		return
	}
	writeCode := `wire.NewSet(
	` + funcName + `,
`

	newCode := strings.Replace(code, "wire.NewSet(", writeCode, 1)
	reg := regexp.MustCompile("\n{2,}")
	newCode = reg.ReplaceAllString(newCode, "\n")
	err = os.WriteFile(providerPath, []byte(newCode), 0600)
	fullPath, _ := filepath.Abs(providerPath)
	if err != nil {
		fmt.Printf("写入注入代码文件失败 %s %s", fullPath, err)
		panic(err)
	}
	wireSuccessCnt++
	fmt.Printf("已注入 funcName, file://%s\n", fullPath)
}
