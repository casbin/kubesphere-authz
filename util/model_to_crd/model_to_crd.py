import sys

def generateModelYamlFromTemplate(modelFileNameNoExt):
    plural=modelFileNameNoExt.replace("_","-")

    templateFile=open("model_item_template.yaml","r")
    template=templateFile.read()
    res=template.replace(r'${plural}',plural)

    modelFile=open(modelFileNameNoExt+".conf","r")
    model=modelFile.read()
    model=model.replace("\n","\n    ")
    res=res.replace(r'${text}',model)
    return res



def main():
    if len(sys.argv)<=1:
        print("Fatal: no input file")
        return
    modelFileName=sys.argv[1]
    splits=modelFileName.split(".")
    if len(splits)==1 or splits[-1] != "conf":
        print("Fatal: input file should be casbin conf file")
        return
    modelName=splits[0]
    result=generateModelYamlFromTemplate(modelName)

    with open(modelName+"_model.yaml","w") as f:
        f.write(result)
    

    
if __name__=="__main__":
    main()
