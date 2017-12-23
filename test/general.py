#!/usr/local/bin/python3
# -*- coding: UTF-8 -*-

import os


def initPath(path):
    filelist = os.listdir(path)
    for filename in filelist:
        filepath = os.path.join(path, filename)
        if os.path.isfile(filepath):
            os.remove(filepath)
        else:
            os.removedirs(filepath)
    pass
