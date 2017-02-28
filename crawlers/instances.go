package crawlers

import (
	"time"

	"github.com/alde/melkor/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

type ec2Client interface {
	DescribeInstances(*ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
}

// The InstancesCrawler struct holds the implementation for the interface
type InstancesCrawler struct {
	instances   []*ec2.Instance
	lastCrawled time.Time
	count       int
	config      *config.Config
	client      ec2Client
}

// NewInstancesCrawler is the constructor of this crawler
func NewInstancesCrawler(c *config.Config) *InstancesCrawler {
	sess := session.Must(session.NewSession())

	client := ec2.New(sess, &aws.Config{Region: aws.String(c.AWSRegion)})
	return &InstancesCrawler{
		config: c,
		client: client,
	}
}

// Resource identifies the name of the crawled resource
func (i *InstancesCrawler) Resource() string {
	return "Instances"
}

// LastCrawled is the timestamp of the most recent crawl
func (i *InstancesCrawler) LastCrawled() time.Time {
	return i.lastCrawled
}

// DoCrawl handles the crawling of AWS
func (i *InstancesCrawler) DoCrawl() error {
	logrus.WithField("resource", i.Resource()).Info("Crawling")

	params := &ec2.DescribeInstancesInput{}
	resp, err := i.client.DescribeInstances(params)

	if err != nil {
		return err
	}

	for _, r := range resp.Reservations {
		for _, ins := range r.Instances {
			i.instances = append(i.instances, ins)
		}
	}
	i.count = len(i.instances)
	i.lastCrawled = time.Now()

	logrus.WithFields(logrus.Fields{
		"resource": i.Resource(),
		"count":    i.Count(),
	}).Info("Done crawling")

	return nil
}

// List instances
func (i *InstancesCrawler) List() []string {
	var data []string
	for _, ins := range i.instances {
		data = append(data, aws.StringValue(ins.InstanceId))
	}
	return data
}

// ListExpanded expands the result
func (i *InstancesCrawler) ListExpanded() []map[string]interface{} {
	var data []map[string]interface{}
	for _, ins := range i.instances {
		data = append(data, structs.Map(ins))
	}
	return data
}

// Get returns a single instance by id
func (i *InstancesCrawler) Get(id string) map[string]interface{} {
	for _, ins := range i.instances {
		if *ins.InstanceId == id {
			return structs.Map(ins)
		}
	}
	return nil
}

// Count the number of instances crawled
func (i *InstancesCrawler) Count() int {
	return i.count
}
