package ninex

import (
	"encoding/json"
	"fmt"
	"gmkdb/global"
	"gmkdb/models/ninex"
	"gmkdb/models/ninex/request"
	"gmkdb/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// NexExpert

// @Author: [keid]
// @Function: GetExpertData
// @Description: 获取高手数据
// @Date: 2026-10-28
func GetExpertData(uid uint) (ninex.NexExpert, []ninex.NexRecord) {
	url := fmt.Sprintf("https://49129.com/forums/h5/prediction/prediction/detail?pageNum=1&pageSize=15&predictionTypeId=31&userId=%d&lotteryType=2", uid)

	client := &http.Client{Timeout: 30 * time.Second}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Referer", "https://49129.com/gaosou/details?typeId=31&uid=3211553&lian=4")

	resp, err := client.Do(req)
	if err != nil {
		// log.Fatal("请求失败:", err)
		global.GIA_LOG.Error("请求失败:", zap.Error(err))
	}
	defer resp.Body.Close()

	// fmt.Println("HTTP 状态码:", resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)

	// // 保存完整响应用到调试
	// os.WriteFile("api_response.json", body, 0644)
	// fmt.Println("已保存完整响应到 api_response.json")

	// 解析
	var base request.BaseResp
	if err := json.Unmarshal(body, &base); err != nil {
		global.GIA_LOG.Error("JSON 解析失败:", zap.Error(err))
	}

	var datas []ninex.NexRecord
	var user ninex.NexExpert

	if len(base.Data.List) != 0 && base.Success {
		user = ninex.NexExpert{
			UserId:         base.Data.UserId,
			NickName:       base.Data.Nickname,
			FromForumsName: base.Data.FromForumsName,
			HeadPic:        base.Data.HeadPic,
			FansCount:      base.Data.FansCount,
			ShareCount:     base.Data.ShareCount,
			ThumbUpCount:   base.Data.ThumbUpCount,
		}

		if uid != base.Data.UserId {
			global.GIA_LOG.Error("GetExpertDatas Uid ultimately inconsistent!", zap.Any(strconv.Itoa(int(uid)), base.Data.UserId))
			uid = base.Data.UserId
		}

		// 逐条写入
		for _, record := range base.Data.List {
			prefix := record.PredictionID[:7] // "2026298"
			num, err := strconv.Atoi(prefix)
			if num%1000 != record.Period || err != nil {
				global.GIA_LOG.Error("Period 自定义期数解析失败:", zap.String("prefix", prefix), zap.Int("num", num), zap.Int("Period", record.Period))
			}
			if num < 2026014 {
				continue
			}
			t, err := utils.ToTime(record.CreateTime)
			if err != nil {
				global.GIA_LOG.Error("时间解析失败:", zap.Error(err))
			}
			// fmt.Printf("CreateTime: %v\n", record.CreateTime)
			// fmt.Printf("t: %v, %v\n", t, t.Unix())

			lotteryNumber := record.LotteryNumber
			if len(lotteryNumber) == 0 {
				lotteryNumber = "末开" // 未来末日
			}

			for _, detail := range record.DetailList {
				var success bool
				if detail.Success == 1 {
					success = true // "命中"
				}
				datas = append(datas, ninex.NexRecord{
					Values:           strings.Join(detail.Value, ""),
					Period:           num,
					Success:          success,
					ContinuousCount:  record.ContinuousCount,
					LotteryNumber:    lotteryNumber,
					CreateTime:       t,
					LotteryType:      record.LotteryType,
					PredictionTypeID: record.PredictionTypeID,
					PredictionID:     record.PredictionID,
					UserId:           uid,
				})
			}
		}

	} else {
		global.GIA_LOG.Warn("list 为空，可能接口需要 Cookie 或参数不对", zap.Uint("uid", uid))
	}

	// // fmt.Printf("code=%d\n", base.Code)
	// fmt.Printf("昵称：%s\n来自论坛: %s\n", base.Data.Nickname, base.Data.FromForumsName)
	// fmt.Printf("粉丝数: %d\n转发数: %d\n点赞数: %d\n", base.Data.FansCount, base.Data.ShareCount, base.Data.ThumbUpCount)
	// fmt.Printf("本页记录数: %d | %d\n", len(base.Data.List), len(datas))
	// fmt.Printf("base.Data.List[0]: %v\n", base.Data.List[0])
	// fmt.Printf("datas[0]: %v\n", datas[0])
	// fmt.Printf("datas[%d]: %v\n", len(datas)-1, datas[len(datas)-1])

	// fmt.Printf("user: %v\n", user)

	// fmt.Printf("base.Msg: %v\n", base.Msg)
	// fmt.Printf("base.ServerTime: %v -> %v\n", base.ServerTime, utils.TransTime(base.ServerTime))
	// fmt.Printf("base.SubCode: %v\n", base.SubCode)
	// fmt.Printf("base.SubMsg: %v\n", base.SubMsg)
	// fmt.Printf("base.Success: %v\n", base.Success)

	return user, datas
}
