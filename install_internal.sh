#install all crd
make install
#load existing models and policies
cd deployments
python3 load_crd.py
#the webhook deployment uses this account
kubectl create serviceaccount my-sa # if necessary, namespace should be added 
kubectl apply -f webhook_register_internal_step1.yaml
kubectl apply -f webhook_register_internal_step2.yaml
#grant necessary access authority to the webhook
kubectl create clusterrolebinding my-sa-view \
  --clusterrole=cluster-admin \
  --serviceaccount=default:my-sa 
