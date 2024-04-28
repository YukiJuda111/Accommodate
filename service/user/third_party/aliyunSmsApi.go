package third_party

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"math/rand"
	"os"
	"strconv"
)

func SendSms(phoneNum string) (string, error) {
	//-------------------阿里云短信发送api-------------------
	config := sdk.NewConfig()

	// TODO: 在系统环境中配置ALIBABA_CLOUD_ACCESS_KEY_ID和ALIBABA_CLOUD_ACCESS_KEY_SECRET环境变量
	// Please ensure that the environment variables ALIBABA_CLOUD_ACCESS_KEY_ID and ALIBABA_CLOUD_ACCESS_KEY_SECRET are set.
	credential := credentials.NewAccessKeyCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))
	/* use STS Token
	credential := credentials.NewStsTokenCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"), os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN"))
	*/
	client, err := sdk.NewClientWithOptions("cn-hangzhou", config, credential)
	if err != nil {
		return "", err
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["SignName"] = "阿里云短信测试"
	request.QueryParams["TemplateCode"] = "SMS_154950909"
	request.QueryParams["PhoneNumbers"] = phoneNum

	smsCode := strconv.Itoa(rand.Intn(1000000))
	request.QueryParams["TemplateParam"] = "{\"code\":\"" + smsCode + "\"}"
	_, err = client.ProcessCommonRequest(request)
	return smsCode, err
}
