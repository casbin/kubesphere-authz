from testframework.TestSuite import TestSuite
from testframework.Tester import Tester
import os,sys
workspacePath = os.path.abspath('../..')
print(workspacePath)
tester = Tester()

#test for allowed_repo_approved
tester.addTestSuite(TestSuite(
    "allowed_repo_approved",
    workspacePath+"/example/allowed_repo/config.json",
    [
        (workspacePath+"/example/allowed_repo/allowed_repo_approved.yaml", True),
        (workspacePath+"/example/allowed_repo/allowed_repo_rejected.yaml", False),
    ]
))
#test for container_resource_limit
tester.addTestSuite(TestSuite(
    "container_resource_limit",
    workspacePath+"/example/container_resource_limit/config.json",
    [
        (workspacePath+"/example/container_resource_limit/container_resource_limit_approved.yaml", True),
        (workspacePath+"/example/container_resource_limit/container_resource_limit_rejected1.yaml", False),
        (workspacePath+"/example/container_resource_limit/container_resource_limit_rejected2.yaml", False),
    ]
))
#test for permission check
tester.addTestSuite(TestSuite(
    "permission",
    workspacePath+"/example/permission/config.json",
    [
        (workspacePath+"/example/permission/permission_approved.yaml", True),
        (workspacePath+"/example/permission/permission_rejected.yaml", False),
    ]
))


#test for block_nodeport_service
tester.addTestSuite(TestSuite(
    "block_nodeport_service",
    workspacePath+"/example/block_nodeport_service/config.json",
    [
        (workspacePath+"/example/block_nodeport_service/block_nodeport_service_approved.yaml", True),
        (workspacePath+"/example/block_nodeport_service/block_nodeport_service_rejected.yaml", False),
    ]
))

res=tester.run()
if res==False:
    sys.exit(1)