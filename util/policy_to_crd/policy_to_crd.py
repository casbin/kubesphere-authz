import sys
#before running this script, please make sure that there is a directory called 'crd', which has a subordinate directory policy
def generateDefinitionTemeplate(csvFileNameNoExt):
    singular=csvFileNameNoExt
    plural=singular#?

    singular_splits=singular.split("_")
    for i in range(0,len(singular_splits)):
        if len(singular_splits[i])!=0 and  singular_splits[i][0].isalpha():
            singular_splits[i]=singular_splits[i].capitalize()
    kind="".join(singular_splits)
    singular=singular.replace("_","-")
    plural=plural.replace("_","-")
    #load template
    templateFile=open("crd_template.yaml","r")
    template=templateFile.read()
    res=template.replace(r'${singular}',singular)
    res=res.replace(r'${plural}',plural)
    res=res.replace(r'${kind}',kind)

    return res,kind

def generatePoicyTemplate(csvFileName,kind):
    #load template
    templateFile=open("policy_template.yaml")
    template=templateFile.read()
    
    f=open(csvFileName,"r")
    lineNum=1
    while True:
        line=f.readline()
        if not line:
            break
        if line.strip()=="":
            continue
        name="policy"+str(lineNum)
        res=template.replace(r"${kind}",kind)
        res=res.replace(r"${name}",name)
        res=res.replace(r"${policy}",line.strip())
        out=open("crd/policy/"+name+".yaml","w")
        out.write(res)
        out.close()
        lineNum+=1

    return res
        
        
    


def main():
    if len(sys.argv)<=1:
        print("Fatal: no input files")
        return
    inputfile=sys.argv[1]
    splits=inputfile.split(".")
    print(splits[-1])
    if len(splits)==1 or splits[-1] != "csv":
        print("Fatal: input policy file must be in csv format")
        return

    #generate crd definition yaml
    definition,kind=generateDefinitionTemeplate(splits[0])
    out=open("crd/"+splits[0]+"_definition"+".yaml","w")
    out.write(definition)
    out.close()

    generatePoicyTemplate(inputfile,kind)

if __name__=="__main__":
    main()
