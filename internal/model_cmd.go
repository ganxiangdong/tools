package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "生成DO、Repo、RepoImpl文件，并自动写入 Impl 构造方法到 provider_set.go 文件",
	Long:  `生成DO、Repo、RepoImpl文件，并自动写入 Impl 构造方法到 provider_set.go 文件，根据已经生成好的 ent 模型文件来生成`,
	Run: func(cmd *cobra.Command, args []string) {
		Generate()
	},
}

var tableName string
var poFileName string
var poFilePath string
var modName string
var isOverwrite bool = false

func init() {
	ModelCmd.Flags().StringVarP(&tableName, "table", "t", "", "指定模型对应的表名，生成前请确定已经生成了 ent 模型，service 项目一般用 make ent t=表名来生成模型")
}

var doTemplate = `package do

type {{.Name}}Do struct {
{{ range $index, $value := .Fields }}	{{$value.name}} {{$value.type}} 
{{ end }}}`

var repoTemplate = `package repository

import (
	"context"
	"{{.ModName}}/internal/biz/bo"
	"{{.ModName}}/internal/biz/do"
)

type {{.Name}}Repo interface {
	// Get 通过 id 获取一条数据，出错则 panic
	Get(ctx context.Context, id int) *do.{{.Name}}Do

	// Find 通过多个 id 获取多条数据，出错则 panic
	Find(ctx context.Context, ids... int) []*do.{{.Name}}Do

	// Create 创建数据，出错则 panic
	Create(ctx context.Context, createData *do.{{.Name}}Do) *do.{{.Name}}Do

	// CreateBulk 批量创建数据，出错则 panic
	CreateBulk(ctx context.Context, dos []*do.{{.Name}}Do) []*do.{{.Name}}Do

	// Update 更新数据，出错则 panic
	Update(ctx context.Context, updateData *do.{{.Name}}Do) int

	// Delete 删除数据，出错则 panic
	Delete(ctx context.Context, ids... int) int

	// SearchList 搜索列表，出错则 panic
	SearchList(ctx context.Context, reqBo *bo.ReqPageBo) (dos []*do.{{.Name}}Do, respPage *bo.RespPageBo)

	// GetE 通过 id 获取一条数据
	GetE(ctx context.Context, id int) (*do.{{.Name}}Do, error)

	// FindE 通过多个 id 获取多条数据
	FindE(ctx context.Context, ids... int) ([]*do.{{.Name}}Do, error)

	// CreateE 创建数据
	CreateE(ctx context.Context, creatData *do.{{.Name}}Do) (*do.{{.Name}}Do, error)

	// CreateBulkE 批量创建数据
	CreateBulkE(ctx context.Context, dos []*do.{{.Name}}Do) ([]*do.{{.Name}}Do, error)

	// UpdateE 更新数据
	UpdateE(ctx context.Context, updateData *do.{{.Name}}Do) (int, error)

	// DeleteE 软删除数据
	DeleteE(ctx context.Context, ids... int) (int, error)

	// DeleteForceE 删除数据(硬删除)
	DeleteForceE(ctx context.Context, ids... int) (int, error)

	// SearchListE 搜索列表
	SearchListE(ctx context.Context, reqBo *bo.ReqPageBo) (dos []*do.{{.Name}}Do, respPage *bo.RespPageBo, err error)
}`

var repoImplTemplate = `package repositoryimpl

import (
	"context"
	"time"

	"{{.ModName}}/internal/biz/bo"
	"{{.ModName}}/internal/biz/do"
	"{{.ModName}}/internal/biz/repository"
	"{{.ModName}}/internal/data"
	"{{.ModName}}/internal/data/ent"
	"{{.ModName}}/internal/data/ent/{{.EntPackage}}"
)

type {{.Name}}RepoImpl struct {
	Base[ent.{{.Name}}, do.{{.Name}}Do, ent.{{.Name}}Query]
	data *data.Data
}

// New{{.Name}}RepoImpl 创建 {{.Name}}Repo的实现者
func New{{.Name}}RepoImpl(data *data.Data) repository.{{.Name}}Repo {
	return &{{.Name}}RepoImpl{data: data}
}

// ToEntity 转换成实体
func ({{.RecName}} *{{.Name}}RepoImpl) ToEntity(po *ent.{{.Name}}) *do.{{.Name}}Do {
	if po == nil {
		return nil
	}
	return {{.RecName}}.Base.ToEntity(po)
}

// ToEntities 转换成实体
// 支持基本类型的值对象
func ({{.RecName}} *{{.Name}}RepoImpl) ToEntities(pos []*ent.{{.Name}}) []*do.{{.Name}}Do {
	if pos == nil {
		return nil
	}
	// 使用循环，以免单个转换有特殊处理，要修改两个地方
	entities := make([]*do.{{.Name}}Do, len(pos))
	for k, p := range pos {
		entities[k] = {{.RecName}}.ToEntity(p)
	}
	return entities
}

// Get 通过 id 获取一条数据，出错则 panic
func ({{.RecName}} *{{.Name}}RepoImpl) Get(ctx context.Context, id int) *do.{{.Name}}Do {
	row := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Query().Where({{.EntPackage}}.ID(id)).FirstX(ctx)
	return {{.RecName}}.ToEntity(row)
}

// Find 通过多个 id 获取多条数据，出错则 panic
func ({{.RecName}} *{{.Name}}RepoImpl) Find(ctx context.Context, ids... int) []*do.{{.Name}}Do {
	if len(ids) == 0 {
		return nil
	}
	rows := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Query().Where({{.EntPackage}}.IDIn(ids...)).AllX(ctx)
	return {{.RecName}}.ToEntities(rows)
}

// Create 创建数据，出错则 panic
func ({{.RecName}} *{{.Name}}RepoImpl) Create(ctx context.Context, createData *do.{{.Name}}Do) *do.{{.Name}}Do {
	row := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Create().
		{{ range $index, $value := .Fields }}{{ if and (ne $value.name "ID") (ne $value.name "Edges")}}Set{{$value.name}}(createData.{{$value.name}}).
		{{ end }}{{ end }}SaveX(ctx)
	return {{.RecName}}.ToEntity(row)
}

// CreateBulk 批量创建数据，出错则 panic
func ({{.RecName}} *{{.Name}}RepoImpl) CreateBulk(ctx context.Context, dos []*do.{{.Name}}Do) []*do.{{.Name}}Do {
	if len(dos) == 0 {
		return nil
	}
	values := make([]*ent.{{.Name}}Create, len(dos))
	for i, item := range dos {
		values[i] = {{.RecName}}.data.GetDb(ctx).{{.Name}}.Create(){{ range $index, $value := .Fields }}{{ if and (ne $value.name "ID") (ne $value.name "Edges")}}.
			Set{{$value.name}}(item.{{$value.name}}){{ end }}{{ end }}
	}
	rows := {{.RecName}}.data.GetDb(ctx).{{.Name}}.CreateBulk(values...).SaveX(ctx)
	return {{.RecName}}.ToEntities(rows)
}

// Update 更新数据，出错则 panic
func ({{.RecName}} *{{.Name}}RepoImpl) Update(ctx context.Context, updateData *do.{{.Name}}Do) int {
	cnt := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Update().Where({{.EntPackage}}.ID(updateData.ID)).
		{{ range $index, $value := .Fields }}{{ if and (ne $value.name "ID") (ne $value.name "Edges")}}Set{{$value.name}}(updateData.{{$value.name}}).
		{{ end }}{{ end }}SaveX(ctx)
	return cnt
}

// Delete 删除数据，出错则 panic
func ({{.RecName}} *{{.Name}}RepoImpl) Delete(ctx context.Context, ids... int) int {
	if len(ids) == 0 {
		return 0
	}
	//物理删除
	effectCnt := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Delete().Where({{.EntPackage}}.IDIn(ids...)).ExecX(ctx)

	//软件删除
	// nowTime := int(time.Now().Unix())
	// deleteVal := -1
	// effectCnt := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Update().
	// 	Where({{.EntPackage}}.IDIn(ids...), {{.EntPackage}}.StatusNEQ(deleteVal)).
	//	SetStatus(deleteVal).
	//	SetUpdateTime(nowTime).
	//	SaveX(ctx)
	return effectCnt
}

// SearchList 搜索列表，出错则 panic
func ({{.RecName}} *{{.Name}}RepoImpl) SearchList(ctx context.Context, reqBo *bo.ReqPageBo) (dos []*do.{{.Name}}Do, respPage *bo.RespPageBo) {
	q := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Query()
	{{.RecName}}.SetPageByBo(q, reqBo)
	if reqBo != nil {
		respPage = {{.RecName}}.QueryRespPage(ctx, q, reqBo)
	}
	pos := q.AllX(ctx)
	dos = {{.RecName}}.ToEntities(pos)
	return
}

// GetE 通过 id 获取一条数据
func ({{.RecName}} *{{.Name}}RepoImpl) GetE(ctx context.Context, id int) (*do.{{.Name}}Do, error) {
	row, err := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Query().Where({{.EntPackage}}.ID(id)).First(ctx)
	if err != nil {
		return nil, err
	}
	return {{.RecName}}.ToEntity(row), nil
}

// FindE 通过多个 id 获取多条数据
func ({{.RecName}} *{{.Name}}RepoImpl) FindE(ctx context.Context, ids... int) ([]*do.{{.Name}}Do, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Query().Where({{.EntPackage}}.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, err
	}
	return {{.RecName}}.ToEntities(rows), nil
}

// CreateE 创建数据
func ({{.RecName}} *{{.Name}}RepoImpl) CreateE(ctx context.Context, creatData *do.{{.Name}}Do) (*do.{{.Name}}Do, error) {
	row, err := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Create().
		{{ range $index, $value := .Fields }}{{ if and (ne $value.name "ID") (ne $value.name "Edges")}}Set{{$value.name}}(creatData.{{$value.name}}).
		{{ end }}{{ end }}Save(ctx)
	if err != nil {
		return nil, err
	}
	return {{.RecName}}.ToEntity(row), nil
}

// CreateBulkE 批量创建数据
func ({{.RecName}} *{{.Name}}RepoImpl) CreateBulkE(ctx context.Context, dos []*do.{{.Name}}Do) ([]*do.{{.Name}}Do, error) {
	if len(dos) == 0 {
		return nil, nil
	}
	values := make([]*ent.{{.Name}}Create, len(dos))
	for i, item := range dos {
		values[i] = {{.RecName}}.data.GetDb(ctx).{{.Name}}.Create(){{ range $index, $value := .Fields }}{{ if and (ne $value.name "ID") (ne $value.name "Edges")}}.
			Set{{$value.name}}(item.{{$value.name}}){{ end }}{{ end }}
	}
	rows, err := {{.RecName}}.data.GetDb(ctx).{{.Name}}.CreateBulk(values...).Save(ctx)
	if err != nil {
		return nil, err
	}
	return {{.RecName}}.ToEntities(rows), nil
}

// UpdateE 更新数据
func ({{.RecName}} *{{.Name}}RepoImpl) UpdateE(ctx context.Context, updateData *do.{{.Name}}Do) (int, error) {
	cnt, err := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Update().Where({{.EntPackage}}.ID(updateData.ID)).
		{{ range $index, $value := .Fields }}{{ if and (ne $value.name "ID") (ne $value.name "Edges")}}Set{{$value.name}}(updateData.{{$value.name}}).
		{{ end }}{{ end }}Save(ctx)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// DeleteE 删除数据
func ({{.RecName}} *{{.Name}}RepoImpl) DeleteE(ctx context.Context, ids... int) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	
	// 软件删除
	return {{.RecName}}.data.GetDb(ctx).{{.Name}}.Update().
		Where({{.EntPackage}}.IDIn(ids...), {{.EntPackage}}.DeletedAt(0)).
		SetDeletedAt(int(time.Now().Unix())).
		Save(ctx)
}

// DeleteForceE 删除数据(硬删除)
func ({{.RecName}} *{{.Name}}RepoImpl) DeleteForceE(ctx context.Context, ids... int) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	// 物理删除
	return {{.RecName}}.data.GetDb(ctx).{{.Name}}.Delete().Where({{.EntPackage}}.IDIn(ids...)).Exec(ctx)
}

// SearchListE 搜索列表
func ({{.RecName}} *{{.Name}}RepoImpl) SearchListE(ctx context.Context, reqBo *bo.ReqPageBo) (dos []*do.{{.Name}}Do, respPage *bo.RespPageBo, err error) {
	q := {{.RecName}}.data.GetDb(ctx).{{.Name}}.Query()
	{{.RecName}}.SetPageByBo(q, reqBo)
	if reqBo != nil {
		respPage = {{.RecName}}.QueryRespPage(ctx, q, reqBo)
	}
	pos, err := q.All(ctx)
	if err != nil {
		return nil, nil, err
	}
	dos = {{.RecName}}.ToEntities(pos)
	return
}
`

func Generate() {
	if !initParams() {
		return
	}
	//解析代码树
	po := parsePoSpec()
	poName := po.Name.String()
	fields := getPoFields(po)

	//生成 DO
	generateDo(poName, fields)

	//生成 repo
	generateRepo(poName)

	//生成 impl
	generateImpl(poName, fields)

	//写入注入 impl 方法
	writeProviderSetImpl(poName)
}

func initParams() bool {
	if tableName == "" {
		fmt.Println("table参数不能为空")
		return false
	}
	poFileName = strings.ToLower(tableName)
	poFileName = strings.ReplaceAll(poFileName, "_", "")
	poFilePath = "./internal/data/ent/" + poFileName + ".go"
	// 文件是否存在
	_, err := os.Stat(poFilePath)
	if err != nil {
		fmt.Println("模型文件不存在：" + poFilePath + "请先生成后重试，make ent t=" + tableName)
		return false
	}
	//解析module 名称
	fileName := "go.mod"
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	parsedModFile, err := modfile.Parse("go.mod", fileContent, nil)
	if err != nil {
		panic(err)
	}
	modName = parsedModFile.Module.Mod.Path
	return true
}

func generateDo(poName string, fields []map[string]string) {
	//创建文件，如果存在则不覆盖
	filePath := "./internal/biz/do/" + tableName + ".go"

	// 判断 filePath 文件是否存在
	_, err := os.Stat(filePath)
	if err == nil && !isOverwrite {
		fmt.Println("跳过创建 Do，文件已存在，请先删除\n\trm -rf " + filePath)
		return
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("do").Parse(doTemplate)
	if err != nil {
		panic(err)
	}
	dataDo := struct {
		Fields []map[string]string
		Name   string
	}{
		Fields: fields,
		Name:   poName,
	}
	err = tmpl.Execute(f, dataDo)
	if err != nil {
		fmt.Println("生成Do代码出错：", err)
	}
	fullPath, _ := filepath.Abs(filePath)
	fmt.Println("创建成功：file://" + fullPath)
}

func generateRepo(poName string) {
	filePath := "./internal/biz/repository/" + tableName + ".go"
	// 判断 filePath 文件是否存在
	_, err := os.Stat(filePath)
	if err == nil && !isOverwrite {
		fmt.Println("跳过创建 Repo，文件已存在，请先删除\n\trm -rf " + filePath)
		return
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("repo").Parse(repoTemplate)
	if err != nil {
		panic(err)
	}

	data := struct {
		Name    string
		RecName string
		ModName string
	}{
		Name:    poName,
		RecName: strings.ToLower(poName[0:1]),
		ModName: modName,
	}
	err = tmpl.Execute(f, data)
	if err != nil {
		fmt.Println("生成Do代码出错：", err)
	}
	fullPath, _ := filepath.Abs(filePath)
	fmt.Println("创建成功：file://" + fullPath)
}

func generateImpl(poName string, fields []map[string]string) {
	filePath := "./internal/data/repositoryimpl/" + tableName + ".go"
	// 判断 filePath 文件是否存在
	_, err := os.Stat(filePath)
	if err == nil && !isOverwrite {
		fmt.Println("跳过创建 Impl，文件已存在，请先删除\n\trm -rf " + filePath)
		return
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("repo").Parse(repoImplTemplate)
	if err != nil {
		panic(err)
	}

	data := struct {
		Methods    []map[string]string
		Name       string
		RecName    string
		EntPackage string
		Fields     []map[string]string
		ModName    string
	}{
		Name:       poName,
		RecName:    strings.ToLower(poName[0:1]),
		EntPackage: strings.ToLower(strings.ReplaceAll(poName, "_", "")),
		Fields:     fields,
		ModName:    modName,
	}
	err = tmpl.Execute(f, data)
	if err != nil {
		fmt.Println("生成Impl代码出错：", err)
	}
	fullPath, _ := filepath.Abs(filePath)
	fmt.Println("创建成功：file://" + fullPath)
}

func writeProviderSetImpl(poName string) {
	setPath := "./internal/data/repositoryimpl/provider_set.go"
	fullSetPath, _ := filepath.Abs(setPath)
	content, err := os.ReadFile(fullSetPath)
	if err != nil {
		fmt.Printf("写入impl注入代码到 file://%s 失败，读取失败，请确定文件是否存在\n", fullSetPath)
		panic(err)
	}
	code := string(content)
	methodName := "New" + poName + "RepoImpl"
	if strings.Contains(code, methodName) {
		fmt.Printf("跳过写入 impl 注入代码，已经存在, file://%s\n", fullSetPath)
		return
	}
	writeCode := `wire.NewSet(
	` + methodName + `,
`

	newCode := strings.Replace(code, "wire.NewSet(", writeCode, 1)
	err = os.WriteFile(setPath, []byte(newCode), 0600)
	if err != nil {
		fmt.Printf("写入注入代码文件 file://%s 失败\n", fullSetPath)
		panic(err)
	}
	fmt.Println("已写入注入代码 NewMemberRepoImpl 到 file://" + fullSetPath)

}

func parsePoSpec() *ast.TypeSpec {
	fs := token.NewFileSet()
	// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	path, _ := filepath.Abs(poFilePath)

	f, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	po := f.Decls[1].(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
	return po
}

func getPoFields(poTs *ast.TypeSpec) []map[string]string {
	var m []map[string]string
	poStructType := poTs.Type.(*ast.StructType)
	for _, field := range poStructType.Fields.List {
		t, isOk := field.Type.(*ast.Ident)
		if field.Names == nil || !isOk {
			continue
		}
		item := map[string]string{
			"name": field.Names[0].Name,
			"type": t.Name,
		}
		m = append(m, item)
	}
	return m
}
