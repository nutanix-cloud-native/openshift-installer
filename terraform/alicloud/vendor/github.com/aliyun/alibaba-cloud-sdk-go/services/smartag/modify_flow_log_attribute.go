package smartag

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyFlowLogAttribute invokes the smartag.ModifyFlowLogAttribute API synchronously
func (client *Client) ModifyFlowLogAttribute(request *ModifyFlowLogAttributeRequest) (response *ModifyFlowLogAttributeResponse, err error) {
	response = CreateModifyFlowLogAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyFlowLogAttributeWithChan invokes the smartag.ModifyFlowLogAttribute API asynchronously
func (client *Client) ModifyFlowLogAttributeWithChan(request *ModifyFlowLogAttributeRequest) (<-chan *ModifyFlowLogAttributeResponse, <-chan error) {
	responseChan := make(chan *ModifyFlowLogAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyFlowLogAttribute(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ModifyFlowLogAttributeWithCallback invokes the smartag.ModifyFlowLogAttribute API asynchronously
func (client *Client) ModifyFlowLogAttributeWithCallback(request *ModifyFlowLogAttributeRequest, callback func(response *ModifyFlowLogAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyFlowLogAttributeResponse
		var err error
		defer close(result)
		response, err = client.ModifyFlowLogAttribute(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ModifyFlowLogAttributeRequest is the request struct for api ModifyFlowLogAttribute
type ModifyFlowLogAttributeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	NetflowVersion       string           `position:"Query" name:"NetflowVersion"`
	Description          string           `position:"Query" name:"Description"`
	InactiveAging        requests.Integer `position:"Query" name:"InactiveAging"`
	SlsRegionId          string           `position:"Query" name:"SlsRegionId"`
	ActiveAging          requests.Integer `position:"Query" name:"ActiveAging"`
	OutputType           string           `position:"Query" name:"OutputType"`
	ProjectName          string           `position:"Query" name:"ProjectName"`
	LogstoreName         string           `position:"Query" name:"LogstoreName"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	NetflowServerPort    requests.Integer `position:"Query" name:"NetflowServerPort"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	NetflowServerIp      string           `position:"Query" name:"NetflowServerIp"`
	Name                 string           `position:"Query" name:"Name"`
	FlowLogId            string           `position:"Query" name:"FlowLogId"`
}

// ModifyFlowLogAttributeResponse is the response struct for api ModifyFlowLogAttribute
type ModifyFlowLogAttributeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyFlowLogAttributeRequest creates a request to invoke ModifyFlowLogAttribute API
func CreateModifyFlowLogAttributeRequest() (request *ModifyFlowLogAttributeRequest) {
	request = &ModifyFlowLogAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Smartag", "2018-03-13", "ModifyFlowLogAttribute", "smartag", "openAPI")
	request.Method = requests.POST
	return
}

// CreateModifyFlowLogAttributeResponse creates a response to parse from ModifyFlowLogAttribute response
func CreateModifyFlowLogAttributeResponse() (response *ModifyFlowLogAttributeResponse) {
	response = &ModifyFlowLogAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}