## OpenCV Bindings

This directory contains stripped-down bindings extracted from GoCV. You must build OpenCV from source and install to this directory.

### Build Instructions

MacOS:

```bash
cmake -G "Unix Makefiles" -DCMAKE_OSX_DEPLOYMENT_TARGET=11.0 -DWITH_JPEG=OFF -DWITH_TIFF=OFF -DWITH_WEBP=OFF -DWITH_OPENJPEG=OFF -DWITH_JASPER=OFF -DWITH_OPENEXR=OFF -DWITH_FFMPEG=OFF -DWITH_GSTREAMER=OFF -DWITH_MSFMF=OFF -DBUILD_LIST=features2d,imgcodecs,flann -DWITH_OPENCL=OFF -DBUILD_SHARED_LIBS=OFF -DCMAKE_BUILD_TYPE=Release -DOPENCV_VS_VERSIONINFO_SKIP=1 -DCMAKE_INSTALL_PREFIX=~/go/src/github.com/nosyliam/revolution/opencv ..```