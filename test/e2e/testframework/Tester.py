from . import TestSuite
from typing import *
class Tester:
    def __init__(self):
        self.testSuiteList=[]
        self.passed=0
        self.failed=0
    def addTestSuite(self, testSuite: TestSuite):
        self.testSuiteList.append(testSuite)
    def run(self)->bool:
        for i in self.testSuiteList:
            res=i.run()
            self.passed+=res[0]
            self.failed+=res[1]
        print("==================================================================")
        print("[E2E Test]: Total %d PASSED, %d FAILED"%(self.passed, self.failed))
        return self.failed==0