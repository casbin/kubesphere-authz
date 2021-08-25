import subprocess
import os
import time
from typing import *


class TestSuite:
    def __init__(self, testName: str, configFilePath: str, testCases: List[Tuple[str, bool]]):
        '''
            constructor of the Testsuite case\n
            params:\n
            testName: name of the test suite, usually the rule name\n
            configFilePath: path of the json config file which shall be passed to webhook server\n
            testCases:List[Tuple[str, bool]], in which\n
                1st element of the tuple should be absolute path of the yaml file to apply to k8s
                2st element of the tuple should be boolean value, marking wherther this yaml should be accepted
        '''
        self.name = testName
        self.configFilePath = configFilePath
        self.testPath = os.path.dirname(__file__)
        self.workspacePath = os.path.abspath(self.testPath+"/../../..")
        self.webhookProcess = None
        self.logFileHandler = None
        self.testCases = testCases

    def setUp(self) -> None:
        '''
        set up external webhook server. Logs of output will be put into ../testlog
        '''

        #print("[E2E Test]:%s: setting up webhook server and log"%(self.name))
        logFileName = self.name+"-"+time.strftime("%Y-%m-%d-%H-%M-%S")+".log"
        self.logFileHandler = open(
            "%s/test/e2e/testlog/%s" % (self.workspacePath, logFileName), "w")
        cmd = [
            "%s/test/e2e/testbuild/main.exe" % (self.workspacePath),
            self.configFilePath,
        ]
        self.webhookProcess = subprocess.Popen(
            cmd, cwd=self.workspacePath, stdout=self.logFileHandler, stderr=self.logFileHandler)
        #print("[E2E Test]:%s: admission webhook started, pid %d"%(self.name,self.webhookProcess.pid))

    def tearDown(self) -> None:
        '''
        shut down external webhook server. 
        '''
        #print("[E2E Test]:%s: shutting webhook server and log"%(self.name))
        self.webhookProcess.kill()
        self.logFileHandler.close()

    def test(self) -> Tuple[int, int]:
        '''
        test each testcase and collect the result\n
        return: Tuple[int,int],in which:\n
            1st value of tuple is the number of passed test\n
            2st value of tuple is the number of failed test\n
        '''
        success = 0
        fail = 0
        time.sleep(1)
        for i in range(0, len(self.testCases)):
            yamlFilePath = self.testCases[i][0]
            shouldSuccess = self.testCases[i][1]

            webhookRunning = self.webhookProcess.poll()
            if webhookRunning != None:
                # webhook server crashed, immediately failed
                print("[E2E Test]:FAILED WEBHOOK HAS CRASHED. Test suit %s, Test case %s" % (
                    self.name, os.path.basename(yamlFilePath)))
                fail += 1
                continue
            cmd = [
                "minikube",
                "kubectl",
                "--",
                "apply",
                "-f",
                yamlFilePath,
                "--dry-run=server"
            ]
            res = subprocess.Popen(
                cmd, stdout=self.logFileHandler, stderr=self.logFileHandler)
            res.wait()
            if (shouldSuccess and res.returncode == 0) or (not shouldSuccess and res.returncode != 0):
                # passed
                print("[E2E Test]:PASSED Test suit %s, Test case %s" %
                      (self.name, os.path.basename(yamlFilePath)))
                success += 1
            else:
                # failed
                print("[E2E Test]:FAILED Test suit %s, Test case %s" %
                      (self.name, os.path.basename(yamlFilePath)))
                fail += 1
        return (success, fail)

    def run(self) -> Tuple[int, int]:
        '''
        run whole testsuite and collect the result\n
        return: Tuple[int,int],in which:\n
            1st value of tuple is the number of passed test\n
            2st value of tuple is the number of failed test\n
        '''
        self.setUp()
        res = self.test()
        self.tearDown()
        #print("[E2E Test]:%s: %d test passed, %d test failed"%(self.name,res[0],res[1]))
        return res
