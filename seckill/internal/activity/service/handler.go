package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"Redrock/seckill/internal/activity/data"
	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/models"
	"Redrock/seckill/internal/pkg/redis"
	activity "Redrock/seckill/kitex_gen/activity"
)

// InternalActivityServiceImpl implements the last service interface defined in the IDL.
type ActivityServiceImpl struct{
	activityData 	*data.ActivityData
	activityRedis 	*data.ActivityRedis
}

// NewInternalActivityServiceImpl 创建服务实例
func NewInternalActivityServiceImpl() *ActivityServiceImpl{
	return &ActivityServiceImpl{
		activityData	: data.NewActivityData(),
		activityRedis	: data.NewActivityRedis(),
	}
}

// NewActivityServiceImpl 创建活动服务实例
func NewActivityServiceImpl() *ActivityServiceImpl {
	return NewInternalActivityServiceImpl()
}

// CreateActivity 创建活动
func (s *ActivityServiceImpl) CreateActivity(ctx context.Context, req *activity.CreateActivityRequest) (*activity.CreateActivityResponse, error){
	response := &activity.CreateActivityResponse{
		BaseResponse : &activity.BaseResponse{},
	}

	// 参数验证
	if req.Name == "" || req.ProductID <= 0 || req.SeckillPrice <0 || req.TotalStock <= 0{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg = "参数错误"

		return response, nil
	}

	// 检查视奸是否合法
	startTime := time.Unix(req.StartTime, 0)
	endTime	:= time.Unix(req.EndTime, 0)
	// now := time.Now()

	// if startTime.Before(now){
	// 	response.BaseResponse.Code = 400
	// 	response.BaseResponse.Msg = "活动开始时间不能早于当前时间"

	// 	return response, nil
	// }

	if endTime.Before(startTime){
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg = "活动结束时间不能早于开始时间"

		return response, nil
	}

	var product models.Product
	result := database.GetDB().First(&product, req.ProductID)
	if result.Error != nil{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg = "商品不存在"

		return response, nil
	}

	activity := &models.Activity{
		Name:			req.Name,
		ProductID:		uint(req.ProductID),
		StartTime:		startTime,
		EndTime:		endTime,
		SeckillPrice:	req.SeckillPrice,
		TotalStock:		req.TotalStock,
		AvailableStock:	req.TotalStock,
		Status:			0, // 未开始
	}

	// 考虑到秒杀系统的高并发，我们选择先将商品信息存入缓存
	database.GetDB().Model(activity).Association("Product").Find(&activity.Product)

	// 将数据写到数据库
	err := s.activityData.Create(ctx, activity)
	if err != nil{
		response.BaseResponse.Code = 500
		response.BaseResponse.Msg = "创建活动失败" + err.Error()
		return response, nil
	}

	// 当活动创建成功后，将活动信息写入Redis
	err = s.activityRedis.SaveActivity(ctx, activity)
	if err != nil{
		// 目的是创建活动，缓存是非关键操作，因此只记录不退出
		log.Printf("缓存活动信息失败：%v\n", err)
	}

	// 将库存写入Redis(预热)
	quantity,err := s.activityRedis.InitStock(ctx, activity.ID,activity.AvailableStock)
	log.Printf("库存数量：%d\n", quantity)
	if err != nil{
		log.Printf("初始化Redis库存失败：%v\n", err)
	}

	response.BaseResponse.Code = 0
	response.BaseResponse.Msg = "活动创建成功"
	response.ActivityID = int64(activity.ID)

	return response, nil
}

// GetActivityList 获取活动列表
func (s *ActivityServiceImpl) GetActivityList(ctx context.Context, req *activity.GetActivityListRequest) (*activity.GetActivityListResponse, error){
	response := &activity.GetActivityListResponse{
		BaseResponse: &activity.BaseResponse{},
		Activities:	[]*activity.ActivityInfo{},
	}

	// 更新所有活动状态
	err := s.activityData.AutoUpdateActivityStatus(ctx)
	if err != nil{
		log.Printf("自动更新活动状态失败：%v\n", err)
	}

	// 查询活动列表
	activities, total, err := s.activityData.List(ctx, int(req.Status))
	if err != nil{
		response.BaseResponse.Code = 500
		response.BaseResponse.Msg  = "查询活动列表失败：" + err.Error()

		return response, nil
	}

	// 构建response
	for _, a := range activities{
		activityInfo := &activity.ActivityInfo{
			Id:					int64(a.ID),
			Name:				a.Name,
			ProductId:			int64(a.ProductID),
			SeckillPrice:		a.SeckillPrice,
			OriginalPrice:		0,
			StartTime:			a.StartTime.Unix(),
			EndTime:			a.EndTime.Unix(),
			TotalStock:			a.TotalStock,
			AvailableStock:		a.AvailableStock,
			Status:				int32(a.Status),
		}
		
		if a.Product.ID != 0{
			activityInfo.ProductName = a.Product.Name
			activityInfo.OriginalPrice = a.Product.Price
		}

		response.Activities = append(response.Activities, activityInfo)
	}

	response.Total = total
	response.BaseResponse.Code = 0
	response.BaseResponse.Msg = "查询活动列表成功"

	return response, nil
}

// GetActivity 获取活动详情
func (s *ActivityServiceImpl) GetActivity(ctx context.Context, req *activity.GetActivityRequest) (*activity.GetActivityResponse, error){
	response := &activity.GetActivityResponse{
		BaseResponse: &activity.BaseResponse{},
	}

	// 更新所有活动状态
	err := s.activityData.AutoUpdateActivityStatus(ctx)
	if err != nil{
		log.Printf("自动更新活动状态失败：%v\n", err)
	}

	if req.ActivityID <= 0{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "活动ID不能为空"
		
		return response, nil
	}

	// 首先考虑从Redis中获取活动信息
	localActivity, err := s.activityRedis.GetActivity(ctx, uint(req.ActivityID))
	if err != nil{
		log.Printf("从Redis中获取活动信息失败：%v", err)
	}

	// 如果无法从Redis中获取信息，则从数据库中获取
	if localActivity == nil{
		localActivity, err = s.activityData.GetByID(ctx, uint(req.ActivityID))

		if err != nil{
			response.BaseResponse.Code = 404
			response.BaseResponse.Msg  = "获取活动信息失败" + err.Error()

			return response, err
		}
	}

	stock, err := s.activityRedis.GetStock(ctx, uint(req.ActivityID))
	if err == nil && stock >= 0{
		// 当从Redis能获取正确的库存时替换，否则使用数据库的库存
		localActivity.AvailableStock = stock
	}

	activityInfo := &activity.ActivityInfo{
		Id:						int64(localActivity.ID),
		Name:					localActivity.Name,
		ProductId:				int64(localActivity.ProductID),
		ProductName:			localActivity.Product.Name,
		OriginalPrice:			localActivity.Product.Price,
		SeckillPrice:			localActivity.SeckillPrice,
		StartTime:				localActivity.StartTime.Unix(),
		EndTime:				localActivity.EndTime.Unix(),
		TotalStock:				localActivity.TotalStock,
		AvailableStock:			localActivity.AvailableStock,
		Status: 				int32(localActivity.Status),
	}

	response.Activity = activityInfo
	response.BaseResponse.Code = 0
	response.BaseResponse.Msg = "查询活动信息成功"
	
	return response, nil
}

// DeductStoc 扣除库存
func (s *ActivityServiceImpl) DeductStock(ctx context.Context, req *activity.DeductStockRequest) (*activity.DeductStockResponse, error){
	response := &activity.DeductStockResponse{
		BaseResponse: &activity.BaseResponse{},
		Success: 		false,	
	}

	// 创建分布式锁
	lockKey := fmt.Sprintf("activity:lock:%d", req.ActivityID)
	lock 	:= redis.NewDistributedLock(s.activityRedis.GetRedis(), lockKey, 1*time.Second)

	try, err := lock.TryLock(ctx)
	// 返回err
	if err != nil{
		log.Printf("获取分布式锁失败：%v", err)
	// 没有成功返回0
	}else if !try{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "系统繁忙，请稍后再试"

		return response, nil
	// 成功defer unlock确保释放锁
	} else{
		defer lock.Unlock(ctx)
	}


	if req.ActivityID <= 0 || req.UserID <= 0{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "活动或用户参数错误"

		return response, nil
	}

	// 检查用户是否参与过活动
	joined, err := s.activityRedis.IsUserJoined(ctx, uint(req.UserID), uint(req.ActivityID))
	if err != nil{
		log.Printf("检查用户是否参与过秒杀活动失败：%v", err)
	}

	if joined{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "您已参与过此秒杀活动"

		return response, nil
	}

	// 获取活动信息以便之后检查活动状态
	localActivity, err := s.activityData.GetByID(ctx, uint(req.ActivityID))
	if err != nil{
		response.BaseResponse.Code = 404
		response.BaseResponse.Msg  = "获取活动信息失败" + err.Error()

		return response, nil
	}

	// 检查活动状态
	now := time.Now()

	// 检查活动是否开始
	if now.Before(localActivity.StartTime){
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "活动尚未开始"

		return response, nil
	}

	// 检查活动是否结束
	if now.After(localActivity.EndTime){
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "活动已结束"

		return response, nil
	}

	// 其他原因关闭活动
	if !localActivity.IsAvailable(){
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "该活动暂不可用"

		return response, nil
	}

	// 从Redis扣除库存
	success, err := s.activityRedis.DeductStock(ctx, uint(req.ActivityID), req.Count)
	if err != nil{
		response.BaseResponse.Code = 500
		response.BaseResponse.Msg  = "扣除库存失败" + err.Error()

		return response, nil
	}

	if !success{
		response.BaseResponse.Code = 400
		response.BaseResponse.Msg  = "库存不足"

		return response, nil
	}

	// 记录用户参与秒杀活动
	err = s.activityRedis.RecordUserJoin(ctx, uint(req.UserID), uint(req.ActivityID))
	if err != nil{
		log.Printf("记录用户参与秒杀活动失败：%v\n", err)
	}

	// 开启一个协程用于异步更新数据库中的库存
	go func(){
		newCtx := context.Background()

		// 获取Redis中的库存
		currentStock, err := s.activityRedis.GetStock(newCtx, uint(req.ActivityID))
		if err != nil{
			log.Printf("获取Redis中的库存失败：%v\n", err)
			return
		}

		// 并更新到数据库
		err = s.activityData.UpdateStock(newCtx, uint(req.ActivityID),currentStock)
		if err != nil{
			log.Printf("更新数据库库存失败：%v\n", err)
		}
	}()

	response.BaseResponse.Code = 0
	response.BaseResponse.Msg  = "扣除库存成功"
	response.Success = true

	return response, nil
}
