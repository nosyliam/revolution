#import <Foundation/Foundation.h>
#import <ScreenCaptureKit/ScreenCaptureKit.h>

NS_ASSUME_NONNULL_BEGIN

typedef void (*FrameCallback)(
    unsigned char* data,
    size_t length,
    int width,
    int height,
    int stride
);

@interface CaptureController : NSObject

@property (nonatomic, assign) FrameCallback frameCallback;

- (instancetype)init NS_DESIGNATED_INITIALIZER;


// Start capturing a given SCWindow
- (void)startWindowCapture:(SCWindow *)scWin;

// Stop capturing
- (void)stopCapture;

@end

NS_ASSUME_NONNULL_END
