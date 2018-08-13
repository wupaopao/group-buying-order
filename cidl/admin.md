### /group_buying/admin/task/list

```json
{
    "Code": 0,
    "Data": {
        "Count": 300,
        "List": [
            {
                "TaskId": 1,
                "Title": "团购任务-1",
                "Introduction": "非常好吃的糖",
                "CoverPicture": "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1515583303718&di=ebb526df71007c6bace7945d627e2b92&imgtype=0&src=http%3A%2F%2Fwww.szthks.com%2Flocalimg%2F687474703a2f2f6777312e616c6963646e2e636f6d2f62616f2f75706c6f616465642f69372f5431796177574664706458585858585858585f2121302d6974656d5f7069632e6a7067.jpg",
                "IllustrationPictures": [
                    "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1515583303717&di=889bf336a115d98f0513df2aeec69f09&imgtype=0&src=http%3A%2F%2Fimg008.hc360.cn%2Fhb%2FMTQ2MDkyMTY5NjkwNTg4NDAwNzM2OQ%3D%3D.jpg",
                    "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1515583303717&di=4679d5b884e016202465e51e28f6d4d7&imgtype=0&src=http%3A%2F%2Fimgsrc.baidu.com%2Fimage%2Fc0%253Dpixel_huitu%252C0%252C0%252C294%252C40%2Fsign%3D558fc1c972310a55d029d6b4de3d26c5%2F060828381f30e9241a20c70b47086e061d95f700.jpg"
                ],
                "Info": [
                    {
                        "Title": "品牌",
                        "Content": "世果汇"
                    },
                    {
                        "Title": "品名",
                        "Content": "墨西哥森林牛肉果"
                    },
                    {
                        "Title": "规格",
                        "Content": "6个"
                    }
                ],
                "Specification": {
                    "ItemIdCounter": 2,
                    "Items": [
                        {
                            "ItemId": "1",
                            "Name": "颜色",
                            "LabelIdCounter": 2,
                            "Labels": [
                                {
                                    "LabelId": "1:1",
                                    "Name": "黑色"
                                },
                                {
                                    "LabelId": "1:2",
                                    "Name": "白色"
                                }
                            ]
                        },
                        {
                            "ItemId": "2",
                            "Name": "尺寸",
                            "LabelIdCounter": 2,
                            "Labels": [
                                {
                                    "LabelId": "2:1",
                                    "Name": "L"
                                },
                                {
                                    "LabelId": "2:2",
                                    "Name": "XL"
                                }
                            ]
                        }
                    ],
                    "SkuMap": {
                        "1:1-2:1": {
                            "ItemId": "1:1-2:1",
                            "MarketPrice": 100,
                            "GroupBuyingPrice": 90,
                            "SettlementPrice": 80,
                            "CostPrice": 70,
                            "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg"
                        },
                        "1:1-2:2": {
                            "ItemId": "1:1-2:2",
                            "MarketPrice": 100,
                            "GroupBuyingPrice": 90,
                            "SettlementPrice": 80,
                            "CostPrice": 70,
                            "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg"
                        },
                        "1:2-2:1": {
                            "ItemId": "1:2-2:1",
                            "MarketPrice": 100,
                            "GroupBuyingPrice": 90,
                            "SettlementPrice": 80,
                            "CostPrice": 70,
                            "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg"
                        },
                        "1:2-2:2": {
                            "ItemId": "1:2-2:2",
                            "MarketPrice": 100,
                            "GroupBuyingPrice": 90,
                            "SettlementPrice": 80,
                            "CostPrice": 70,
                            "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg"
                        }
                    },
                  	"MarketPriceRange": { // 市场价价格区间
                        "Min": 70,
                        "Max": 200
                    },
                    "GroupBuyingPriceRange": { // 团购价价格区间
                        "Min": 71,
                        "Max": 201
                    },
                    "SettlementPriceRange": { // 结算价价格区间
                        "Min": 72,
                        "Max": 202
                    },
                    "CostPriceRange": { // 成本价价格区间
                        "Min": 73,
                        "Max": 203
                    }
                },
                "Detail": "古法酿造黑糖",
                "Notes": "周二21:00截单，下周一配送",
                "StartTime": "2018-01-10T17:39:59.3853796+08:00",
                "EndTime": "2018-01-11T17:39:59.3853796+08:00",
                "IsAutoPost": true,
                "GroupState": 0,
                "OrderState": 0,
                "Version": 1,
                "Sales": 1000
            }
        ]
    },
    "Message": "success"
}
```

详细字段意义请参考`db_model.cidl`



###/api/v1/group_buying/admin/indent/summary/:indent_id

```json
{
    "Code": 0,
    "Data": {
        "TaskStatistics": [
            {
                "IndentId": "20180112", // 订货单ID
                "TaskId": 1, // 团购任务ID
                "TaskContent": {
                    "TaskId": 1,
                    "Title": "团购任务-1",//商品标题
                    "Introduction": "非常好吃的糖",
                    "CoverPicture": "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1515583303718&di=ebb526df71007c6bace7945d627e2b92&imgtype=0&src=http%3A%2F%2Fwww.szthks.com%2Flocalimg%2F687474703a2f2f6777312e616c6963646e2e636f6d2f62616f2f75706c6f616465642f69372f5431796177574664706458585858585858585f2121302d6974656d5f7069632e6a7067.jpg",
                    "IllustrationPictures": [
                        "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1515583303717&di=889bf336a115d98f0513df2aeec69f09&imgtype=0&src=http%3A%2F%2Fimg008.hc360.cn%2Fhb%2FMTQ2MDkyMTY5NjkwNTg4NDAwNzM2OQ%3D%3D.jpg",
                        "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1515583303717&di=4679d5b884e016202465e51e28f6d4d7&imgtype=0&src=http%3A%2F%2Fimgsrc.baidu.com%2Fimage%2Fc0%253Dpixel_huitu%252C0%252C0%252C294%252C40%2Fsign%3D558fc1c972310a55d029d6b4de3d26c5%2F060828381f30e9241a20c70b47086e061d95f700.jpg"
                    ],
                    "Info": [ // 产品信息
                        {
                            "Title": "品牌",
                            "Content": "世果汇"
                        },
                        {
                            "Title": "品名",
                            "Content": "墨西哥森林牛肉果"
                        },
                        {
                            "Title": "规格",
                            "Content": "6个"
                        }
                    ],
                    "Specification": {// 产品规格
                        "ItemIdCounter": 2,
                        "Items": {
                            "1": { // 规格一
                                "ItemId": "1",
                                "Name": "颜色",
                                "LabelIdCounter": 2,
                                "Labels": { // 标签一
                                    "1:1": {
                                        "LabelId": "1:1",
                                        "Name": "黑色"
                                    },
                                    "1:2": { // 标签二
                                        "LabelId": "1:2",
                                        "Name": "白色"
                                    }
                                }
                            },
                            "2": {
                                "ItemId": "2",
                                "Name": "尺寸",
                                "LabelIdCounter": 2,
                                "Labels": {
                                    "2:1": {
                                        "LabelId": "2:1",
                                        "Name": "L"
                                    },
                                    "2:2": {
                                        "LabelId": "2:2",
                                        "Name": "XL"
                                    }
                                }
                            }
                        },
                        "SkuMap": { // 价格映射表
                            "1:1-2:1": { // 标签ID以“-”拼接起来的ID
                                "ItemId": "1:1-2:1",
                                "Labels": {
                                    "1:1": {
                                        "LabelId": "1:1",
                                        "Name": "黑色"
                                    },
                                    "2:1": {
                                        "LabelId": "2:1",
                                        "Name": "L"
                                    }
                                },
                                "MarketPrice": 100,
                                "GroupBuyingPrice": 90,
                                "SettlementPrice": 80,
                                "CostPrice": 70,
                                "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg"
                            },
                            "1:1-2:2": {
                                "ItemId": "1:1-2:2",
                                "Labels": {
                                    "1:1": {
                                        "LabelId": "1:1",
                                        "Name": "黑色"
                                    },
                                    "2:2": {
                                        "LabelId": "2:2",
                                        "Name": "XL"
                                    }
                                },
                                "MarketPrice": 100,
                                "GroupBuyingPrice": 90,
                                "SettlementPrice": 80,
                                "CostPrice": 70,
                                "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg"
                            },
                            "1:2-2:1": {
                                "ItemId": "1:2-2:1",
                                "Labels": {
                                    "1:2": {
                                        "LabelId": "1:2",
                                        "Name": "白色"
                                    },
                                    "2:1": {
                                        "LabelId": "2:1",
                                        "Name": "L"
                                    }
                                },
                                "MarketPrice": 100,
                                "GroupBuyingPrice": 90,
                                "SettlementPrice": 80,
                                "CostPrice": 70,
                                "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg"
                            },
                            "1:2-2:2": {
                                "ItemId": "1:2-2:2",
                                "Labels": {
                                    "1:2": {
                                        "LabelId": "1:2",
                                        "Name": "白色"
                                    },
                                    "2:2": {
                                        "LabelId": "2:2",
                                        "Name": "XL"
                                    }
                                },
                                "MarketPrice": 100,
                                "GroupBuyingPrice": 90,
                                "SettlementPrice": 80,
                                "CostPrice": 70,
                                "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg"
                            }
                        }
                    },
                    "Detail": "古法酿造黑糖",
                    "Notes": "周二21:00截单，下周一配送",
                    "StartTime": "2018-01-12T11:02:16.178783+08:00",
                    "EndTime": "2018-01-13T11:02:16.178783+08:00",
                    "Version": 1
                },
                "Result": [ // 统计结果
                    {
                        "ItemId": "1:1-2:1",
                        "Labels": { // 标签
                            "1:1": {
                                "LabelId": "1:1",
                                "Name": "黑色"
                            },
                            "2:1": {
                                "LabelId": "2:1",
                                "Name": "L"
                            }
                        },
                        "MarketPrice": 100, // 市场价
                        "GroupBuyingPrice": 90, // 团购价
                        "SettlementPrice": 80, // 结算价
                        "CostPrice": 70, // 总结算价
                        "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg",//配图
                        "Sales": 300,//销量
                        "CommunityCount": 60,//社群数目
                        "TotalCost": 2000,//总成本
                        "TotalSettlement": 3000,//总结算金额
                        "TaskId": 1,//团购任务ID
                        "StartTime": "2018-01-12T11:02:16.178783+08:00",//开始时间
                        "EndTime": "2018-01-13T11:02:16.178783+08:00",//结束时间
                        "Title": "团购任务-1"//团购任务名称
                    },
                    {
                        "ItemId": "1:2-2:2",
                        "Labels": {
                            "1:2": {
                                "LabelId": "1:2",
                                "Name": "白色"
                            },
                            "2:2": {
                                "LabelId": "2:2",
                                "Name": "XL"
                            }
                        },
                        "MarketPrice": 100,
                        "GroupBuyingPrice": 90,
                        "SettlementPrice": 80,
                        "CostPrice": 70,
                        "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg",
                        "Sales": 200,
                        "CommunityCount": 300,
                        "TotalCost": 1000,
                        "TotalSettlement": 2000,
                        "TaskId": 1,
                        "StartTime": "2018-01-12T11:02:16.178783+08:00",
                        "EndTime": "2018-01-13T11:02:16.178783+08:00",
                        "Title": "团购任务-1"
                    }
                ],
                "Version": 1
            }
        ]
    },
    "Message": "success"
}
```



###/api/v1/group_buying/admin/line_community/list/:organization_id/:line_id

线路社群

```json
{
    "Code": 0,
    "Data": {
        "Count": 100,
        "List": [
            {
                "LineId": 1, // 线路ID
                "GroupId": 1, // 社群ID
                "LineName": "线路-1", // 线路名称
                "GroupName": "青阳社区一群", // 社群名称
                "ManagerUid": "3", // 管理员ID
                "ManagerName": "马云", // 管理员名称
                "ManagerMobile": "18676726608" // 管理员手机
            },
            {
                "LineId": 2,
                "GroupId": 2,
                "LineName": "线路-1",
                "GroupName": "青阳社区一群",
                "ManagerUid": "3",
                "ManagerName": "马云",
                "ManagerMobile": "18676726608"
            }
        ]
    },
    "Message": "success"
}
```

