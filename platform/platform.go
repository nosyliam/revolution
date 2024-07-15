package platform

//#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations
//#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit
//#cgo darwin LDFLAGS: -framework Carbon -framework CoreFoundation
import "C"
