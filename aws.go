package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"golang.org/x/net/websocket"
	"fmt"
)

func createInstance (tag string, ws *websocket.Conn) {
	ami:= "ami-0b33d91d" // us-east-1  -  Amazon Linux AMI 2016.09.1 (HVM), SSD
	instType:="t2.micro"

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})

	svc := ec2.New(sess)
	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String(ami),
		InstanceType: aws.String(instType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		var m Message
		m.Text = fmt.Sprintf("Could not create instance\n")
		postMessage(ws, m)
		log.Println("Could not create instance", err)
		return
	}

	var m Message
	m.Text = fmt.Sprintf("Created instance\n")
	postMessage(ws, m)

	log.Println("Created instance", *runResult.Instances[0].InstanceId)

	// Add tags to the created instance
	_ , errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(tag),
			},
		},
	})
	if errtag != nil {
		m.Text = fmt.Sprintf("Could not create tags for instance\n")
		postMessage(ws, m)
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return
	}

	//var m Message
	m.Text = "Successfully tagged instance"
	postMessage(ws, m)
	log.Println("Successfully tagged instance")
}
