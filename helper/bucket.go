package helper

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func PutBucketPolicy(bucket string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast2RegionID),
	}))

	svc := s3.New(sess)
	policy := `{
	    "Version": "2012-10-17",
	    "Statement": [
	        {
	            "Sid": "statement1",
	            "Effect": "Allow",
	            "Principal": {
	                "AWS": "arn:aws:iam::200869506108:user/rino"
	            },
	            "Action": [
	                "s3:GetBucketLocation",
	                "s3:ListBucket"
	            ],
	            "Resource": "arn:aws:s3:::karina-cohive-backup-us-east-2"
	        },
	        {
	            "Sid": "statement2",
	            "Effect": "Allow",
	            "Principal": {
	                "AWS": "arn:aws:iam::200869506108:user/rino"
	            },
	            "Action": "s3:GetObject",
	            "Resource": "arn:aws:s3:::karina-cohive-backup-us-east-2/*"
	        }
	    ]
	}`
	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		// Policy: aws.String("{\"Version\": \"2012-10-17\", \"Statement\": [{ \"Sid\": \"id-1\",\"Effect\": \"Allow\",\"Principal\": {\"AWS\": \"arn:aws:iam::123456789012:root\"}, \"Action\": [ \"s3:PutObject\",\"s3:PutObjectAcl\"], \"Resource\": [\"arn:aws:s3:::acl3/*\" ] } ]}"),
		Policy: aws.String(policy),
	}

	result, err := svc.PutBucketPolicy(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
