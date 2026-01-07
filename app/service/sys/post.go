/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 岗位
 */

package sys

import (
	"context"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/domain/vo"
	base2 "github.com/sufo/bailu-admin/app/service/base"
)

var PostSet = wire.NewSet(wire.Struct(new(PostOption), "*"), NewPostService)

type PostOption struct {
	PostRepo  *repo.PostRepo
	UserRepo  *repo.UserRepo
	TransRepo *repo.Trans
}
type PostService struct {
	base2.BaseService[entity.Post]
	PostOption
}

func NewPostService(opt PostOption) *PostService {
	return &PostService{base2.BaseService[entity.Post]{opt.PostRepo.Repository}, opt}
}

func (p *PostService) Create(ctx context.Context, post *entity.Post) error {
	return p.PostRepo.Create(ctx, post)
}

// 查询所有符合条件的，无dataScope限制
func (p *PostService) List(ctx context.Context, post dto.PostParams) (*resp.PageResult[entity.Post], error) {
	builder := base.NewQueryBuilder()
	builder.WithWhereStruct(post).
		WithPagination(ctx)
	if result, err := p.PostRepo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		return result.(*resp.PageResult[entity.Post]), err
	}
}

func (p *PostService) Update(ctx context.Context, post *entity.Post) (*entity.Post, error) {
	err := p.PostRepo.Update(ctx, post)
	return post, err
}

func (p *PostService) Delete(ctx context.Context, ids []uint64) error {
	return p.TransRepo.Exec(ctx, func(ctx context.Context) error {
		//解除用户和岗位关系
		err := p.UserRepo.UntiedPost(ctx, ids)
		if err != nil {
			return err
		}
		return p.PostRepo.Delete(ctx, ids)
	})
}

func (p *PostService) FindOptions(ctx context.Context) ([]*vo.Option[uint64], error) {
	builder := base.NewQueryBuilder()
	builder.WithSelect("id value").WithSelect("name label").
		WithWhere("status=?", 1). //启用
		WithOrder("sort asc")
	options := make([]*vo.Option[uint64], 0)
	err := p.PostRepo.FindModelByBuilder(ctx, builder, &options)
	return options, err
}
