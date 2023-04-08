#!/bin/sh
set -e
if test "$CONFIGURATION" = "Debug"; then :
  cd /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p3-2-pathtracer-sp22-whatissleep-3-2/xcode
  make -f /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p3-2-pathtracer-sp22-whatissleep-3-2/xcode/CMakeScripts/ReRunCMake.make
fi
if test "$CONFIGURATION" = "Release"; then :
  cd /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p3-2-pathtracer-sp22-whatissleep-3-2/xcode
  make -f /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p3-2-pathtracer-sp22-whatissleep-3-2/xcode/CMakeScripts/ReRunCMake.make
fi
if test "$CONFIGURATION" = "MinSizeRel"; then :
  cd /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p3-2-pathtracer-sp22-whatissleep-3-2/xcode
  make -f /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p3-2-pathtracer-sp22-whatissleep-3-2/xcode/CMakeScripts/ReRunCMake.make
fi
if test "$CONFIGURATION" = "RelWithDebInfo"; then :
  cd /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p3-2-pathtracer-sp22-whatissleep-3-2/xcode
  make -f /Users/imankonjo2020/Documents/CS_184_Foundation_of_Computer_Graphics_and_Imaging/p3-2-pathtracer-sp22-whatissleep-3-2/xcode/CMakeScripts/ReRunCMake.make
fi

