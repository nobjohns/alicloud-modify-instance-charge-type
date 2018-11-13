package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/howeyc/gopass"
	"os"
	"strings"
)

var (
	ecsClient        *ecs.Client
	instanceIds      = make([]string, 0)
	period           string
	periodUnit       string
	duration         = 12
	chargeType       string
	includeDataDisks = true
	dryRun           = false
	autoPay          = true
	region           string
)

func main() {

	region, key, secret := GetCreds()
	ecsClient = createEcsClient(region, key, secret)
	instanceIds, period, periodUnit, chargeType = getParams()
	modifyInstanceChargeType()
	ModifyInstanceAutoRenewAttribute()
}

func GetCreds() (string, string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter Alicloud region: ")
	r, _ := reader.ReadString('\n')
	r = strings.TrimSpace(r)

	fmt.Printf("Alicloud access key: ")
	k, _ := reader.ReadString('\n')
	k = strings.TrimSpace(k)

	fmt.Printf("Alicloud access secret: ")
	secretByte, err := gopass.GetPasswd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s := string(secretByte)
	return r, k, s
}

func createEcsClient(region, key, secret string) *ecs.Client {
	ecsClient, err := ecs.NewClientWithAccessKey(
		region,
		key,
		secret)
	if err != nil {
		panic(err)
	}
	return ecsClient
}

func getParams() ([]string, string, string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter instance IDs, separated by comma: ")
	i, _ := reader.ReadString('\n')
	for _, v := range strings.Split(i, ",") {
		instanceIds = append(instanceIds, strings.TrimSpace(v))
	}

	fmt.Printf("Period(1/2/3 ...): ")
	period, _ := reader.ReadString('\n')
	period = strings.TrimSpace(period)

	fmt.Printf("Unit of Period(Year/Month): ")
	periodUnit, _ := reader.ReadString('\n')
	periodUnit = strings.TrimSpace(periodUnit)

	fmt.Printf("Instance Charge Type(PrePaid/PostPaid): ")
	chargeType, _ := reader.ReadString('\n')
	chargeType = strings.TrimSpace(chargeType)

	return instanceIds, period, periodUnit, chargeType
}

func modifyInstanceChargeType() {
	request := ecs.CreateModifyInstanceChargeTypeRequest()
	request.Period = requests.Integer(period)
	request.PeriodUnit = periodUnit
	request.InstanceChargeType = chargeType
	request.IncludeDataDisks = requests.NewBoolean(includeDataDisks)
	request.DryRun = requests.NewBoolean(dryRun)
	request.AutoPay = requests.NewBoolean(autoPay)

	a, _ := json.Marshal(instanceIds)
	request.InstanceIds = string(a)
	response, err := ecsClient.ModifyInstanceChargeType(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}

func ModifyInstanceAutoRenewAttribute() {
	request := ecs.CreateModifyInstanceAutoRenewAttributeRequest()
	request.RegionId = region
	request.Duration = requests.NewInteger(duration)
	request.AutoRenew = requests.NewBoolean(autoPay)

	for _, i := range instanceIds {
		request.InstanceId = i
		response, err := ecsClient.ModifyInstanceAutoRenewAttribute(request)
		if err != nil {
			panic(err)
		}
		fmt.Println(response)
	}
}
