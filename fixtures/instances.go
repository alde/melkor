package fixtures

import (
	"encoding/json"
	"fmt"
	"time"
)

var (
	t = time.Date(2017, 02, 28, 10, 15, 22, 23, time.Local)
)

func createAWSResponse(i int, team string) map[string]interface{} {
	account := "123456789"

	return map[string]interface{}{
		"AmiLaunchIndex": i,
		"Architecture":   "x86_64",
		"BlockDeviceMappings": []map[string]interface{}{
			{
				"DeviceName": "/dev/sda1",
				"Ebs": map[string]interface{}{
					"AttachTime":          t,
					"DeleteOnTermination": true,
					"Status":              "attached",
					"VolumeId":            fmt.Sprintf("vol-%d", i),
				},
			},
		},
		"ClientToken":  fmt.Sprintf("smart-%s-2M3GIISBQQ3U", team),
		"EbsOptimized": false,
		"EnaSupport":   nil,
		"Hypervisor":   "xen",
		"IamInstanceProfile": map[string]string{
			"Arn": fmt.Sprintf("arn:aws:iam::%s:instance-profile/%s", account, team),
			"Id":  "KIRUNAVARA",
		},
		"ImageId":           fmt.Sprintf("ami-%d", i),
		"InstanceId":        fmt.Sprintf("i-%d", i),
		"InstanceLifecycle": nil,
		"InstanceType":      "m3.medium",
		"KernelId":          nil,
		"KeyName":           fmt.Sprintf("ssh_key_%s", team),
		"LaunchTime":        t,
		"Monitoring": map[string]string{
			"State": "disabled",
		},
		"NetworkInterfaces": []map[string]interface{}{
			{
				"Association": nil,
				"Attachment": map[string]interface{}{
					"AttachTime":          t,
					"AttachmentId":        fmt.Sprintf("eni-attach-%d", i),
					"DeleteOnTermination": true,
					"DeviceIndex":         0,
					"Status":              "attached",
				},
				"Description": "",
				"Groups": []map[string]interface{}{
					{
						"GroupId":   fmt.Sprintf("sg-%d", i),
						"GroupName": fmt.Sprintf("%s-eu-staging-ServerSecurityGroup-%d", team, i),
					},
				},
				"Ipv6Addresses":      []map[string]interface{}{},
				"MacAddress":         "02:a6:fd:3a:b7:4f",
				"NetworkInterfaceId": fmt.Sprintf("eni-%d", i),
				"OwnerId":            account,
				"PrivateDnsName":     fmt.Sprintf("ip-10-20-30-%d.eu-west-1.compute.internal", i),
				"PrivateIpAddress":   fmt.Sprintf("10.20.30.%d", i),
				"PrivateIpAddresses": []map[string]interface{}{
					{
						"Association":      nil,
						"Primary":          true,
						"PrivateDnsName":   fmt.Sprintf("ip-10-20-30-%d.eu-west-1.compute.internal", i),
						"PrivateIpAddress": fmt.Sprintf("10.20.30.%d", i),
					},
				},
				"SourceDestCheck": true,
				"Status":          "in-use",
				"SubnetId":        "subnet-11aa22bb",
				"VpcId":           "vpc-987654321",
			},
		},
		"Placement": map[string]interface{}{
			"Affinity":         nil,
			"AvailabilityZone": "eu-west-1b",
			"GroupName":        "",
			"HostId":           nil,
			"Tenancy":          "default",
		},
		"Platform":         nil,
		"PrivateDnsName":   fmt.Sprintf("ip-10-20-30-%d.eu-west-1.compute.internal", i),
		"PrivateIpAddress": fmt.Sprintf("10.20.30.%d", i),
		"ProductCodes":     []map[string]interface{}{},
		"PublicDnsName":    "",
		"PublicIpAddress":  nil,
		"RamdiskId":        nil,
		"RootDeviceName":   "/dev/sda1",
		"RootDeviceType":   "ebs",
		"SecurityGroups": []map[string]interface{}{
			{
				"GroupId":   fmt.Sprintf("sg-%d", i),
				"GroupName": fmt.Sprintf("%s-eu-staging-ServerSecurityGroup-%d", team, i),
			},
		},
		"SourceDestCheck":       true,
		"SpotInstanceRequestId": nil,
		"SriovNetSupport":       nil,
		"State": map[string]interface{}{
			"Code": 16,
			"Name": "running",
		},
		"StateReason":           nil,
		"StateTransitionReason": "",
		"SubnetId":              "subnet-11aa22bb",
		"Tags": []map[string]string{
			{
				"Key":   "Name",
				"Value": fmt.Sprintf("fake-service-%d", i),
			},
			{
				"Key":   "Service",
				"Value": "fake-service",
			},
			{
				"Key":   "Team",
				"Value": team,
			},
			{
				"Key":   "Environment",
				"Value": "staging",
			},
		},
		"VirtualizationType": "hvm",
		"VpcId":              "vpc-987654321",
	}
}

// FullCrawlerData returns the full structure of `count` number of instances
func FullCrawlerData(count int) []map[string]interface{} {
	var data []map[string]interface{}
	for i := 0; i < count; i++ {
		data = append(data, createAWSResponse(i, fmt.Sprintf("team%d", i)))
	}

	return data

}

// ExpectedFullResponse actually just marshals the FullCrawlerData into a json string
func ExpectedFullResponse(count int) string {
	b, _ := json.Marshal(FullCrawlerData(count))
	return string(b)
}
