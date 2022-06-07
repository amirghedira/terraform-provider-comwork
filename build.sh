#!/bin/bash
echo "=> deleting previous version"
rm -rf /Users/ghediraamir/.terraform.d/plugins/terraform.local/local/comwork/*
echo "=> deleting terraform setup"
rm -rf .terraform
rm -rf terraform.*
rm -rf *tfstate*

rm .terraform.lock.hcl

VERSION=${1}
echo "=> building terraform-provider-comwork_${VERSION}"
go build -o terraform-provider-comwork_${VERSION}
chmod +x terraform-provider-comwork_${VERSION}
mkdir -p /Users/ghediraamir/.terraform.d/plugins/terraform.local/local/comwork/${VERSION}
mkdir -p /Users/ghediraamir/.terraform.d/plugins/terraform.local/local/comwork/${VERSION}/darwin_amd64
mv terraform-provider-comwork_${VERSION} /Users/ghediraamir/.terraform.d/plugins/terraform.local/local/comwork/${VERSION}/darwin_amd64/terraform-provider-comwork_${VERSION}

echo "=> Initializing terraform-provider-comwork_${VERSION}"
terraform init
