package service

import (
	"context"
	"encoding/json"
	"github.com/Etpmls/EM-Attachment/src/application"
	"github.com/Etpmls/EM-Attachment/src/application/model"
	"github.com/Etpmls/EM-Attachment/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"net/http"
)

type ServiceAttachment struct {
	protobuf.UnimplementedAttachmentServer
}

// Upload Image
// 上传图片
// Upload pictures via http, only save the picture path to the database
// 通过http上传图片，仅保存图片路径到数据库
func (this ServiceAttachment) UploadImage(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	// Get file from request
	f, file, err := r.FormFile("file")
	if err != nil {
		return
	}

	f.Close()

	// Validate Image and get extension
	var attachment model.Attachment
	extension, err := attachment.AttachmentValidateImage(file)
	if err != nil {
		em.LogDebug.Output(em.MessageWithLineNum(err.Error()))
		b, _ := em.ErrorHttp(em.ERROR_Code, "Image verification failed!", nil, err)
		w.Write(b)
		return
	}

	path, err := attachment.AttachmentUploadImage(file, extension)
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		b, _ := em.ErrorHttp(em.ERROR_Code, "Image upload failed!", nil, err)
		w.Write(b)
		return
	}

	b, err := json.Marshal(map[string]string{"path" : application.ServiceConfig.Service.Host + path})
	if err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
		b, _ := em.ErrorHttp(em.ERROR_Code, "Image upload failed!", nil, err)
		w.Write(b)
		return
	}

	b2, _ := em.SuccessHttp(em.SUCCESS_Code, "Upload success", string(b))
	w.Write(b2)
}

// Create attachment link
// 创建附件关联
// The relationship between the distribution path and the service
// 分配路径与所属服务的关联关系
type validate_AttachmentCreate struct {
	Service string	`json:"service" validate:"required,max=255"`
	Path string	`json:"path" validate:"required"`
	OwnerID uint	`json:"owner_id" validate:"required"`
	OwnerType string	`json:"owner_type" validate:"required,max=255"`
}
func (this *ServiceAttachment) Create(ctx context.Context, request *protobuf.AttachmentCreate) (*em_protobuf.Response, error) {
	// Validate
	var vd validate_AttachmentCreate
	err := em.Validator.Validate(request, &vd)
	if err != nil {
		em.LogWarn.Output(em.MessageWithLineNum_OneRecord(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
	}

	// Set storage method
	request.StorageMethod = application.ServiceConfig.Service.FileStorageMethod
	var attachment model.Attachment
	request.Path = attachment.GetUrlPath(request.Path)

	err = em.DB.Transaction(func(tx *gorm.DB) error {
		var old model.Attachment
		result := tx.Where("service = ?", request.Service).Where("owner_id = ?", request.OwnerId).Where("owner_type = ?", request.OwnerType).First(&old)

		// If the form contains attachment path
		// 如果表单包含附件路径
		if len(request.Path) > 0 {
			// Same as the previous attachment path, skip
			// 和以前附件路径一样，跳过
			if old.Path == request.Path {
				return nil
			} else {
				// Delete old attachments, update new attachments
				// 删除旧附件，更新新附件
				if result.RowsAffected > 0 {
					// Delete attachments and databases according to Path
					// 根据Path删除附件和数据库
					err := old.AttachmentBatchDelete([]string{old.Path})
					if err != nil {
						em.LogError.Output(em.MessageWithLineNum(err.Error()))
						return err
					}
				}

				m, err := em.StructToMap(request)
				if err != nil {
					em.LogError.Output(em.MessageWithLineNum(err.Error()))
					return err
				}
				// Note: json to map int format will be converted to float
				// 注意：json转map int格式会转换为float
				m["owner_id"] = uint(m["owner_id"].(float64))

				result := tx.Model(model.Attachment{}).Where("path = ?", request.Path).Updates(m)
				if result.Error != nil {
					return result.Error
				}
			}

		} else {
			if result.RowsAffected > 0 {
				// Delete attachments and databases according to Path
				// 根据Path删除附件和数据库
				err := old.AttachmentBatchDelete([]string{old.Path})
				if err != nil {
					em.LogError.Output(em.MessageWithLineNum(err.Error()))
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), nil)
}
/*func (this *ServiceAttachment) Create(ctx context.Context, request *protobuf.AttachmentCreate) (*em_protobuf.Response, error) {
	// Validate
	var vd validate_AttachmentCreate
	err := em.Validator.Validate(request, &vd)
	if err != nil {
		em.LogWarn.Output(em.MessageWithLineNum_OneRecord(err.Error()))
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
	}

	// Set storage method
	request.StorageMethod = application.ServiceConfig.Service.FileStorageMethod
	var attachment model.Attachment
	request.Path = attachment.GetUrlPath(request.Path)

	err = em.DB.Transaction(func(tx *gorm.DB) error {
		// Delete historical thumbnails
		// 则删除历史缩略图
		var old model.Attachment
		result := tx.Where("service = ?", request.Service).Where("owner_id = ?", request.OwnerId).Where("owner_type = ?", request.OwnerType).First(&old)
		// Delete if found
		// 如果找到记录则删除
		if result.RowsAffected > 0 {
			// Delete attachments and databases according to Path
			// 根据Path删除附件和数据库
			err := old.AttachmentBatchDelete([]string{old.Path})
			if err != nil {
				em.LogError.Output(em.MessageWithLineNum(err.Error()))
				return err
			}
		}

		// If the form contains thumbnails
		// 如果表单包含缩略图，
		if len(request.Path) > 0 {
			m, err := em.StructToMap(request)
			if err != nil {
				em.LogError.Output(em.MessageWithLineNum(err.Error()))
				return err
			}
			// Note: json to map int format will be converted to float
			// 注意：json转map int格式会转换为float
			m["owner_id"] = uint(m["owner_id"].(float64))

			result := tx.Model(model.Attachment{}).Where("path = ?", attachment.GetUrlPath(request.Path)).Updates(m)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), nil)
}*/

type validate_AttachmentGetOne struct {
	Service string	`json:"service" validate:"required,max=255"`
	OwnerID uint	`json:"owner_id" validate:"required"`
	OwnerType string	`json:"owner_type" validate:"required,max=255"`
}
func (this *ServiceAttachment) GetOne(ctx context.Context, request *protobuf.AttachmentGetOne) (*em_protobuf.Response, error) {
	// Validate
	{
		var vd validate_AttachmentGetOne
		err := em.Validator.Validate(request, &vd)
		if err != nil {
			em.LogWarn.Output(em.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	var a model.Attachment
	result := em.DB.Where("service = ?", request.GetService()).Where("owner_id = ?", request.GetOwnerId()).Where("owner_type = ?", request.GetOwnerType()).First(&a)
	if result.RowsAffected == 0 {
		em.LogInfo.Output(em.MessageWithLineNum("No record"))
	}

	// If it is stored locally, the path plus the domain name
	// 如果是本地储存，则路径加上域名
	a.MakeUrlPath(&a)

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), a)
}

type validate_AttachmentDiskCleanUp struct {
	Service string	`json:"service" validate:"required,max=255"`
}
func (this *ServiceAttachment) DiskCleanUp(ctx context.Context, request *protobuf.AttachmentDiskCleanUp) (*em_protobuf.Response, error) {
	{
		// Validate
		var vd validate_AttachmentDiskCleanUp
		err := em.Validator.Validate(request, &vd)
		if err != nil {
			em.LogWarn.Output(em.MessageWithLineNum(err.Error()))
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	var attachment model.Attachment
	err := attachment.DeleteUnused(request.GetService())
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Clear"), nil, err)
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Clear"), nil)
}