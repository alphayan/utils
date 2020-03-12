package xcjswag

import (
	"encoding/json"
	"github.com/alphayan/iris"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
	"utils/config"
	"utils/dir"
)

const (
	SRCPKG       = "src" //源码默认放在src下边
	PROPATH      = "./"
	CTRPKG       = "controller" //默认扫描controller包,此处现在废弃不用
	UtilResponse = "utils/response"
)

// Gen presents a generate tool for swag.
type Gen struct {
}

func InitSwagger(app *iris.Application) {
	os.Getenv("GOPATH")
	//TODO 扫描包需要修改 PROPATH这样写可能会有问题，这里尝试没有问题
	//PROPATH, _ := filepath.Abs("./")
	if config.Config.Swagger.Controller == "" {
		config.Config.Swagger.Controller = CTRPKG
	}
	NewGen().Build(config.Config.Swagger.Apipackage, filepath.Join(PROPATH, SRCPKG, config.Config.Swagger.Mainpackage))
}

// New creates a new Gen.
func NewGen() *Gen {
	return &Gen{}
}

// Build builds swagger json file  for gived searchDir and mainApiFile.
func (g *Gen) Build(searchDir, mainApiFile string) error {
	log.Println("Generate swagger docs....")
	p := New()
	p.ParseApi(searchDir, mainApiFile)
	swagger := p.GetSwagger()

	b, _ := json.MarshalIndent(swagger, "", "    ")
	//flag := dir.IsExist(config.Config.SwaggerPath)
	flag := dir.IsExist(config.Config.Swagger.Swaggerpath)
	if !flag {
		dir.MkdirAll(config.Config.Swagger.Swaggerpath)
	}
	//os.MkdirAll(filepath.Join(searchDir, "docs"), os.ModePerm)
	docs, _ := os.Create(filepath.Join(config.Config.Swagger.Swaggerpath, "swagger.json"))
	defer docs.Close()

	packageTemplate.Execute(docs, struct {
		Timestamp time.Time
		Doc       string
	}{
		Timestamp: time.Now(),
		//Doc:       "`" + string(b) + "`",
		Doc: string(b),
	})

	log.Printf("create swagger.json at  %+v", docs.Name())
	return nil
}

//Json文件不支持注解
var packageTemplate = template.Must(template.New("").Parse(`{{.Doc}}`))
