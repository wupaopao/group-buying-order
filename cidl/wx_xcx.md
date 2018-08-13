### /group_buying/wx_xcx/task/community_sales/:task_Id/:group_id

团长查看团购任务的销量

```json
{
    "Code": 0,
    "Data": {
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
            "Items": {
                "1": {
                    "ItemId": "1",
                    "Name": "颜色",
                    "LabelIdCounter": 2,
                    "Labels": {
                        "1:1": {
                            "LabelId": "1:1",
                            "Name": "黑色"
                        },
                        "1:2": {
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
            "SkuMap": {
                "1:1-2:1": {
                    "SkuId": "1:1-2:1",
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
                    "SkuId": "1:1-2:2",
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
                    "SkuId": "1:2-2:1",
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
                    "SkuId": "1:2-2:2",
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
        "StartTime": "2018-01-12T17:41:10.2016351+08:00",
        "EndTime": "2018-01-13T17:41:10.2016351+08:00",
        "Version": 1,
        "IsAutoPost": true,
        "GroupState": 0,
        "OrderState": 0,
        "IsDelete": false,
        "Sales": { // 销量
            "1:1-2:1": { // 商品规格
                "SkuId": "1:1-2:1",
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
                "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg",
                "Sales": 100 // 销量
            },
            "1:1-2:2": {
                "SkuId": "1:1-2:2",
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
                "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg",
                "Sales": 100
            },
            "1:2-2:1": {
                "SkuId": "1:2-2:1",
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
                "IllustrationPicture": "https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=2109908802,1486768288&fm=27&gp=0.jpg",
                "Sales": 100
            },
            "1:2-2:2": {
                "SkuId": "1:2-2:2",
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
                "Sales": 100
            }
        }
    },
    "Message": "success"
}
```



###/api/v1/group_buying/wx_xcx/task/community_buy/1/3

```json
{
	"Items":[
		{
			"SkuId":"1:1-2:1", # sku_id
			"Sales":100 # 购买数量
		},
		{
			"SkuId":"1:2-2:1",
			"Sales":50
		}
	]
}
```

