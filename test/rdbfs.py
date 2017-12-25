#!/usr/local/bin/python3
# -*- coding: UTF-8 -*-

import os
import general
import filefun
import folderfun

from general import outputAnswer
from time import sleep


def getPWD():
    outputAnswer.itemPrint(os.getcwd())
    return os.getcwd()


def msgOutput(item, status):
    outputAnswer.itemPrint(item)
    if status[0] is True:
        outputAnswer.passCheck(status[1])
    else:
        outputAnswer.failCheck(status[1])


def msgFinalOutput(status):
    if status[0] is True:
        outputAnswer.checkMsg(status[1])
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


def createFolder(initpath, srcfolder):
    outputAnswer.headerPrint('CREATE FOLDER')
    status = folderfun.createFolder(initpath, srcfolder)
    msgOutput("createFolder", status)
    finalstatus = folderfun.checkFolderExist(initpath, srcfolder)
    status = folderfun.removeFolder(initpath, srcfolder)
    msgOutput("removeFolder", status)
    msgFinalOutput(finalstatus)
    return True


def renameFolder(initpath, srcfolder, dstfolder):
    outputAnswer.headerPrint('RENAME FOLDER')
    status = folderfun.createFolder(initpath, srcfolder)
    msgOutput("createFolder", status)
    sleep(2)
    status = folderfun.renameFolder(initpath, srcfolder, dstfolder)
    msgOutput("renameFolder", status)
    finalstatus = folderfun.checkRenameFolder(initpath, srcfolder, dstfolder)
    status = folderfun.removeFolder(initpath, dstfolder)
    msgOutput("removeFolder", status)
    msgFinalOutput(finalstatus)
    return True


def changeFolderMode(initpath, srcfolder, Mode=0o644):
    outputAnswer.headerPrint('CHANGE FOLDER MODE')
    status = folderfun.createFolder(initpath, srcfolder)
    msgOutput("createFolder", status)
    status = folderfun.checkFolderMode(initpath, srcfolder, '755')
    msgOutput("checkFolderMode", status)
    sleep(2)
    status = folderfun.changeFolderMode(initpath, srcfolder, 0o644)
    msgOutput("changeFolderMode", status)
    finalstatus = folderfun.checkFolderMode(initpath, srcfolder, '644')
    status = folderfun.removeFolder(initpath, srcfolder)
    msgOutput("removeFolder", status)
    msgFinalOutput(finalstatus)
    return True


def main():
    initpath = getPWD()
    status = general.checkInit(initpath)
    msgOutput("checkInit", status)
    # file test
    editFile(initpath, 'test1')
    renameFile(initpath, 'test1', 'test2')
    editFileAndEditAgain(initpath, 'test1')
    editFileAndEditAgainAndRename(initpath, 'test1', 'test2')
    # folder test
    createFolder(initpath, 'testF1')
    renameFolder(initpath, 'testF1', 'testF2')
    changeFolderMode(initpath, 'testF1')


if __name__ == '__main__':
    main()