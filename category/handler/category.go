package handler

import (
	"context"
	common "github.com/HuangQinTang/micro_shop_common"
	"github.com/micro_shop/category/domain/model"
	"github.com/micro_shop/category/domain/service"
	category "github.com/micro_shop/category/proto/category"
	"github.com/prometheus/common/log"
)

type Category struct {
	CategoryDataService service.ICategoryDataService
}

//提供创建分类的服务
func (c *Category) CreateCategory(ctx context.Context, Req *category.CategoryReq, Res *category.CreateCategoryRes) error {
	category := &model.Category{}
	//赋值
	err := common.SwapTo(Req, category)
	if err != nil {
		return err
	}
	categoryId, err := c.CategoryDataService.AddCategory(category)
	if err != nil {
		return err
	}
	Res.Message = "分类添加成功"
	Res.CategoryId = categoryId
	Res.TraceId = common.WithTrace(ctx)
	return nil
}

//提供分类更新服务
func (c *Category) UpdateCategory(ctx context.Context, Req *category.CategoryReq, Res *category.UpdateCategoryRes) error {
	category := &model.Category{}
	err := common.SwapTo(Res, category)
	if err != nil {
		return err
	}
	err = c.CategoryDataService.UpdateCategory(category)
	if err != nil {
		return err
	}
	Res.Message = "分类更新成功"
	Res.TraceId = common.WithTrace(ctx)
	return nil
}

//提供分类删除服务
func (c *Category) DeleteCategory(ctx context.Context, Req *category.DeleteCategoryReq, Res *category.DeleteCategoryRes) error {
	err := c.CategoryDataService.DeleteCategory(Req.CategoryId)
	if err != nil {
		return nil
	}
	Res.Message = "删除成功"
	Res.TraceId = common.WithTrace(ctx)
	return nil
}

//根据分类名称查找分类
func (c *Category) FindCategoryByName(ctx context.Context, Req *category.FindByNameReq, Res *category.CategoryRes) error {
	category, err := c.CategoryDataService.FindCategoryByName(Req.CategoryName)
	if err != nil {
		return err
	}
	Res.TraceId = common.WithTrace(ctx)
	return common.SwapTo(category, Res)

}

//根据分类ID查找分类
func (c *Category) FindCategoryByID(ctx context.Context, Req *category.FindByIdReq, Res *category.CategoryRes) error {
	category, err := c.CategoryDataService.FindCategoryByID(Req.CategoryId)
	if err != nil {
		return err
	}
	Res.TraceId = common.WithTrace(ctx)
	return common.SwapTo(category, Res)
}

func (c *Category) FindCategoryByLevel(ctx context.Context, Req *category.FindByLevelReq, Res *category.FindAllRes) error {
	categorySlice, err := c.CategoryDataService.FindCategoryByLevel(Req.Level)
	if err != nil {
		return err
	}
	categoryToRes(categorySlice, Res)
	Res.TraceId = common.WithTrace(ctx)
	return nil
}

func (c *Category) FindCategoryByParent(ctx context.Context, Req *category.FindByParentReq, Res *category.FindAllRes) error {
	categorySlice, err := c.CategoryDataService.FindCategoryByParent(Req.ParentId)
	if err != nil {
		return err
	}
	categoryToRes(categorySlice, Res)
	Res.TraceId = common.WithTrace(ctx)
	return nil
}

func (c *Category) FindAllCategory(ctx context.Context, Req *category.FindAllReq, Res *category.FindAllRes) error {
	categorySlice, err := c.CategoryDataService.FindAllCategory()
	if err != nil {
		return err
	}
	categoryToRes(categorySlice, Res)
	Res.TraceId = common.WithTrace(ctx)
	return nil
}

func categoryToRes(categorySlice []model.Category, Res *category.FindAllRes) {
	for _, cg := range categorySlice {
		cr := &category.CategoryDesc{}
		err := common.SwapTo(cg, cr)
		if err != nil {
			log.Error(err)
			break
		}
		Res.Category = append(Res.Category, cr)
	}
}
