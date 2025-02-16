# DeepLLM Tourism Assistant

基于多智能体的杭州旅游助手系统，提供智能行程规划、餐饮住宿推荐、天气建议等功能。

## 项目结构

```
├── components/
│   └── agent/
│       ├── accommodation/ # 住宿推荐智能体
│       ├── dining/        # 餐饮推荐智能体
│       ├── weather/       # 天气建议智能体
│       ├── planner/       # 行程规划智能体
│       └── coordinator/   # 多智能体协调器
├── internal/
│   └── data/
│       ├── loader.go      # 数据加载器
│       ├── models.go      # 数据模型定义
│       ├── query.go       # 数据查询接口
│       └── distance.go    # 距离计算工具
├── data/
│   ├── # 区域数据 景点数据 餐厅数据 酒店数据 天气数据
└── cmd/
    ├── guide/            # 基础使用示例
    └── multiagent/       # 多智能体示例
```

## 功能特点

1. 智能行程推荐
   - 基于用户偏好的景点推荐
   - 考虑天气因素的行程调整
   - 合理的时间安排

2. 餐饮住宿建议
   - 根据预算和偏好推荐餐厅
   - 智能住宿匹配
   - 特色美食推荐

3. 路线规划优化
   - 基于距离的智能排序
   - 考虑交通便利性
   - 景点间最优路线

4. 天气相关建议
   - 实时天气预报
   - 室内外活动建议
   - 出行装备提醒

5. 个性化定制
   - 多样化的筛选条件
   - 灵活的预算控制
   - 特殊需求适配

## 快速开始

1. 安装依赖
```bash
go mod tidy
```

2. 运行示例
```bash
# 运行基础示例
go run cmd/guide/main.go

# 运行多智能体示例
go run cmd/multiagent/main.go
```

3. 示例输出
```
=== 行程概览 ===
行程天数: 3天
亮点推荐:
- 2025-02-18适合：徒步西湖、灵隐寺参观

=== 住宿安排 ===
酒店: 杭州柏悦酒店 (五星级)
地址: 西湖区浦春路28号

=== 每日行程 ===
2025-02-18:
天气情况:
- 适合活动: 户外观光、徒步游览
...
```

## 开发指南

1. 添加新的智能体
   - 在 components/agent 下创建新的包
   - 实现对应的接口
   - 在 coordinator 中集成

2. 扩展数据模型
   - 在 internal/data/models.go 中添加新的结构体
   - 在 data/ 目录下添加对应的 JSON 数据文件
   - 在 internal/data/loader.go 中添加加载逻辑

3. 自定义推荐算法
   - 修改各个智能体中的评分和推荐逻辑
   - 在 coordinator 中调整协调策略

## 数据格式

所有数据文件采用 JSON 格式，具体结构请参考 data/ 目录下的示例文件。

## 注意事项

1. 所有坐标采用 WGS84 坐标系
2. 价格单位统一使用 CNY（人民币）
3. 时间格式遵循 RFC3339 标准
4. 距离单位统一使用公里（km）
