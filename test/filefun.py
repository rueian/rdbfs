#!/usr/local/bin/python3
# -*- coding: UTF-8 -*-

import os
import shutil


def editFile(path, filename, inputSting='123\n'):
    filepath = os.path.join(path, filename)
    try:
        with open(filepath, 'a', encoding='UTF-8') as openfile:
            openfile.write(inputSting)
    except Exception as e:
        return False, e
    return True, "success"


def reEditFile(path, filename, inputSting='123\n'):
    filepath = os.path.join(path, filename)
    try:
        with open(filepath, 'w') as openfile:
            openfile.write(inputSting)
    except Exception as e:
        return False, e
    return True, "success"


def renameFile(path, srcfile, dstfile):
    srcpath = os.path.join(path, srcfile)
    dstpath = os.path.join(path, dstfile)
    try:
        os.rename(srcpath, dstpath)
    except Exception as e:
        return False, e
    return True, "success"


def removeFile(path, filename):
    filepath = os.path.join(path, filename)
    try:
        os.remove(filepath)
    except Exception as e:
        return False, e
    return True, "success"


def cmpDiffTwoFile(firstfile, secondfile, detail=False):
    if detail:
        commandline = 'cmp -l '+firstfile+' '+secondfile
    else:
        commandline = 'cmp '+firstfile+' '+secondfile
    status = os.system(commandline)
    if status:
        return True
    else:
        return False


def copyFile(path, topath):
    if os.path.isfile(path) and not os.path.isfile(topath):
        try:
            shutil.copyfile(path, topath)
            return True, "success"
        except Exception as e:
            return False, e


def copyFileToFolder(path, folderpath):
    if os.path.isfile(path) and os.path.isdir(folderpath):
        try:
            shutil.copy(path, folderpath)
            return True, "success"
        except Exception as e:
            return False, e


def checkFileExist(path, filename):
    filepath = os.path.join(path, filename)
    if os.path.isfile(filepath):
        return True, "Check success"
    return False, "File is not exist"


def checkFileRemove(path, filename):
    filepath = os.path.join(path, filename)
    if not os.path.isfile(filepath):
        return True, "Check success"
    return False, "File is not removed"


def checkCopyFileSuccess(path, topath):
    if os.path.isfile(path):
        if os.path.isfile(path):
            return True, "Check success"
        else:
            return False, "{} is not exist".format(topath)
    else:
        return False, "{} is not exist".format(path)


def moveFile(path, topath):
    pass
