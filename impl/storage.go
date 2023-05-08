package impl

import (
	"fmt"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func init() {
	rpc.Post(types.CmdPathResourcePut, PutResource, types.RoleTypeStudent, types.RoleTypeTeacher, types.RoleTypeAdmin)

	rpc.Use(types.CmdPathResourceGet, GetResource, types.RoleTypePublic)
}

func PutResource(ctx *goon.Ctx, req *types.PutResourceReq) (*types.PutResourceRsp, error) {
	var rsp types.PutResourceRsp

	var allowDup bool
	if req.Path == "" {
		req.Path = fmt.Sprintf("%s/%s/%s",
			ctx.GetReqHeader(types.HeaderRoleType),
			ctx.GetReqHeader(types.HeaderUserId),
			utils.Sha512(string(req.Body)))

		allowDup = true
	}

	if req.ContentType == "" {
		req.ContentType = http.DetectContentType(req.Body[:512])
	}

	if req.Path[0] == '/' {
		req.Path = req.Path[1:]
	}

	if req.Path[len(req.Path)-1] == '/' {
		req.Path = req.Path[:len(req.Path)-1]
	}

	var id uint64
	err := db.Model(&types.ModelResource{}).
		Where("path = ?", req.Path).
		Select("Id").Scan(&id).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if id > 0 {
		if allowDup {
			rsp.Path = types.CmdPathResourceGet + req.Path
			return &rsp, nil
		}

		return nil, types.CreateError(types.ErrPathAlreadyExists)
	}

	resource := &types.ModelResource{
		Path: req.Path,
		ResourceDetail: &types.ResourceDetail{
			Body: req.Body,

			ContentType: req.ContentType,
		},
	}

	err = db.Create(resource).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Path = types.CmdPathResourceGet + req.Path
	return &rsp, nil
}

func GetResource(ctx *goon.Ctx) error {
	switch ctx.Method() {
	case goon.MethodGet, goon.MethodHead:
		path := ctx.Path()
		path = path[len(types.CmdPathResourceGet):]
		if path[0] == '/' {
			path = path[1:]
		}
		if path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}

		var resource types.ModelResource
		err := db.Where("path = ?", path).First(&resource).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.SetStatusCode(http.StatusNotFound)
				return nil
			}
			log.Errorf("err:%v", err)
			return err
		}

		ctx.SetResHeader("Content-Type", resource.ResourceDetail.ContentType)
		ctx.SetResHeader("Last-Modified", time.Unix(int64(resource.UpdatedAt), 0).UTC().Format(http.TimeFormat))
		ctx.SetResHeader("Cache-Control", "public, max-age=86400")
		ctx.SetResHeader("Content-Disposition", "inline")

		args := ctx.Context().Request.URI().QueryArgs()
		if args != nil {
			filename := args.Peek("filename")
			if len(filename) > 0 {
				ctx.SetResHeader("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
			}
		}

		if ctx.Method() == goon.MethodGet {
			return ctx.Send(string(resource.ResourceDetail.Body))
		}

		return nil
	case goon.MethodOption:
		return nil
	default:
		return ctx.Next()
	}
}
