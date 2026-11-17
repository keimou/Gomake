package request

// detailList 里的每一项
type DetailItem struct {
	Success      int      `json:"success"`      // 1=命中 0=未命中
	SuccessValue []string `json:"successValue"` // 命中的生肖: "虎"
	Value        []string `json:"value"`        // 预测的生肖数组: ["鼠","马","鸡","龙","蛇","兔","虎","牛","猪"]
}

// list 里的每一项
type PredictionRecord struct {
	ContinuousCount  string       `json:"continuousCount"` // "1连中"
	CreateTime       string       `json:"createTime"`      // "2026年01月12日 15:50:58"
	DetailList       []DetailItem `json:"detailList"`
	Period           int          `json:"period"`           // 298
	PredictionID     string       `json:"predictionId"`     // 20262982313211553 => 2026 298 2 31 3211553 （date period lotteryType  predictionType Uid）
	LotteryType      int          `json:"lotteryType"`      // 2
	PredictionTypeID int          `json:"predictionTypeId"` // 31
	LotteryNumber    string       `json:"lotteryNumber"`    // "08猪"
}

// data 层
type DataStruct struct {
	Concern        int                `json:"concern"`        // 0
	ConcernCount   int                `json:"concernCount"`   // 0
	FansCount      int                `json:"fansCount"`      // 969
	FromForumsName string             `json:"fromForumsName"` // "118论坛"
	HeadPic        string             `json:"headPic"`        // "https://img.image049.com:4949/forumsfiles/53.png"
	List           []PredictionRecord `json:"list"`
	// ---
	LotteryType  int    `json:"lotteryType"`  // 2
	Nickname     string `json:"nickname"`     // "忧心的钢笔阹"
	ShareCount   int    `json:"shareCount"`   // 226
	ThumbUp      int    `json:"thumbUp"`      // 0
	ThumbUpCount int    `json:"thumbUpCount"` // 246
	UserId       uint   `json:"userId"`       // 3211553
}

// 接口外层
type BaseResp struct {
	Code       int        `json:"code"` // 10000
	Data       DataStruct `json:"data"`
	Msg        string     `json:"msg"`
	ServerTime int64      `json:"serverTime"`
	SubCode    int        `json:"subCode"`
	SubMsg     string     `json:"subMsg"`
	Success    bool       `json:"success"`
}

// {
// 	"code":10000,
// 	"data":{
// 		"concern":0,
// 		"concernCount":0,
// 		"fansCount":969,
// 		"fromForumsName":"118论坛",
// 		"headPic":"https://img.image049.com:4949/forumsfiles/53.png",
// 		"list":[
// 			{
// 				"continuousCount":"1连中",
// 				"createTime":"2026年01月12日 15:50:58",
// 				"detailList":[
// 					{
// 						"success":1,
// 						"successValue":[
// 							"虎"
// 						],
// 						"value":[
// 							"鼠",
// 							"马",
// 							"鸡",
// 							"龙",
// 							"蛇",
// 							"兔",
// 							"虎",
// 							"牛",
// 							"猪"
// 						]
// 					}
// 				],
// 				"lotteryNumber":"16虎",
// 				"lotteryType":2,
// 				"paid":0,
// 				"pay":0,
// 				"payCount":0,
// 				"period":310,
// 				"predictionId":"20253102313211495",
// 				"predictionTypeId":31,
// 				"predictionTypeName":"九肖中特",
// 				"price":0.00,
// 				"rate":3.00,
// 				"rateStr":"300%",
// 				"success":1
// 			},
// 			....
// 			{
// 				"continuousCount":"-",
// 				"createTime":"2026年10月25日 21:48:07",
// 				"detailList":[
// 					{
// 						"success":0,
// 						"value":[
// 							"鸡",
// 							"鼠",
// 							"猪",
// 							"狗",
// 							"兔",
// 							"龙",
// 							"蛇",
// 							"羊",
// 							"猴"
// 						]
// 					}
// 				],
// 				"lotteryType":2,
// 				"paid":0,
// 				"pay":0,
// 				"payCount":0,
// 				"period":299,
// 				"predictionId":"20262992313211495",
// 				"predictionTypeId":31,
// 				"predictionTypeName":"九肖中特",
// 				"price":0.00,
// 				"rate":0.00,
// 				"rateStr":"0%"
// 			}
// 		],
// 		"lotteryType":2,
// 		"nickname":"涂钦百合搚",
// 		"shareCount":26,
// 		"thumbUp":0,
// 		"thumbUpCount":42,
// 		"userId":3211495
// 	},
// 	"msg":"操作成功",
// 	"serverTime":1792984442960,
// 	"subCode":10000,
// 	"subMsg":"操作成功",
// 	"success":true
// }
