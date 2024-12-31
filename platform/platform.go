package platform

//#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations
//#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit
//#cgo darwin LDFLAGS: -framework Carbon -framework CoreFoundation
//#cgo darwin LDFLAGS: -L. -L./capture_darwin -lCapture -framework Accelerate -framework UniformTypeIdentifiers -framework Foundation -framework ScreenCaptureKit -framework CoreMedia -framework CoreVideo
import "C"
