package examples

import (
	"testing"
)



func TestGVM(t *testing.T) {
	// 登录
	gmpClient, err := gvmClient("192.168.102.137:9390", "admin", "webpassword")
	if err != nil {
		panic(err)
	}
	//获取scanner  08b69003-5fc2-4037-a479-93b440211c73

	//创建新资产：使用自带的端口列表（传入id）
	//assetsId, err := createAssetsFromExistPort(gmpClient,"测试域名-Ping存活", "localhost","4a4717fe-57d2-11e1-9a26-406186ea4fc5")
	////c8724451-6e6e-4a8e-975e-b8bf39f3d453 <nil>
	//fmt.Println(assetsId, err)

	//创建新资产：自定义扫描端口
	//assetsId,err := createAssetsFromNewPort(gmpClient,"测试域名-自定义端口", "localhost","T:7120,7106,7103,7117,7110,139,7102,80,7104,7100,3306,445,22")
	////ca0faed7-504b-4242-9918-c1b4afc1eaeb <nil>
	//fmt.Println(assetsId, err)
	//
	////获取 规则id
	//configId, err := getConfigId(gmpClient, "system_config")
	//// 92946eff-b8f4-40ba-acb5-662adf17441c <nil>
	//fmt.Println(configId, err)
	//
	////添加任务
	//reportId, err := createTask(
	//	gmpClient,
	//	"测试任务1",
	//	"92946eff-b8f4-40ba-acb5-662adf17441c",
	//	"ca0faed7-504b-4242-9918-c1b4afc1eaeb",
	//	"08b69003-5fc2-4037-a479-93b440211c73")
	//
	//// 5f9d2242-faf6-4c6f-82c7-b49a3fa2a6e1 <nil>
	//fmt.Println(reportId, err)
	//
	////启动任务
	//success, err := startTask(gmpClient, "5f9d2242-faf6-4c6f-82c7-b49a3fa2a6e1")
	//fmt.Println(success, err)

	//// 获取全部结果
	getAllResult(gmpClient)
}