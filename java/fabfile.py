import os
import glob
from datetime import datetime

from fabric.api import *
from fabric.contrib import files
from fabric.utils import abort


def clean():
    local('mvn clean')


def package():
    local('mvn package')


def run():
    package()
    local('java -jar target/java-interp-1.0-SNAPSHOT.jar')
