#!/usr/local/bin/python3
# -*- coding: UTF-8 -*-

import os
import general
import filefun

from general import outputAnswer


def getPWD():
    outputAnswer.itemPrint(os.getcwd())
    return os.getcwd()


def msgOutput(item, status):
    outputAnswer.itemPrint(item)
    if status[0] is True:
        outputAnswer.passCheck(status[1])
    else:
        outputAnswer.failCheck(status[1])


def itemOutput(item, path, filename, status):
    outputAnswer.itemPrint(item)
    if status[0] is True:
        status = filefun.checkFileRemove(path, filename)
        if status[0] is True:
            outputAnswer.passCheck(status[1])
        else:
            outputAnswer.failCheck(status[1])
    else:
        outputAnswer.failCheck(status[1])


def main():
    initpath = getPWD()
    status = general.checkInit(initpath)
    msgOutput("checkInit", status)
    status = filefun.editFile(initpath, 'test1', 'I am good guy\n')
    msgOutput("editFile", status)
    status = filefun.renameFile(initpath, 'test1', 'test2')
    msgOutput("renameFile", status)
    status = filefun.removeFile(initpath, 'test2')
    itemOutput("removeFile", initpath, 'test2', status)


if __name__ == '__main__':
    main()
