#!/usr/local/bin/python3
# -*- coding: UTF-8 -*-

import os


def createFolder(path, subpath):
    folderpath = os.path.join(path, subpath)
    try:
        os.mkdir(folderpath, 0777)
        return True, 'success'
    except Exception as e:
        return False, e


def createFolders(path, *subpath):
    for sub in subpath:
        path = os.path.join(path, sub)
    try:
        os.makedirs(path, 0777)
        return True, 'success'
    except Exception as e:
        return False, e


def renameFolder(srcpath, dstpath):

    pass


def deleteFolder(srcpath):
    try:
        os.removedirs(srcpath)
        return True, "success"
    except Exception as e:
        return False, e
