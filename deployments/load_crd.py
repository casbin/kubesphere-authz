
import os
import time
time.sleep(10)
workspacePath = os.path.abspath('../')
print("Load Policy start")
os.system("kubectl create namespace kubesphere-authz-system")
os.system("kubectl apply -f %s"%(workspacePath+"/config/crd/bases/auth.casbin.org_casbinmodels.yaml"))
exampleDirs=[name for name in  os.listdir(workspacePath+"/example") if os.path.isdir(workspacePath+"/example/"+name)]
print(exampleDirs)
for rule in exampleDirs:
    ruleDir=workspacePath+"/example/"+rule
    # check whether the crd yaml files exists:
    ruleCrdDir=ruleDir+"/crd"
    if not os.path.exists( ruleCrdDir):
        continue
    # the definition yaml rule is supposed to be this
    # the model yaml rule is supposed to be this
    ruleModelFileName=ruleCrdDir+"/"+rule+"_model.yaml"
    rulePolicyDir=ruleCrdDir+"/policy"
    policyCrdFile=[]
    if os.path.exists(rulePolicyDir):
        policyCrdFiles=[rulePolicyDir+"/"+name for name in os.listdir(rulePolicyDir) if os.path.isfile(rulePolicyDir+"/"+name)]

    os.system("kubectl apply -f %s "%(ruleModelFileName))
    time.sleep(0.5)
    for policyCrdFile in policyCrdFiles:
        os.system("kubectl apply -f %s"%(policyCrdFile))
