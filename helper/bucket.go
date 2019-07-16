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
	                "AWS": "arn:aws:iam::xxxxxxxxxxxx:user/yourusername"
	            },
	            "Action": [
	                "s3:GetBucketLocation",
	                "s3:ListBucket"
	            ],
	            "Resource": "arn:aws:s3:::yourbucketname"
	        },
	        {
	            "Sid": "statement2",
	            "Effect": "Allow",
	            "Principal": {
	                "AWS": "arn:aws:iam::xxxxxxxxxxxx:user/yourusername"
	            },
	            "Action": "s3:GetObject",
	            "Resource": "arn:aws:s3:::yourbucketname/*"
	        }
	    ]
	}`
	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
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
