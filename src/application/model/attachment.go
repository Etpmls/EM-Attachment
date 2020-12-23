package model

import (
	"errors"
	"github.com/Etpmls/EM-Attachment/src/application"
	em "github.com/Etpmls/Etpmls-Micro"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type Attachment struct {
	ID        uint `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Service string	`json:"service"`
	StorageMethod string	`json:"storage_method"`
	Path string	`json:"path"`
	OwnerID uint	`json:"owner_id"`
	OwnerType string	`json:"owner_type"`
}

// Validate Path is a file in storage/upload
// 严重路径是否在storage/upload中
func (this *Attachment) AttachmentValidatePath(path string) error {
	const upload_path = "storage/upload/"
	// 截取前十五个字符判断和Path是否相同
	if len(path) <= len(upload_path) || !strings.Contains(path[:len(upload_path)], upload_path) {
		em.LogError.Output(em.MessageWithLineNum("Illegal request path!"))
		return  errors.New("Illegal request path!")
	}
	// 删除前缀
	// p := strings.TrimPrefix(path, upload_path)
	f, err := os.Stat(path)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	// 如果文件是目录
	if f.IsDir() {
		em.LogError.Output(em.MessageWithLineNum("Cannot delete directory!"))
		return errors.New("Cannot delete directory!")
	}
	return nil
}

// Batch delete any type of files in storage/upload/
// 批量删除storage/upload/中的任何类型文件
func (this *Attachment) AttachmentBatchDelete(s []string) (err error) {
	for _, v := range s {
		// Validate If a File
		err = this.AttachmentValidatePath(v)
		if err != nil {
			return err
		}
		// Delete Image
		_ = os.Remove(v)
	}

	// Delete Database
	if err = em.DB.Unscoped().Where("path IN (?)", s).Delete(Attachment{}).Error; err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return err
	}

	return err
}

// Delete unused attachments
// 删除未使用的附件
func (this *Attachment) DeleteUnused(service_name string) error {
	var a []Attachment
	em.DB.Where("service = ?", service_name).Where("owner_id = ?", 0).Or("owner_type = ?", "").Find(&a)

	// If there is no value, return directly
	// 如果没有值，则直接返回
	if len(a) == 0 {
		em.LogDebug.Output(em.MessageWithLineNum("No files need to be deleted!"))
		return nil
	}

	var pt []string
	for _, v := range a {
		pt = append(pt, v.Path)
	}

	err := this.AttachmentBatchDelete(pt)
	if err != nil {
		return err
	}

	return nil
}

// Validate if file is a image
// 验证文件是否为图片
func (this *Attachment) AttachmentValidateImage(h *multipart.FileHeader) (s string, err error) {
	f, err := h.Open()
	if err != nil {
		return s, err
	}

	// 识别图片类型
	_, image_type, _ := image.Decode(f)

	// 获取图片的类型
	switch image_type {
	case `jpeg`:
		return "jpeg", nil
	case `png`:
		return "png", nil
	case `gif`:
		return "git", nil
	case `bmp`:
		return "bmp", nil
	default:
		err := errors.New("This is not an image file, or the image file format is not supported!")
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return "", err
	}
}

// Upload Image
// 上传图片
func (this *Attachment) AttachmentUploadImage(file *multipart.FileHeader, extension string) (p string, err error) {
	// Make Dir
	t := time.Now().Format("20060102")
	path := "storage/upload/" + t + "/"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return p, err
	}
	// UUID File name
	u := strings.ReplaceAll(uuid.New().String(), "-", "")

	file_path := path + u + "." + extension
	err = this.SaveUploadedFile(file, file_path)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return p, errors.New("Failed to save file!")
	}

	if err = em.DB.Create(&Attachment{Path: file_path}).Error; err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		return p, err
	}

	return file_path, err
}

// Copy GIN c.SaveUploadedFile
func (this *Attachment) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// If it is stored locally, the path plus the domain name
// 如果是本地储存，则路径加上域名
func (this *Attachment) MakeUrlPath(attachment *Attachment) {
	if attachment.Path == "" || application.ServiceConfig.Service.FileStorageMethod != "local" {
		return
	}

	attachment.Path = application.ServiceConfig.Service.Host + attachment.Path
	return
}

// If it is stored locally, save the path without host, if it is not stored locally, save the full URL path
// 如果是本地储存，则保存不带host的路径，若非本地储存，则保存完整url路径
func (this *Attachment) GetUrlPath(urlPath string) string {
	if application.ServiceConfig.Service.FileStorageMethod == "local" {
		return em.GetUrlPath(urlPath, true)
	}
	return urlPath
}
