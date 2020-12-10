package examples

import (
	"errors"
	"fmt"
	"github.com/filewalkwithme/go-gmp/pkg/9/gmp"
	"github.com/filewalkwithme/go-gmp/pkg/9/gmp/client"
	"github.com/filewalkwithme/go-gmp/pkg/9/gmp/connections"
)

// 构建客户端连接
func gvmClient(gvmServerAddr string, gvmUsername string, gvmPwd string)(gmp.Client, error){
	// 连接到 GVMD
	conn, err := connections.NewTLSConnection(gvmServerAddr,true)
	if err != nil {
		return nil, err
	}
	//defer conn.Close()
	// 认证
	gmpClient := client.New(conn)
	auth := &gmp.AuthenticateCommand{}
	auth.Credentials.Username = gvmUsername
	auth.Credentials.Password = gvmPwd
	_, err = gmpClient.Authenticate(auth)
	if err != nil {
		return nil, err
	}
	return gmpClient, nil
}

// 获取scanner id list
// 要去掉默认带的CVE scanner，所以也可以去数据库取数据
func getScannersIdList(gmpClient gmp.Client) ([]string,error) {
	scannerCmd := &gmp.GetScannersCommand{}
	//scannerCmd.Filter = `name="OpenVAS Default"`
	scannersResp, err := gmpClient.GetScanners(scannerCmd)
	if err != nil {
		return nil, err
	}
	if scannersResp != nil {
		var scannerIdList []string
		for _, scanner:= range scannersResp.Scanner {
			scannerIdList = append(scannerIdList,scanner.ID)
		}
		return scannerIdList, nil
	}
	return nil, nil
}

// 获取 规则id 规则名称为system_config
func getConfigId(gmpClient gmp.Client, configName string)(string,error){
	configCmd := &gmp.GetConfigsCommand{}
	//configCmd.Filter = `name="system_config"`
	configCmd.Filter = fmt.Sprintf(fmt.Sprintf("name=\"%s\"", configName))
	configResp, err := gmpClient.GetConfigs(configCmd)
	if err != nil {
		return "", err
	}
	if len(configResp.Config) > 0 {
		return configResp.Config[0].ID, nil
	}
	return "", errors.New("规则未创建")
}

// 添加资产，第一步先把所有资产录入，资产名称和域名列表（可以一个名称，对应多个资产）
func createAssetsFromExistPort(gmpClient gmp.Client, assetsName string, assetsHosts string, useExistPortListID string)(string,error){
	assetsCmd := &gmp.CreateTargetCommand{}
	assetsCmd.Name = assetsName
	assetsCmd.Hosts = assetsHosts
	//assetsCmd.AliveTests = "ICMP Ping"
	portList := &gmp.CreateTargetPortList{ID: useExistPortListID}
	assetsCmd.PortList = portList
	createTargetResp, err := gmpClient.CreateTarget(assetsCmd)
	if err != nil {
		panic(err)
	}
	return createTargetResp.ID,nil
}

// 添加资产，第一步先把所有资产录入，资产名称和域名列表（可以一个名称，对应多个资产）
func createAssetsFromNewPort(gmpClient gmp.Client, assetsName string, assetsHosts string, newPortRange string)(string,error){
	assetsCmd := &gmp.CreateTargetCommand{}
	assetsCmd.Name = assetsName
	assetsCmd.Hosts = assetsHosts
	assetsCmd.PortRange = newPortRange
	createTargetResp, err := gmpClient.CreateTarget(assetsCmd)
	if err != nil {
		panic(err)
	}
	return createTargetResp.ID,nil
}

// 添加任务
func createTask(gmpClient gmp.Client, taskName string, configId string, assetsId string, scannerId string )(string, error) {
	newTask := &gmp.CreateTaskCommand{}
	newTask.Name = taskName
	config := &gmp.CreateTaskConfig{ID: configId}
	newTask.Config = config
	target := &gmp.CreateTaskTarget{ID: assetsId}
	newTask.Target = target
	scanner := &gmp.CreateTaskScanner{ID: scannerId}
	newTask.Scanner = scanner
	newTaskResp, err := gmpClient.CreateTask(newTask)
	if err != nil {
		return "", err
	}
	return newTaskResp.ID, nil
}

func startTask(gmpClient gmp.Client, taskId string)(string, error){
	startTask:= &gmp.StartTaskCommand{}
	startTask.TaskID = taskId
	taskResp, err := gmpClient.StartTask(startTask)
	if err != nil {
		return "", err
	}
	return taskResp.ReportID, nil
}

func getAllResult(gmpClient gmp.Client)  {
	getResults := &gmp.GetResultsCommand{}
	getResults.Filter = "severity>3.9 apply_overrides=0 min_qod=70 sort-reverse=severity rows=10 first=1"
	results, err := gmpClient.GetResults(getResults)
	if err != nil {
		panic(err)
	}
	res := results.Result
	//Show results
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i])
		//fmt.Printf("Result[%d]: (host:%s) (vul:%s) (score: %s)\n", i, res[i].Host,res[i].Name, res[i].Severity)
	}
}

func getSingleTaskResult(gmpClient gmp.Client, taskId string){
	getResults := &gmp.GetResultsCommand{}
	getResults.Filter = fmt.Sprintf("task_id=%s levels=hml min_qod=70 apply_overrides=0 sort-reverse=severity rows=10 first=1",taskId)
	//getResults.Filter = "task_id=8025e8a5-05e9-43f7-81da-f3e45d7f2d32 levels=hml min_qod=70 apply_overrides=0 sort-reverse=severity rows=10 first=1"
	results, err := gmpClient.GetResults(getResults)
	if err != nil {
		panic(err)
	}
	res := results.Result
	for i := 0; i < len(res); i++ {
		resSingle := res[i]
		fmt.Println(resSingle)
		//fmt.Printf("Result[%d]: (host:%s) (vul:%s) (score: %s)\n", i, res[i].Host,res[i].Name, res[i].Severity)
	}
}

func getTaskProcess(gmpClient gmp.Client, taskId string)(string,error){
	gt := &gmp.GetTasksCommand{}
	gt.TaskID = taskId
	getTasksResp, err := gmpClient.GetTasks(gt)
	if err != nil {
		return "", err
	}
	return getTasksResp.Task[0].Progress.Value, nil
}
