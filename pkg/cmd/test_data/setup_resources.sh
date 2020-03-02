#!/bin/bash

set -e


echo "setting up the cloud resources for ecluster $CLUSTER_NAME in project $PROJECT_ID"

# CLI-DOC-GEN-START
gcloud iam service-accounts create $CLUSTER_NAME-jb

retry gcloud iam service-accounts add-iam-policy-binding \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:$PROJECT_ID.svc.id.goog[$NAMESPACE/externaldns-sa]" \
  $CLUSTER_NAME-ex@$PROJECT_ID.iam.gserviceaccount.com
# CLI-DOC-GEN-END

# change to the new jx namespace
jx ns $NAMESPACE