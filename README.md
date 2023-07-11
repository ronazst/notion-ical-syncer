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
    ./scripts/get-ical-url.sh <customized-cloudformation-stack-suffix> <aws-region>
```

After the `get-ical-url.sh` command, you will see the url like this: `https://<id>.lambda-url.us-sourtheast-1.on.aws/<uuid>`

Then open `https://<id>.lambda-url.us-sourtheast-1.on.aws/<uuid>/add-config` to add new config

After you config has been added you will get a config id which can used to query ical content.

The query format like this: `https://<id>.lambda-url.us-sourtheast-1.on.aws/<uuid>/ical?config_ids=<config_id_1>&&config_ids=<config_id_2>`

### Note

* Deploy AWS resource may charge you monthly
  * I already make the resource under free tier as possible as I can. But sadly, I can't avoid to use S3 inorder to simplify deployment process
* Each deployment(with different suffix) will create 2 different AWS Cloudformation stack
  * notion-ical-syncer-infra (this shared between all lambda stack)
  * notion-ical-syncer-stack-<customized-cloudformation-stack-suffix>
* Once you want to delete it, go to AWS Cloudformation to delete related stack manually
* The S3 bucket will not be deleted after you deleted AWS Cloudformation, so you have to delete related S3 manually