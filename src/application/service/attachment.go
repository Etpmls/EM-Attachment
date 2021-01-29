package service

import (
	"context"
	"encoding/json"
	"github.com/Etpmls/EM-Attachment/src/application"
	"github.com/Etpmls/EM-Attachment/src/application/model"
	"github.com/Etpmls/EM-Attachment/src/application/protobuf"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/v2/protobuf"
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
		em.LogError.OutputSimplePath(err.Error())
		b, _ := em.ErrorHttp(em.ERROR_Code, "Image upload failed!", nil, err)
		w.Write(b)
		return
	}

	k, err := em.Kv.ReadKey(define.MakeServiceConfField(em.Micro.Config.Service.RpcName, application.KvServiceHost))
	if err != nil {
		em.LogError.OutputSimplePath(err)
	}
	b, err := json.Marshal(map[string]string{"path" : k + path})
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
	fsm, err := em.Kv.ReadKey(define.MakeServiceConfField(em.Micro.Config.Service.RpcName, application.KvServiceFileStorageMethod))
	if err != nil {
		em.LogError.OutputSimplePath(err)
	}
	request.StorageMethod = fsm
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
						em.LogError.OutputSimplePath(err.Error())
						return err
					}
				}

				m, err := em.StructToMap(request)
				if err != nil {
					em.LogError.OutputSimplePath(err.Error())
					return err
				}
				// Note: json to map int format will be converted to float
				// 注意：json转map int格式会转换为float
				m["owner_id"] = uint(m["owner_id"].(float64))

				result := tx.Model(model.Attachment{}).Where("path = ?", request.Path).Updates(m)
				if result.Error != nil {
					em.LogError.OutputSimplePath(result.Error.Error())
					return result.Error
				}
			}

		} else {
			if result.RowsAffected > 0 {
				// Delete attachments and databases according to Path
				// 根据Path删除附件和数据库
				err := old.AttachmentBatchDelete([]string{old.Path})
				if err != nil {
					em.LogError.OutputSimplePath(err.Error())
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

type validate_AttachmentCreateMany struct {
	Service string	`json:"service" validate:"required,max=255"`
	Paths []string	`json:"paths"`
	OwnerID uint	`json:"owner_id" validate:"required"`
	OwnerType string	`json:"owner_type" validate:"required,max=255"`
}
func (this *ServiceAttachment) CreateMany(ctx context.Context, request *protobuf.AttachmentCreateMany) (*em_protobuf.Response, error) {
	{
		// Validate
		var vd validate_AttachmentCreateMany
		err := em.Validator.Validate(request, &vd)
		if err != nil {
			em.LogInfo.OutputSimplePath(err.Error())
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Set storage method
	fsm, err := em.Kv.ReadKey(define.MakeServiceConfField(em.Micro.Config.Service.RpcName, application.KvServiceFileStorageMethod))
	if err != nil {
		em.LogError.OutputSimplePath(err)
	}

	request.StorageMethod = fsm
	var attachment model.Attachment
	var tmp []string
	for _, v := range request.Paths {
		tmp = append(tmp, attachment.GetUrlPath(v))
	}
	request.Paths = tmp

	err = em.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Find old attachments
		// 1.查找旧附件
		var old []model.Attachment
		tx.Where("service = ?", request.Service).Where("owner_id = ?", request.OwnerId).Where("owner_type = ?", request.OwnerType).Find(&old)

		// 2. Remove the attachments with the same path, and the rest is the different attachments
		// 2.去除path相同的附件，剩下的就是差异的附件
		var same_req = make(map[int]bool)
		var same_old = make(map[int]bool)
		for k, v := range request.Paths {
			for ok, ov := range old {
				// same attachment
				// 相同的附件
				if v == ov.Path {
					same_req[k] = true
					same_old[ok] = true
				}
			}
		}
		// 2-1 Get old attachments
		// 2-1 获取老的附件
		var old_paths []string
		for k, v := range old {
			if same_old[k] != true {
				old_paths = append(old_paths, v.Path)
			}
		}
		// 2-2 Get new attachment
		// 2-2 获取新的附件
		var new_paths []string
		for k, v := range request.Paths {
			if same_req[k] != true {
				new_paths = append(new_paths, v)
			}
		}

		// 3. Delete old attachments
		// 3.删除老的附件
		if len(old_paths) > 0 {
			// Delete attachments and databases according to Path
			// 根据Path删除附件和数据库
			var a model.Attachment
			err := a.AttachmentBatchDelete(old_paths)
			if err != nil {
				em.LogError.OutputSimplePath(err.Error())
				return err
			}
		}

		// 4. If the new attachment list is not empty, add a new attachment
		// 4.如果新附件列表不为空，则增加新的附件
		if len(new_paths) > 0 {
			m, err := em.StructToMap(request)
			if err != nil {
				em.LogError.OutputSimplePath(err.Error())
				return err
			}
			// Note: json to map int format will be converted to float
			// 注意：json转map int格式会转换为float
			m["owner_id"] = uint(m["owner_id"].(float64))
			delete(m, "paths")

			result := tx.Model(model.Attachment{}).Where("path IN ?", new_paths).Updates(m)
			if result.Error != nil {
				em.LogError.OutputSimplePath(result.Error.Error())
				return result.Error
			}
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), nil)
}

type validate_AttachmentDelete struct {
	Service string	`json:"service" validate:"required,max=255"`
	OwnerIds []uint	`json:"owner_ids" validate:"required"`
	OwnerType string	`json:"owner_type" validate:"required,max=255"`
}
func (this *ServiceAttachment) Delete(ctx context.Context, request *protobuf.AttachmentDelete) (*em_protobuf.Response, error) {
	{
		// Validate
		var vd validate_AttachmentDelete
		err := em.Validator.Validate(request, &vd)
		if err != nil {
			em.LogWarn.OutputSimplePath(err.Error())
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// 查找全部attachment
	var list []model.Attachment
	em.DB.Model(&model.Attachment{}).Where("service = ?", request.GetService()).Where("owner_type = ?", request.GetOwnerType()).Where("owner_id IN ?", request.GetOwnerIds()).Find(&list)
	if len(list) > 0 {
		// 根据Path删除附件
		var attachment model.Attachment
		err := attachment.BatchDelete(list)
		if err != nil {
			em.LogError.OutputSimplePath(err)
			return em.ErrorRpc(codes.Internal, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Delete"), nil, err)
		}
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Delete"), nil)
}

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
			em.LogWarn.OutputSimplePath(err)
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	var a model.Attachment
	result := em.DB.Where("service = ?", request.GetService()).Where("owner_id = ?", request.GetOwnerId()).Where("owner_type = ?", request.GetOwnerType()).First(&a)
	if result.RowsAffected == 0 {
		em.LogInfo.OutputSimplePath("No record")
	}

	// If it is stored locally, the path plus the domain name
	// 如果是本地储存，则路径加上域名
	a.MakeUrlPath(&a)

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), a)
}

type validate_AttachmentGetMany struct {
	Service string	`json:"service" validate:"required,max=255"`
	OwnerIds []uint	`json:"owner_ids"`
	OwnerType string	`json:"owner_type" validate:"required,max=255"`
}
func (this *ServiceAttachment) GetMany(ctx context.Context, request *protobuf.AttachmentGetMany) (*em_protobuf.Response, error) {
	// Validate
	{
		var vd validate_AttachmentGetMany
		err := em.Validator.Validate(request, &vd)
		if err != nil {
			em.LogWarn.OutputSimplePath(err)
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// If no ids, skip
	if len(request.GetOwnerIds()) == 0 {
		return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), nil)
	}

	var a []model.Attachment
	result := em.DB.Where("service = ?", request.GetService()).Where("owner_id IN ?", request.GetOwnerIds()).Where("owner_type = ?", request.GetOwnerType()).Find(&a)
	if result.RowsAffected == 0 {
		em.LogInfo.OutputSimplePath("No record")
	}

	for k, v := range a {
		// If it is stored locally, the path plus the domain name
		// 如果是本地储存，则路径加上域名
		v.MakeUrlPath(&a[k])
	}

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
		em.LogError.OutputSimplePath(err)
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Clear"), nil, err)
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Clear"), nil)
}

type validate_AttachmentAppend struct {
	Service string	`json:"service" validate:"required,max=255"`
	Paths []string	`json:"paths"`
	OwnerID uint	`json:"owner_id" validate:"required"`
	OwnerType string	`json:"owner_type" validate:"required,max=255"`
}
func (this *ServiceAttachment) Append(ctx context.Context, request *protobuf.AttachmentAppend) (*em_protobuf.Response, error) {
	{
		// Validate
		var vd validate_AttachmentAppend
		err := em.Validator.Validate(request, &vd)
		if err != nil {
			em.LogWarn.OutputSimplePath(err)
			return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Validate"), nil, err)
		}
	}

	// Set storage method
	fsm, err := em.Kv.ReadKey(define.MakeServiceConfField(em.Micro.Config.Service.RpcName, application.KvServiceFileStorageMethod))
	if err != nil {
		em.LogError.OutputSimplePath(err)
	}

	request.StorageMethod = fsm
	var attachment model.Attachment
	var tmp []string
	for _, v := range request.Paths {
		tmp = append(tmp, attachment.GetUrlPath(v))
	}
	request.Paths = tmp

	err = em.DB.Transaction(func(tx *gorm.DB) error {

		// 4. If the new attachment list is not empty, add a new attachment
		// 4.如果新附件列表不为空，则增加新的附件
		if len(request.Paths) > 0 {
			m, err := em.StructToMap(request)
			if err != nil {
				em.LogError.OutputSimplePath(err.Error())
				return err
			}
			// Note: json to map int format will be converted to float
			// 注意：json转map int格式会转换为float
			m["owner_id"] = uint(m["owner_id"].(float64))
			delete(m, "paths")

			result := tx.Model(model.Attachment{}).Where("path IN ?", request.Paths).Updates(m)
			if result.Error != nil {
				em.LogError.OutputSimplePath(result.Error)
				return result.Error
			}
		}

		return nil
	})
	if err != nil {
		return em.ErrorRpc(codes.InvalidArgument, em.ERROR_Code, em.I18n.TranslateFromRequest(ctx, "ERROR_Create"), nil, err)
	}

	return em.SuccessRpc(em.SUCCESS_Code, em.I18n.TranslateFromRequest(ctx, "SUCCESS_Create"), nil)
}
