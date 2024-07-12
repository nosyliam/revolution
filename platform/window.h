#include "base.h"
#include <stdio.h>

extern AXError _AXUIElementGetWindow(AXUIElementRef, CGWindowID* out);

typedef struct Frame {
	int width;
    int height;
    int x;
    int y;
} Frame;

typedef struct Frames {
	int len;
	Frame* frames;
} Frames;

typedef struct Screenshot {
    size_t width;
	size_t height;
	size_t stride;

	unsigned int len;
	unsigned char* data;
} Screenshot;

#if defined(IS_MACOSX)

typedef struct Window {
	AXUIElementRef window;
	CGWindowID     id;
} Window;

Frames* get_display_frames() {
	CGDirectDisplayID displays[32];
	uint32_t count;

    if (CGGetActiveDisplayList(32, displays, &count) != kCGErrorSuccess)
    {
        return NULL;
    }


    Frames* frames = malloc(sizeof(Frames));
    Frame* data = malloc(sizeof(Frame) * count);
    frames->len = count;
    frames->frames = data;

    for (int i = 0; i < count; i++) {
    	CGRect bounds = CGDisplayBounds(displays[i]);
    	Frame *frame = malloc(sizeof(Frame));
    	frame->width = (int)bounds.size.width;
    	frame->height = (int)bounds.size.height;
    	frame->x = bounds.origin.x;
    	frame->y = bounds.origin.y;
    	data[i] = *frame;
    }

    return frames;
}

void set_window_frame(const Window* window, const int width, const int height, const int x, const int y) {
	CGPoint* position = malloc(sizeof(CGPoint));
	position->x = x;
	position->y = y;
	CFTypeRef positionStorage = AXValueCreate(kAXValueCGPointType, position);
	AXUIElementSetAttributeValue(window->window, (CFStringRef)kAXPositionAttribute, positionStorage);
	CFRelease(positionStorage);
	free(position);
	
	CGSize* size = malloc(sizeof(CGSize));
	size->width = width;
	size->height = height;
	CFTypeRef sizeStorage = AXValueCreate(kAXValueCGSizeType, size);
	AXUIElementSetAttributeValue(window->window, (CFStringRef)kAXSizeAttribute, sizeStorage);
	CFRelease(sizeStorage);
	free(size);
}

Frame* get_window_frame(const Window* window) {
	CFTypeRef sizeStorage;
	AXError result = AXUIElementCopyAttributeValue(window->window, (CFStringRef)kAXSizeAttribute, &sizeStorage);

	CGSize size;
	if (result == kAXErrorSuccess) {
		if (!AXValueGetValue(sizeStorage, kAXValueCGSizeType, (void *)&size)) {
			size = CGSizeZero;
		}
	}
	else {
		size = CGSizeZero;
	}

	if (sizeStorage)
		CFRelease(sizeStorage);

	CFTypeRef positionStorage;
	result = AXUIElementCopyAttributeValue(window->window, (CFStringRef)kAXPositionAttribute, &positionStorage);

	CGPoint topLeft;
	if (result == kAXErrorSuccess) {
		if (!AXValueGetValue(positionStorage, kAXValueCGPointType, (void *)&topLeft)) {
			topLeft = CGPointZero;
		}
	}
	else {
		topLeft = CGPointZero;
	}

	if (positionStorage)
		CFRelease(positionStorage);

	Frame* frame = malloc(sizeof(Frame));
	frame->width = size.width;
	frame->height = size.height;
	frame->x = topLeft.x;
	frame->y = topLeft.y;

	return frame;
}

Screenshot* screenshot_window(const Window* window) {
    int bgraDataLen = 0;
    CGImageRef windowImage = CGWindowListCreateImage(CGRectNull, kCGWindowListOptionIncludingWindow, window->id, kCGWindowImageBoundsIgnoreFraming & kCGWindowImageNominalResolution);
	CGColorSpaceRef colorSpace = CGColorSpaceCreateDeviceRGB();
	size_t width = CGImageGetWidth(windowImage);
	size_t height = CGImageGetHeight(windowImage);
	size_t stride = CGImageGetBytesPerRow(windowImage);
	size_t len = sizeof(unsigned char) * stride * height;
	unsigned char* data = malloc(len);

	CGContextRef context = CGBitmapContextCreate(data, width, height,
												 8, stride, colorSpace,
											     kCGImageAlphaPremultipliedLast | kCGBitmapByteOrder32Big);

	CGContextDrawImage(context, CGRectMake(0, 0, width, height), windowImage);

	Screenshot* screenshot = malloc(sizeof(Screenshot));
	screenshot->width = width;
	screenshot->height = height;
	screenshot->stride = stride;
	screenshot->data = data;
	screenshot->len = len;

	CGImageRelease(windowImage);
	CGContextRelease(context);
    return screenshot;
}

void activate_window(const Window* window) {
	if (AXUIElementPerformAction(window->window, kAXRaiseAction) != kAXErrorSuccess) {
		pid_t pid = 0;
		if (AXUIElementGetPid(window->window, &pid) != kAXErrorSuccess || !pid) { return; }

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"

		ProcessSerialNumber psn;
		if (GetProcessForPID(pid, &psn) == 0) {
			SetFrontProcessWithOptions(&psn, kSetFrontProcessFrontWindowOnly);
		}

#pragma clang diagnostic pop
	}
}

int get_window_count(pid_t pid) {
	AXUIElementRef application = AXUIElementCreateApplication(pid);
	if (application == 0) {return 0;}

	CFArrayRef windows = NULL;
	AXUIElementCopyAttributeValues(application, kAXWindowsAttribute, 0, 1024, &windows);

	if (windows != NULL) {
	    CFRelease(windows);
	    CFRelease(application);
		return CFArrayGetCount(windows);
    }
    CFRelease(application);
    return 0;
}

Window* get_window_with_pid(pid_t pid) {
	AXUIElementRef application = AXUIElementCreateApplication(pid);
	if (application == 0) {return NULL;}

	CFArrayRef windows = NULL;
	// Get all windows associated with the app
	AXUIElementCopyAttributeValues(application, kAXWindowsAttribute, 0, 1024, &windows);
	CGWindowID win = 0;

	if (windows != NULL) {
		int count = CFArrayGetCount(windows);
		if (count == 1) {
            AXUIElementRef element = (AXUIElementRef) CFArrayGetValueAtIndex(windows, 0);

			CGWindowID temp = 0;
			_AXUIElementGetWindow(element, &temp);

            CFRetain(element);
            CFRelease(windows);
            Window* window = malloc(sizeof(window));
            window->window = element;
            window->id = temp;
            return window;
		}
	}

	return NULL;
};
#endif