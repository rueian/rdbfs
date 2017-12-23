#!/usr/local/bin/python3
# -*- coding: UTF-8 -*-

import os
import shutil


def editFile(path, filename, inputSting='123\n'):
    filepath = os.path.join(path, filename)
    with open(filepath, 'a', encoding='UTF-8') as openfile:
        openfile.write(inputSting)
    return True


def reEditFile(path, filename, inputSting='123\n'):
    filepath = os.path.join(path, filename)
    try:
        with open(filepath, 'w') as openfile:
            openfile.write(inputSting)
        return True, 'success'
    except Exception as e:
        return False, e


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
    pass


def moveFile(path, topath):
    pass
