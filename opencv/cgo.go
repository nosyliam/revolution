package opencv

/*
#cgo windows CXXFLAGS:   --std=c++11
#cgo windows CPPFLAGS:   -IC:./include
#cgo windows LDFLAGS:    -LC: -L./lib -lopencv_core4100 -lopencv_imgproc4100 -lopencv_features2d4100
#cgo darwin CXXFLAGS:   --std=c++11
#cgo darwin CPPFLAGS:   -I./include/opencv4
#cgo darwin LDFLAGS:    -L./lib/opencv4/3rdparty -llibpng -littnotify -ltegra_hal -lzlib -L./lib -lopencv_core -lopencv_imgproc -lopencv_features2d -lopencv_flann -lopencv_imgcodecs
*/
import "C"
