package models

import "time"

type CloudtrailData struct {
	Records []Records `json:"Records"`
}

type Records struct {
	EventVersion      string            `json:"eventVersion"`
	UserIdentity      UserIdentity      `json:"userIdentity"`
	EventTime         time.Time         `json:"eventTime"`
	EventSource       string            `json:"eventSource"`
	EventName         string            `json:"eventName"`
	AwsRegion         string            `json:"awsRegion"`
	SourceIPAddress   string            `json:"sourceIPAddress"`
	UserAgent         string            `json:"userAgent"`
	RequestParameters RequestParameters `json:"requestParameters"`
	ResponseElements  ResponseElements  `json:"responseElements"`
}

type UserIdentity struct {
	Type        string `json:"type"`
	PrincipalID string `json:"principalId"`
	Arn         string `json:"arn"`
	AccessKeyID string `json:"accessKeyId"`
	AccountID   string `json:"accountId"`
	UserName    string `json:"userName"`
}

type RequestParameters struct {
	InstancesSet InstanceSet `json:"instancesSet"`
}

type InstanceSet struct {
	Items []Items `json:"items"`
}

type Items struct {
	InstanceID    string        `json:"instanceId"`
	CurrentState  CurrentState  `json:"currentState"`
	PreviousState PreviousState `json:"previousState"`
}

type CurrentState struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type PreviousState struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type ResponseElements struct {
	InstancesSet InstanceSet `json:"instancesSet"`
}

type Ip2CountryResponse struct {
	CountryName string `json:"countryName"`
}

type Country struct {
	Name   Name   `json:"name"`
	Region string `json:"region"`
}

type Name struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}
