#!/bin/bash
java -classpath h2.jar org.h2.tools.Server -?
java -classpath h2.jar org.h2.tools.Server -webAllowOthers -tcpAllowOthers -pgAllowOthers -ifNotExists -trace
