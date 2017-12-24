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


def editFile(initpath, srcfile):
    outputAnswer.headerPrint('EDIT FILE')
    status = filefun.editFile(initpath, srcfile, 'I am good guy\n')
    msgOutput("editFile", status)
    status = filefun.removeFile(initpath, srcfile)
    msgOutput("removeFile", status)
    return True


def renameFile(initpath, srcfile, dstfile):
    outputAnswer.headerPrint('REMOVE FILE')
    status = filefun.editFile(initpath, srcfile, 'I am good guy\n')
    msgOutput("editFile", status)
    status = filefun.renameFile(initpath, srcfile, dstfile)
    msgOutput("renameFile", status)
    status = filefun.removeFile(initpath, dstfile)
    itemOutput("removeFile", initpath, dstfile, status)
    return True


def editFileAndEditAgain(initpath, srcfile):
    outputAnswer.headerPrint('EDIT FILE AND EDIT AGAIN')
    status = filefun.editFile(initpath, srcfile, 'I am good guy\n')
    msgOutput("editFile", status)
    status = filefun.reEditFile(initpath, srcfile, 'actually not')
    msgOutput("reEditFile", status)
    status = filefun.removeFile(initpath, srcfile)
    itemOutput("removeFile", initpath, srcfile, status)
    return True


def editFileAndEditAgainAndRename(initpath, srcfile, dstfile):
    outputAnswer.headerPrint('EDIT FILE AND EDIT AGAIN AND RENAME')
    status = filefun.editFile(initpath, srcfile, 'I am good guy\n')
    msgOutput("editFile", status)
    status = filefun.reEditFile(initpath, srcfile, 'actually not')
    msgOutput("reEditFile", status)
    status = filefun.renameFile(initpath, srcfile, dstfile)
    msgOutput("renameFile", status)
    status = filefun.removeFile(initpath, dstfile)
    itemOutput("removeFile", initpath, dstfile, status)
    return True


def main():
    initpath = getPWD()
    status = general.checkInit(initpath)
    msgOutput("checkInit", status)
    editFile(initpath, 'test1')
    renameFile(initpath, 'test1', 'test2')
    editFileAndEditAgain(initpath, 'test1')
    editFileAndEditAgainAndRename(initpath, 'test1', 'test2')


if __name__ == '__main__':
    main()
