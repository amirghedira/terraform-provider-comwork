#!/bin/bash
echo "=> deleting previous version"
rm -rf /Users/ghediraamir/.terraform.d/plugins/terraform.local/local/myprovider/*
echo "=> deleting terraform setup"
rm -rf .terraform
rm .terraform.lock.hcl

VERSION=${1}
echo "=> building terraform-provider-myprovider_${VERSION}"
go build -o terraform-provider-myprovider_${VERSION}
chmod +x terraform-provider-myprovider_${VERSION}
mkdir /Users/ghediraamir/.terraform.d/plugins/terraform.local/local/myprovider/${VERSION}
mkdir /Users/ghediraamir/.terraform.d/plugins/terraform.local/local/myprovider/${VERSION}/darwin_amd64
mv terraform-provider-myprovider_${VERSION} /Users/ghediraamir/.terraform.d/plugins/terraform.local/local/myprovider/${VERSION}/darwin_amd64/terraform-provider-myprovider_${VERSION}

echo "=> Initializing terraform-provider-myprovider_${VERSION}"
terraform init
