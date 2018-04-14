import os
import glob
from datetime import datetime

from fabric.api import *
from fabric.contrib import files
from fabric.utils import abort


def clean():
    local('mvn clean')


def generate_ast():
    compile()
    local('cd target/classes && java com/christilden/craftinginterpreters/tool/GenerateAst ../../src/main/java/com/christilden/lox')


def compile():
    local('mvn compile')


def package():
    local('mvn package')


def test():
    package()
    local('java -jar target/java-interp-1.0-SNAPSHOT.jar test.lox')


def run():
    package()
    local('java -jar target/java-interp-1.0-SNAPSHOT.jar')
