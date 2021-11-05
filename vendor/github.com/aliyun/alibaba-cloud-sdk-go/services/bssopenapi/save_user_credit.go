package bssopenapi

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

// SaveUserCredit invokes the bssopenapi.SaveUserCredit API synchronously
func (client *Client) SaveUserCredit(request *SaveUserCreditRequest) (response *SaveUserCreditResponse, err error) {
	response = CreateSaveUserCreditResponse()
	err = client.DoAction(request, response)
	return
}

// SaveUserCreditWithChan invokes the bssopenapi.SaveUserCredit API asynchronously
func (client *Client) SaveUserCreditWithChan(request *SaveUserCreditRequest) (<-chan *SaveUserCreditResponse, <-chan error) {
	responseChan := make(chan *SaveUserCreditResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveUserCredit(request)
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

// SaveUserCreditWithCallback invokes the bssopenapi.SaveUserCredit API asynchronously
func (client *Client) SaveUserCreditWithCallback(request *SaveUserCreditRequest, callback func(response *SaveUserCreditResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveUserCreditResponse
		var err error
		defer close(result)
		response, err = client.SaveUserCredit(request)
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

// SaveUserCreditRequest is the request struct for api SaveUserCredit
type SaveUserCreditRequest struct {
	*requests.RpcRequest
	AvoidExpiration          requests.Boolean `position:"Query" name:"AvoidExpiration"`
	Description              string           `position:"Query" name:"Description"`
	AvoidPrepaidNotification requests.Boolean `position:"Query" name:"AvoidPrepaidNotification"`
	AvoidPrepaidExpiration   requests.Boolean `position:"Query" name:"AvoidPrepaidExpiration"`
	AvoidNotification        requests.Boolean `position:"Query" name:"AvoidNotification"`
	Operator                 string           `position:"Query" name:"Operator"`
	CreditValue              string           `position:"Query" name:"CreditValue"`
	CreditType               string           `position:"Query" name:"CreditType"`
}

// SaveUserCreditResponse is the response struct for api SaveUserCredit
type SaveUserCreditResponse struct {
	*responses.BaseResponse
	Code      string `json:"Code" xml:"Code"`
	Success   bool   `json:"Success" xml:"Success"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Message   string `json:"Message" xml:"Message"`
}

// CreateSaveUserCreditRequest creates a request to invoke SaveUserCredit API
func CreateSaveUserCreditRequest() (request *SaveUserCreditRequest) {
	request = &SaveUserCreditRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("BssOpenApi", "2017-12-14", "SaveUserCredit", "", "")
	request.Method = requests.POST
	return
}

// CreateSaveUserCreditResponse creates a response to parse from SaveUserCredit response
func CreateSaveUserCreditResponse() (response *SaveUserCreditResponse) {
	response = &SaveUserCreditResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}