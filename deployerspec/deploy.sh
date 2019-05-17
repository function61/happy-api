#!/bin/bash -eu

# the zip name needs to change from previous deployment for it to be considered new
newZipName="lambdafunc-$FRIENDLY_REV_ID.zip"

if [ ! -e "$newZipName" ]; then
	ln -s "lambdafunc.zip" "$newZipName"
fi

echo "zip_filename = \"$newZipName\"" > terraform.tfvars

terraform init

planFilename="update.plan"

terraform plan -state /state/terraform.tfstate -out "$planFilename"

# wait for enter
echo "[press any key to deploy or ctrl-c to abort]"
read DUMMY

terraform apply "$planFilename"
