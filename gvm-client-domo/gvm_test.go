package examples

import (
	"fmt"
	"testing"
	"time"
)



func TestGVM(t *testing.T) {
	// 登录
	gmpClient, err := gvmClient("192.168.102.137:9390", "admin", "strongpassword")
	if err != nil {
		panic(err)
	}
	//获取scanner
	scannerList, err := getScannersIdList(gmpClient)
	scannerId := scannerList[1]

	//创建新资产：使用自带的端口列表（传入默认端口的id）
	assetsId, err := createAssetsFromExistPort(gmpClient,"测试域名-Ping存活", "localhost","4a4717fe-57d2-11e1-9a26-406186ea4fc5")
	fmt.Println(assetsId, err)

	//创建新资产：自定义扫描端口
	//assetsId,err := createAssetsFromNewPort(gmpClient,"测试域名-自定义端口", "localhost","T:7120,7106,7103,7117,7110,139,7102,80,7104,7100,3306,445,22")
	//fmt.Println(assetsId, err)

	//提前在创建好要扫描的规则  根据规则名称获取 规则id
	configId, err := getConfigId(gmpClient, "system_config")
	fmt.Println(configId, err)

	//添加任务
	taskId, err := createTask(gmpClient, "测试任务", configId, assetsId, scannerId)
	fmt.Println(taskId, err)

	//启动任务
	success, err := startTask(gmpClient, taskId)
	fmt.Println(success, err)

	// 获取所有任务结果
	getAllResult(gmpClient)


	// 循环获取任务进度 status == "-1"代表完成
	for {
		status, _ := getTaskProcess(gmpClient,taskId)
		time.Sleep(10 * time.Second)
		if status == "-1"{
			// 获取单个任务的结果
			getSingleTaskResult(gmpClient, taskId)
			break
		}
	}
}
