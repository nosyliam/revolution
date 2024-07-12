#pragma once
#ifndef OS_H
#define OS_H

#if !defined(IS_MACOSX) && defined(__APPLE__) && defined(__MACH__)
	#define IS_MACOSX
	#include <dlfcn.h>
    #include <CoreFoundation/CoreFoundation.h>
    #include <CoreGraphics/CoreGraphics.h>
    #include <ApplicationServices/ApplicationServices.h>

    static Boolean(*gAXIsProcessTrustedWithOptions) (CFDictionaryRef);
    static CFStringRef* gkAXTrustedCheckOptionPrompt;

    bool check_ax_enabled(bool showPrompt) {
    	// Statically load all required functions one time
    	static dispatch_once_t once; dispatch_once (&once,
    	^{
    		// Open the framework
    		void* handle = dlopen("/System/Library/Frameworks/Application"
    			"Services.framework/ApplicationServices", RTLD_LAZY);

    		// Validate the handle
    		if (handle != NULL) {
    			*(void**) (&gAXIsProcessTrustedWithOptions) = dlsym (handle, "AXIsProcessTrustedWithOptions");
    			gkAXTrustedCheckOptionPrompt = (CFStringRef*) dlsym (handle, "kAXTrustedCheckOptionPrompt");
    		}
    	});

    	// Check for new OSX 10.9 function
    	if (gAXIsProcessTrustedWithOptions) {
    		// Check whether to show prompt
    		CFBooleanRef displayPrompt = showPrompt ? kCFBooleanTrue : kCFBooleanFalse;

    		// Convert display prompt value into a dictionary
    		const void* k[] = { *gkAXTrustedCheckOptionPrompt };
    		const void* v[] = { displayPrompt };
    		CFDictionaryRef o = CFDictionaryCreate(NULL, k, v, 1, NULL, NULL);

    		// Determine whether the process is actually trusted
    		bool result = (*gAXIsProcessTrustedWithOptions)(o);
    		// Free memory
    		CFRelease(o);
    		return result;
    	} else {
    		// Ignore deprecated warnings
    		#pragma clang diagnostic push
    		#pragma clang diagnostic ignored "-Wdeprecated-declarations"

    		// Check whether we have accessibility access
    		return AXAPIEnabled() || AXIsProcessTrusted();
    		#pragma clang diagnostic pop
    	}
    }
#endif /* IS_MACOSX */

#if !defined(IS_WINDOWS) && (defined(WIN32) || defined(_WIN32) || \
                            defined(__WIN32__) || defined(__WINDOWS__) || defined(__CYGWIN__))
	#define IS_WINDOWS
#endif /* IS_WINDOWS */
#endif