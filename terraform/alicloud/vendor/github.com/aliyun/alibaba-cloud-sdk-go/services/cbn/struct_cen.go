package cbn

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

// Cen is a nested struct in cbn response
type Cen struct {
	CenId                  string                 `json:"CenId" xml:"CenId"`
	Name                   string                 `json:"Name" xml:"Name"`
	Description            string                 `json:"Description" xml:"Description"`
	ProtectionLevel        string                 `json:"ProtectionLevel" xml:"ProtectionLevel"`
	Status                 string                 `json:"Status" xml:"Status"`
	CreationTime           string                 `json:"CreationTime" xml:"CreationTime"`
	Ipv6Level              string                 `json:"Ipv6Level" xml:"Ipv6Level"`
	ResourceGroupId        string                 `json:"ResourceGroupId" xml:"ResourceGroupId"`
	CenBandwidthPackageIds CenBandwidthPackageIds `json:"CenBandwidthPackageIds" xml:"CenBandwidthPackageIds"`
	Tags                   TagsInDescribeCens     `json:"Tags" xml:"Tags"`
}