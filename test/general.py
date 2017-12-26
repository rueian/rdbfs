#!/usr/local/bin/python3
# -*- coding: UTF-8 -*-

import os


def checkInit(path):
    if os.path.ismount(path):
        if not os.listdir(path):
            return True, 'success'
        else:
            return False, 'dir is not empty'
    else:
        return False, 'dir is not mount'


def initPath(path):
    filelist = os.listdir(path)
    for filename in filelist:
        filepath = os.path.join(path, filename)
        if os.path.isfile(filepath):
            os.remove(filepath)
        else:
            os.removedirs(filepath)
    pass


class bColor:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


class outputAnswer(object):
    def headerPrint(head):
        print("\nHEADER: {}{}{}".format(bColor.HEADER, head, bColor.ENDC))
        return True

    def itemPrint(item):
        print("> item: {}".format(item))
        return True

    def checkMsg(msg):
        print(">> {}PASS{} {}".format(bColor.OKGREEN, bColor.ENDC, msg))
        return True

    def passCheck(msg):
        print(">> {}PASS{} {}".format(bColor.OKBLUE, bColor.ENDC, msg))
        return True

    def failCheck(msg):
        print(">> {}FAIL{} {}".format(bColor.FAIL, bColor.ENDC, msg))
        return True


def printPWD():
    outputAnswer.itemPrint(os.getcwd())
    outputAnswer.checkMsg(os.getcwd())
    outputAnswer.passCheck(os.getcwd())
    outputAnswer.failCheck(os.getcwd())


def main():
    printPWD()


if __name__ == '__main__':
    main()
