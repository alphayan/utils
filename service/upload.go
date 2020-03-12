package service

import (
	"errors"
	"fmt"
	"github.com/alphayan/iris"
	"github.com/mholt/archiver"
	"github.com/rs/xid"
	"github.com/ungerik/go-dry"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"github.com/alphayan/utils/config"
	"github.com/alphayan/utils/dir"
)

//设置默认路径
var (
	BathPath     string = "static"
	UploadPath   string = "uploads"   // 小车匠资源
	QuestionPath string = "questions" // 小车匠试题
	TrianPath    string = "trains"    // 培训资料
	TempPath     string = "temp"      // 导出成绩模板
	ResumePath   string = "resume"    //师傅简历
	PracticePath   string = "practice"    //任务提交图片或视频资源
	PhotoPath   string = "photo"    //任务提交图片或视频资源
)

// 将文件存储到服务器
func SaveFile(file multipart.File, Filename string, ExtentName string, Flag string) map[string]interface{} {
	guid := xid.New()
	var filePath string
	// 一级、二级目录
	if Flag == "question" { // 小车匠试题
		filePath = filepath.Join(config.Config.StaticPath, QuestionPath)
	} else if Flag == "resource" { // 小车匠资源
		filePath = filepath.Join(config.Config.StaticPath, UploadPath)
	} else if Flag == "train" { // 培训资料
		filePath = filepath.Join(config.Config.StaticPath, TrianPath)
	}else if Flag == "practice" { //上传任务资源
		filePath = filepath.Join(config.Config.StaticPath, PracticePath)
	}else if Flag == "photo" { //上传任务资源
		filePath = filepath.Join(config.Config.StaticPath, PhotoPath)
	}else if Filename != "" { //上传师傅简历
		filePath = filepath.Join(config.Config.StaticPath, ResumePath)
	}
	// 三级目录：按文件上传日期生成
	fileSubPath := guid.Time().Format("2006-01")
	// 四级目录：按文件上传生成的UID生成
	fourthDir := guid.String()
	// 文件上传的最终目录
	fileDir := filepath.Join(filePath, fileSubPath, fourthDir)
	flag := dir.IsExist(fileDir)
	if flag == false {
		dir.MkdirAll(fileDir)
	}

	var convertUrl string // 需要转换的路径
	var convertUrl2 string
	var resourceUrl string
	var resourceUrl2 string
	var fileUrl string
	var fileName string
	var reviewUrl string
	if Flag == "question" { // 小车匠试题
		convertUrl = path.Join(BathPath, QuestionPath, fileSubPath, fourthDir)

		convertUrl2 = filepath.Join(QuestionPath, fileSubPath, fourthDir)
		// path.Join /
		resourceUrl = path.Join(convertUrl, Filename)

		resourceUrl2 = filepath.Join(convertUrl2, Filename)
		// filepath.Join \
		fileUrl = filepath.Join(fileDir, Filename)
	} else if Flag == "resource" {
		convertUrl = path.Join(BathPath, UploadPath, fileSubPath, fourthDir)
		convertUrl2 = filepath.Join(UploadPath, fileSubPath, fourthDir)
		// path.Join /
		guidName := fmt.Sprintf("%s%s", guid.String(), ".")
		fileName = fmt.Sprintf("%s%s", guidName, ExtentName)
		resourceUrl = path.Join(convertUrl, fileName)
		resourceUrl2 = filepath.Join(convertUrl2, fileName)
		// 转换后的路径
		pdf := fmt.Sprintf("%s%s", guidName, "pdf")
		reviewUrl = path.Join(convertUrl, pdf)

		// filepath.Join \
		fileUrl = filepath.Join(fileDir, fileName)
	} else if Flag == "train" {
		convertUrl = path.Join(BathPath, TrianPath, fileSubPath, fourthDir)
		convertUrl2 = filepath.Join(TrianPath, fileSubPath, fourthDir)
		// path.Join /
		guidName := fmt.Sprintf("%s%s", guid.String(), ".")
		fileName = fmt.Sprintf("%s%s", guidName, ExtentName)
		resourceUrl = path.Join(convertUrl, fileName)
		resourceUrl2 = filepath.Join(convertUrl2, fileName)
		// 转换后的路径
		pdf := fmt.Sprintf("%s%s", guidName, "pdf")
		reviewUrl = path.Join(convertUrl, pdf)

		// filepath.Join \
		fileUrl = filepath.Join(fileDir, fileName)
	} else if Flag == "practice"  {//实习任务资源
		convertUrl = path.Join(BathPath, PracticePath, fileSubPath, fourthDir)
		convertUrl2 = filepath.Join(PracticePath, fileSubPath, fourthDir)
		// path.Join /
		guidName := fmt.Sprintf("%s%s", guid.String(), ".")
		fileName = fmt.Sprintf("%s%s", guidName, ExtentName)
		resourceUrl = path.Join(convertUrl, fileName)
		resourceUrl2 = filepath.Join(convertUrl2, fileName)
		// 转换后的路径
		pdf := fmt.Sprintf("%s%s", guidName, "pdf")
		reviewUrl = path.Join(convertUrl, pdf)

		// filepath.Join \
		fileUrl = filepath.Join(fileDir, fileName)
	}else if Flag == "photo"  {//上传头像
		convertUrl = path.Join(BathPath, PhotoPath, fileSubPath, fourthDir)
		convertUrl2 = filepath.Join(PhotoPath, fileSubPath, fourthDir)
		// path.Join /
		guidName := fmt.Sprintf("%s%s", guid.String(), ".")
		fileName = fmt.Sprintf("%s%s", guidName, ExtentName)
		resourceUrl = path.Join(convertUrl, fileName)
		resourceUrl2 = filepath.Join(convertUrl2, fileName)
		// 转换后的路径
		pdf := fmt.Sprintf("%s%s", guidName, "pdf")
		reviewUrl = path.Join(convertUrl, pdf)

		// filepath.Join \
		fileUrl = filepath.Join(fileDir, fileName)
	}else if Filename != "" {
		convertUrl = path.Join(BathPath, ResumePath, fileSubPath, fourthDir)
		convertUrl2 = filepath.Join(ResumePath, fileSubPath, fourthDir)
		// path.Join /
		guidName := fmt.Sprintf("%s%s", guid.String(), ".")
		fileName = fmt.Sprintf("%s%s", guidName, ExtentName)
		resourceUrl = path.Join(convertUrl, fileName)
		resourceUrl2 = filepath.Join(convertUrl2, fileName)
		// 转换后的路径
		pdf := fmt.Sprintf("%s%s", guidName, "pdf")
		reviewUrl = path.Join(convertUrl, pdf)

		// filepath.Join \
		fileUrl = filepath.Join(fileDir, fileName)
	}
	out, err := os.OpenFile(fileUrl, os.O_WRONLY|os.O_CREATE, 0771)
	if err != nil {
		return map[string]interface{}{"Status": iris.StatusBadRequest, "Error: ": err}
	}
	defer out.Close()
	_, err = io.Copy(out, file) // 存储到服务器
	if err != nil {
		return map[string]interface{}{"Status": iris.StatusBadRequest, "Error: ": err}
	}

	if ExtentName == "zip" { // 小车匠仿真题
		info := doXcjSimulation(resourceUrl, convertUrl)
		if info["Status"].(int) == iris.StatusOK {
			return info
		} else {
			return map[string]interface{}{"Status": iris.StatusBadRequest, "Info": info}
		}
	} else {
		return map[string]interface{}{"Status": iris.StatusOK, "ResourceUrl": resourceUrl,
			"ConvertUrl": convertUrl, "FileName": fileName, "ReviewUrl": reviewUrl, "ResourceUrl2": resourceUrl2,
			"ConvertUrl2": convertUrl2}
	}
}

func doXcjSimulation(resourceUrl string, destinationUrl string) map[string]interface{} {
	basePath := config.Config.StaticPath // 项目路径
	resourceUrls := strings.Split(resourceUrl, "static/")
	destinationUrls := strings.Split(destinationUrl, "static/")
	sour := filepath.Join(basePath, resourceUrls[1])
	dest := filepath.Join(basePath, destinationUrls[1])
	log.Println("sour = ", sour)
	log.Println("dest = ", dest)
	err := archiver.Zip.Open(sour, dest)
	//err := archiver.Zip.Open(resourceUrl, destinationUrl)
	//err := archiver.Zip.Open("static/questions/2018-01/xcj_003001_001.zip", "static/questions/2018-01")
	if err != nil {
		return map[string]interface{}{"Status": iris.StatusBadRequest, "Error: ": err}
	}
	child, _ := dry.ListDirDirectories(dest)
	if len(child) == 0 {
		err = errors.New("仿真题错误")
		return map[string]interface{}{"Status": iris.StatusBadRequest, "Error: ": err.Error()}
	}
	ReviewUrl := path.Join(destinationUrl, child[0])
	return map[string]interface{}{"Status": iris.StatusOK, "ResourceUrl": resourceUrl, "ReviewUrl": ReviewUrl}
}

// 存储字节流文件
func SaveBytes(bodyBytes []byte, Flag string, fileUrl string) map[string]interface{} {
	//BathPath := config.Config.StaticPath
	guid := xid.New()
	var filePath string
	// 一级、二级目录
	if Flag == "question" { // 小车匠试题
		filePath = filepath.Join(config.Config.StaticPath, QuestionPath)
	} else if Flag == "resource" { // 小车匠资源
		filePath = filepath.Join(config.Config.StaticPath, UploadPath)
	} else if Flag == "train" { // 培训资料
		filePath = filepath.Join(config.Config.StaticPath, TrianPath)
	}
	// 三级目录：按文件上传日期生成
	fileSubPath := guid.Time().Format("2006-01")
	// 四级目录：按文件上传生成的UID生成
	fileUrls := strings.Split(fileUrl, "/")
	fileNames := fileUrls[len(fileUrls)-1]
	Name := strings.Split(fileNames, ".")
	// 文件上传的最终目录
	fileDir := filepath.Join(filePath, fileSubPath, Name[0])
	flag := dir.IsExist(fileDir)
	if flag == false {
		dir.MkdirAll(fileDir)
	}
	urlFile := filepath.Join(fileDir, fileNames)
	err := ioutil.WriteFile(urlFile, bodyBytes, 0666) //写入文件(字节数组)

	var reviewUrl string
	if Flag == "question" { // 小车匠试题
		reviewUrl = path.Join("/", BathPath, QuestionPath, fileSubPath, Name[0], fileNames)
	} else if Flag == "resource" {
		reviewUrl = path.Join("/", BathPath, UploadPath, fileSubPath, Name[0], fileNames)
	} else if Flag == "train" {
		reviewUrl = path.Join("/", BathPath, TrianPath, fileSubPath, Name[0], fileNames)
	}
	if err != nil {
		return map[string]interface{}{"Status": iris.StatusBadRequest, "Err": err, "ReviewUrl": reviewUrl}
	}
	return map[string]interface{}{"Status": iris.StatusOK, "Err": nil, "ReviewUrl": reviewUrl}
}

func GetMediaFullPathWithMonth(filename string) string {
	fileSubPath := time.Now().Format("2006-01")
	fileDir := filepath.Join(BathPath, TempPath, fileSubPath)
	flag := dir.IsExist(fileDir)
	if flag == false {
		dir.MkdirAll(fileDir)
	}
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		os.Mkdir(fileDir, 0771)
	}
	return filepath.Join(fileDir, filename)
}
