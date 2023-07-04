# notion-ical-syncer
Sync notion calendar to apple calendar

### How to use it

This project depends on AWS service, so you have to deploy it to AWS Lambda in order to use it.

I already created a script to build and deploy this project, you just need to run following commands.

```shell
    git@github.com:ronazst/notion-ical-syncer.git
    cd notion-ical-syncer
    #such as: ./scripts/build-and-deploy-to-aws.sh ronazst us-east-1
    ./scripts/build-and-deploy-to-aws.sh <customized-cloudformation-stack-suffix> <aws-region>
```

### Note

* Deploy AWS resource may charge you monthly
  * I already make the resource under free tier as possible as I can. But sadly, I can't avoid to use S3 inorder to simplify deployment process
* Each deployment(with different suffix) will create 2 different AWS Cloudformation stack
  * notion-ical-syncer-infra-<customized-cloudformation-stack-suffix>
  * notion-ical-syncer-<customized-cloudformation-stack-suffix>
* Once you want to delete it, go to AWS Cloudformation to delete related stack manually
* The S3 bucket will not be deleted after you deleted AWS Cloudformation, so you have to delete related S3 manually