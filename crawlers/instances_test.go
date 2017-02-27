package crawlers

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alde/melkor/config"
	"github.com/alde/melkor/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

func Test_Resource(t *testing.T) {
	c := &config.Config{}
	ic := NewInstancesCrawler(c)

	assert.Equal(t, ic.Resource(), "Instances")
}

func Test_Count(t *testing.T) {
	ic := &InstancesCrawler{
		count: 5,
	}

	assert.Equal(t, ic.Count(), 5)
}

func Test_LastCrawled(t *testing.T) {
	ic := &InstancesCrawler{
		lastCrawled: time.Now(),
	}

	assert.False(t, ic.LastCrawled().IsZero())
}

func setupCrawler() *InstancesCrawler {
	return &InstancesCrawler{
		count: 3,
		instances: []*ec2.Instance{
			{
				InstanceId:       aws.String("i-0"),
				PrivateIpAddress: aws.String("10.0.0.1"),
			},
			{
				InstanceId:       aws.String("i-1"),
				PrivateIpAddress: aws.String("10.0.0.2"),
			},
			{
				InstanceId:       aws.String("i-2"),
				PrivateIpAddress: aws.String("10.0.0.3"),
			},
		},
	}
}

func Test_List_Expanded(t *testing.T) {
	ic := setupCrawler()

	actual := ic.List(0, true).([]map[string]interface{})
	aLen := len(actual)
	assert.Equal(t, aLen, 3)
	for i := 0; i < aLen; i++ {
		assert.Equal(
			t,
			aws.StringValue(actual[i]["PrivateIpAddress"].(*string)),
			fmt.Sprintf("10.0.0.%d", i+1))
	}
}

func Test_List_Unexpanded(t *testing.T) {
	ic := setupCrawler()

	actual := ic.List(0, false).([]string)
	aLen := len(actual)
	assert.Equal(t, aLen, 3)
	for i := 0; i < aLen; i++ {
		assert.Equal(t, actual[i], fmt.Sprintf("i-%d", i))
	}
}

func Test_List_Empty(t *testing.T) {
	ic := &InstancesCrawler{count: 0}

	actual := ic.List(0, false).([]string)
	assert.Empty(t, actual, "Should be empty")
}

func Test_List_Expanded_Limit(t *testing.T) {
	ic := setupCrawler()

	actual := ic.List(1, true).([]map[string]interface{})
	aLen := len(actual)
	assert.Equal(t, aLen, 1)
	for i := 0; i < aLen; i++ {
		assert.Equal(
			t,
			aws.StringValue(actual[i]["PrivateIpAddress"].(*string)),
			fmt.Sprintf("10.0.0.%d", i+1))
	}
}

func Test_List_Unexpanded_Limit(t *testing.T) {
	ic := setupCrawler()

	actual := ic.List(1, false).([]string)
	aLen := len(actual)
	assert.Equal(t, aLen, 1)
	for i := 0; i < aLen; i++ {
		assert.Equal(t, actual[i], fmt.Sprintf("i-%d", i))
	}
}

func Test_Get(t *testing.T) {
	ic := setupCrawler()

	actual := ic.Get("i-0")

	assert.Equal(
		t,
		aws.StringValue(actual["PrivateIpAddress"].(*string)),
		"10.0.0.1")
}

func Test_Get_NotFound(t *testing.T) {
	ic := setupCrawler()

	actual := ic.Get("i-5")
	assert.Nil(t, actual)
}

func Test_DoCrawl(t *testing.T) {
	c := &config.Config{}
	mc := &mock.EC2Client{}
	ic := &InstancesCrawler{
		config: c,
		client: mc,
	}

	err := ic.DoCrawl()
	assert.Nil(t, err)

	assert.Equal(t, ic.Count(), 2)
	assert.False(t, ic.LastCrawled().IsZero())
}

func Test_DoCrawl_Fail(t *testing.T) {
	c := &config.Config{}
	mc := &mock.EC2Client{
		DescribeInstancesFn: func(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
			return nil, errors.New("something went terribly wrong")
		},
	}
	ic := &InstancesCrawler{
		config: c,
		client: mc,
	}

	err := ic.DoCrawl()
	assert.NotNil(t, err)
}
