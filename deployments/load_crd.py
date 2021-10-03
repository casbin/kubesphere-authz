
import os
os.system("kubectl create namespace policy")
workspacePath = os.path.abspath('../')
exampleDirs=[name for name in  os.listdir(workspacePath+"/example") if os.path.isdir(workspacePath+"/example/"+name)]
print(exampleDirs)
for rule in exampleDirs:
    ruleDir=workspacePath+"/example/"+rule
    # check whether the crd yaml files exists:
    ruleCrdDir=ruleDir+"/crd"
    if not os.path.exists( ruleCrdDir):
        continue
    # the definition yaml rule is supposed to be this
    rulePolicyDir=ruleCrdDir+"/policy"
    ruleDefinitionFileName=ruleCrdDir+"/"+rule+"_definition.yaml"
    policyCrdFiles=[rulePolicyDir+"/"+name for name in os.listdir(rulePolicyDir) if os.path.isfile(rulePolicyDir+"/"+name)]

    os.system("kubectl apply -f %s "%(ruleDefinitionFileName))
    for policyCrdFile in policyCrdFiles:
        os.system("kubectl apply -f %s"%(policyCrdFile))



