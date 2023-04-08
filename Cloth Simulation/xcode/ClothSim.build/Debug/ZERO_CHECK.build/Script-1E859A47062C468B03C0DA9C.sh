#!/bin/sh
set -e
if test "$CONFIGURATION" = "Debug"; then :
  cd /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p4-clothsim-sp22-almost-done/xcode
  make -f /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p4-clothsim-sp22-almost-done/xcode/CMakeScripts/ReRunCMake.make
fi
if test "$CONFIGURATION" = "Release"; then :
  cd /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p4-clothsim-sp22-almost-done/xcode
  make -f /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p4-clothsim-sp22-almost-done/xcode/CMakeScripts/ReRunCMake.make
fi
if test "$CONFIGURATION" = "MinSizeRel"; then :
  cd /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p4-clothsim-sp22-almost-done/xcode
  make -f /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p4-clothsim-sp22-almost-done/xcode/CMakeScripts/ReRunCMake.make
fi
if test "$CONFIGURATION" = "RelWithDebInfo"; then :
  cd /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p4-clothsim-sp22-almost-done/xcode
  make -f /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p4-clothsim-sp22-almost-done/xcode/CMakeScripts/ReRunCMake.make
fi

