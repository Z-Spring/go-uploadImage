package upload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"upload2/global"
)

//var AppSetting *settings.AppSettings

type FileInfo struct {
	FileName string
	FileUrl  string
}

func ImageUpload(fileHeader *multipart.FileHeader, file multipart.File) (*FileInfo, error) {
	fileName := GetFileName(fileHeader.Filename)
	//fileName = Md5(fileName)
	if !CheckImageExt(fileName) {
		return nil, errors.New("图片格式不支持")
	}
	if !CheckMaxSize(file) {
		return nil, errors.New("图片太大，请压缩后再上传")
	}
	storePath := GetSavePath()
	if err := CreateSavePath(storePath, os.ModePerm); err != nil {
		return nil, errors.New("创建文件路径失败")
	}
	dst := storePath + "/" + fileName
	if err := CreateFile(fileHeader, dst); err != nil {
		return nil, errors.New("创建文件失败")
	}
	fileUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{fileName, fileUrl}, nil
}

func Md5(value string) string {
	m := md5.New()
	_, err := io.WriteString(m, value)
	if err != nil {
		log.Println("md5 error: ", err)
	}
	return hex.EncodeToString(m.Sum(nil))
}

func GetFileName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = Md5(fileName)
	return fileName + ext
}

func GetSavePath() string {
	return global.AppSetting.UploadStoreUrl
}

func CheckMaxSize(f multipart.File) bool {
	all, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	size := len(all)
	// todo 1
	if size >= global.AppSetting.UploadMaxSize*1024*1024 {
		return false
	}
	return true
}

func CheckImageExt(fileName string) bool {
	ext := path.Ext(fileName)
	ext = strings.ToUpper(ext)
	for _, suffix := range global.AppSetting.UploadPageSuffix {
		if strings.ToUpper(suffix) == ext {
			return true
		}
	}
	return false
}

func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		log.Printf("error: %s", err)
	}
	return nil
}

func CreateFile(fileHeader *multipart.FileHeader, dst string) error {
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		log.Println("openFile error", err)
	}
	create, err := os.Create(dst)
	defer create.Close()
	if err != nil {
		log.Println("create file error", err)
	}
	_, err = io.Copy(create, file)
	return err
}
