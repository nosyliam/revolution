#include "base.h"

#if defined(IS_MACOSX)

#include <mach/mach.h>
#include <mach/mach_time.h>
#include <unistd.h>
#include <stdio.h>

void microsleep(int ms, int* interrupt) {
    static mach_timebase_info_data_t timebaseInfo;
    mach_timebase_info(&timebaseInfo);
    uint64_t start = mach_continuous_time();

    for (;;) {
        usleep(0);
        if (*interrupt > 0) {break;}
        uint64_t now = mach_continuous_time();
        uint64_t nanos = (now - start) * (timebaseInfo.numer / timebaseInfo.denom);
        if (nanos / 1e6 >= ms) {break;}
    }
}

void move_mouse(int x, int y) {
    CGEventRef get = CGEventCreate(NULL);
    CGPoint mouse = CGEventGetLocation(get);

    CGPoint position = { .x = x - mouse.x, .y = y - mouse.y};
    CGEventRef move = CGEventCreateMouseEvent(NULL, kCGEventMouseMoved, position, kCGMouseButtonLeft);

    CGEventPost(kCGSessionEventTap, move);
    CFRelease(move);
}

void scroll_mouse(int x, int y) {
		CGEventRef event;
		event = CGEventCreateScrollWheelEvent(NULL, kCGScrollEventUnitPixel, 2, y, x);
		CGEventPost(kCGHIDEventTap, event);
		CFRelease(event);
}

void send_key_event(int pid, bool down, int key) {
        CGEventRef keyEvent = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)key, down);
        assert(keyEvent != NULL);

        CGEventSetType(keyEvent, down ? kCGEventKeyDown : kCGEventKeyUp);
        if (pid == 0) {
            CGEventPost(kCGSessionEventTap, keyEvent);
        } else {
            CGEventPostToPid(pid, keyEvent);
        }
        CFRelease(keyEvent);
}

#endif