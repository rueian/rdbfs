#!/usr/local/bin/python3
# -*- coding: UTF-8 -*-

import os
import shutil


def createFolder(path, srcfolder):
    folderpath = os.path.join(path, srcfolder)
    try:
        # os.mkdir(folderpath, mode=0777)
        os.mkdir(folderpath)
        return True, 'success'
    except Exception as e:
        return False, e


def createFolders(path, *srcfolder):
    for sub in srcfolder:
        path = os.path.join(path, sub)
    try:
        # os.makedirs(path, mode=0777)
        os.makedirs(path)
        return True, 'success'
    except Exception as e:
        return False, e


def renameFolder(path, srcfolder, detfolder):
    srcfolderpath = os.path.join(path, srcfolder)
    detfolderpath = os.path.join(path, detfolder)
    if os.path.isdir(srcfolderpath) and not os.path.isdir(detfolderpath):
        try:
            os.rename(srcfolderpath, detfolderpath)
        except Exception as e:
            return False, e
        return True, "success"
    else:
        return False, "{} is not exist".format(srcfolderpath)


def copyWholeFolder(path, srcfolder, dstfolder):
    srcfolderpath = os.path.join(path, srcfolder)
    dstfolderpath = os.path.join(path, dstfolder)
    if os.path.isdir(srcfolderpath) and os.path.isdir(dstfolderpath):
        dstfolderpath = os.path.join(dstfolderpath, srcfolder)
        try:
            shutil.copytree(srcfolderpath, dstfolderpath)
        except Exception as e:
            return False, e
        return True, 'success'
    else:
        return False, "{} or {} is not exist".format(srcfolder, dstfolder)


def changeFolderMode(path, srcfolder, mode=0o644):
    folderpath = os.path.join(path, srcfolder)
    if os.path.isdir(folderpath):
        try:
            os.chmod(folderpath, mode)
        except Exception as e:
            return False, e
        return True, "success"
    else:
        return False, "{} is not exist".format(srcfolder)


def removeFolder(path, srcfolder):
    folderpath = os.path.join(path, srcfolder)
    if os.path.isdir(folderpath):
        try:
            os.rmdir(folderpath)
            return True, "success"
        except Exception as e:
            return False, e
    return False, "{} is not exist".format(folderpath)


def removeWholeFolder(path, srcfolder):
    folderpath = os.path.join(path, srcfolder)
    if os.path.isdir(folderpath):
        try:
            shutil.rmtree(folderpath)
            return True, "success"
        except Exception as e:
            return False, e
    return False, "{} is not exist".format(folderpath)


def checkFolderExist(path, srcfolder):
    if type(srcfolder) is list:
        for val in srcfolder:
            path = os.path.join(path, val)
        folderpath = path
    else:
        folderpath = os.path.join(path, srcfolder)
    if os.path.isdir(folderpath):
        return True, "Check success"
    return False, "{} is not exist".format(folderpath)


def checkFolderNotExist(path, srcfolder):
    folderpath = os.path.join(path, srcfolder)
    if not os.path.isdir(folderpath):
        return True, "Check success"
    return False, "{} is exist".format(folderpath)


def checkRenameFolder(path, srcfolder, dstfolder):
    if checkFolderExist(path, dstfolder)[0]:
        if checkFolderNotExist(path, srcfolder)[0]:
            return True, "Check success"
        else:
            return False, "{} is exist".format(srcfolder)
    else:
        return False, "{} is not exist".format(dstfolder)


def checkFolderMode(path, srcfolder, mustbe):
    folderpath = os.path.join(path, srcfolder)
    if os.path.isdir(folderpath):
        # folderMode = oct(stat.S_IMODE(os.stat(folderpath).st_mode))
        folderMode = oct(os.stat(folderpath).st_mode)[-3:]
        if folderMode == mustbe:
            return True, "Check success > folder mode is {}".format(folderMode)
        else:
            return False, "Folder mode is {}".format(folderMode)
    return False, "{} is not exist".format(folderpath)
