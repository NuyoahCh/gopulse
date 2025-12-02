package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// 餐厅数据结构
type Restaurant struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location string   `json:"location"`
	Cuisine  string   `json:"cuisine"`
	Rating   float64  `json:"rating"`
	Tags     []string `json:"tags"`
}

type Dish struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Desc  string  `json:"desc"`
	Spicy bool    `json:"spicy"`
}

// 模拟餐厅数据库
var restaurantDB = []Restaurant{
	{
		ID:       "r1",
		Name:     "川香阁",
		Location: "北京",
		Cuisine:  "川菜",
		Rating:   4.8,
		Tags:     []string{"辣", "正宗", "环境好"},
	},
	{
		ID:       "r2",
		Name:     "粤味轩",
		Location: "北京",
		Cuisine:  "粤菜",
		Rating:   4.6,
		Tags:     []string{"清淡", "海鲜", "精致"},
	},
	{
		ID:       "r3",
		Name:     "麻辣香锅",
		Location: "北京",
		Cuisine:  "川菜",
		Rating:   4.5,
		Tags:     []string{"辣", "实惠", "分量足"},
	},
}

var dishDB = map[string][]Dish{
	"r1": {
		{Name: "水煮鱼", Price: 88, Desc: "鲜嫩鱼肉，麻辣鲜香", Spicy: true},
		{Name: "宫保鸡丁", Price: 48, Desc: "经典川菜，香辣可口", Spicy: true},
		{Name: "蒜泥白肉", Price: 38, Desc: "肥而不腻，蒜香浓郁", Spicy: false},
	},
	"r2": {
		{Name: "清蒸鲈鱼", Price: 128, Desc: "鱼肉鲜嫩，原汁原味", Spicy: false},
		{Name: "白切鸡", Price: 68, Desc: "皮爽肉滑，鸡味浓郁", Spicy: false},
		{Name: "广式烧鹅", Price: 98, Desc: "皮脆肉嫩，香味四溢", Spicy: false},
	},
	"r3": {
		{Name: "麻辣香锅", Price: 58, Desc: "自选食材，麻辣鲜香", Spicy: true},
		{Name: "干锅牛蛙", Price: 68, Desc: "肉质细嫩，香辣入味", Spicy: true},
		{Name: "毛血旺", Price: 48, Desc: "麻辣鲜香，食材丰富", Spicy: true},
	},
}

func main() {
	ctx := context.Background()

	// 创建模型
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
		Model:   "deepseek-chat",
		BaseURL: "https://api.deepseek.com",
	})
	if err != nil {
		log.Fatalf("创建模型失败: %v", err)
	}

	// 工具1: 查询餐厅
	restaurantTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "query_restaurants",
			Desc: "根据条件查询餐厅信息（位置、菜系、是否辣）",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"location": {
					Type:     "string",
					Desc:     "城市位置，例如：北京、上海",
					Required: false,
				},
				"cuisine": {
					Type:     "string",
					Desc:     "菜系类型，例如：川菜、粤菜",
					Required: false,
				},
				"spicy": {
					Type:     "boolean",
					Desc:     "是否要辣的",
					Required: false,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			fmt.Printf("\\n[工具执行] query_restaurants\\n")

			var results []Restaurant
			location, _ := params["location"].(string)
			cuisine, _ := params["cuisine"].(string)
			spicy, _ := params["spicy"].(bool)

			for _, r := range restaurantDB {
				match := true
				if location != "" && r.Location != location {
					match = false
				}
				if cuisine != "" && r.Cuisine != cuisine {
					match = false
				}
				if spicy {
					hasSpicy := false
					for _, tag := range r.Tags {
						if tag == "辣" {
							hasSpicy = true
							break
						}
					}
					if !hasSpicy {
						match = false
					}
				}

				if match {
					results = append(results, r)
				}
			}

			resultJSON, _ := json.Marshal(results)
			return string(resultJSON), nil
		},
	)

	// 工具2: 查询菜品
	dishTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "query_dishes",
			Desc: "查询指定餐厅的菜品信息（需要餐厅ID）",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"restaurant_id": {
					Type:     "string",
					Desc:     "餐厅ID",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params map[string]any) (string, error) {
			restaurantID := params["restaurant_id"].(string)
			fmt.Printf("\\n[工具执行] query_dishes(restaurant_id=%s)\\n", restaurantID)

			dishes, exists := dishDB[restaurantID]
			if !exists {
				return `{"error": "餐厅不存在"}`, nil
			}

			resultJSON, _ := json.Marshal(dishes)
			return string(resultJSON), nil
		},
	)

	// 创建 React Agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{restaurantTool, dishTool},
		},
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	// 使用 Agent
	messages := []*schema.Message{
		schema.UserMessage("我在北京，想吃辣一点的，给我推荐几家餐厅和特色菜"),
	}

	fmt.Println("=== 用户: 我在北京，想吃辣一点的，给我推荐几家餐厅和特色菜 ===")

	response, err := agent.Generate(ctx, messages)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("\\n=== Agent 回答 ===\\n%s\\n", response.Content)
}
