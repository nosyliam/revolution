// CaptureController.m
#import "CaptureController.h"

@interface StreamOutput : NSObject <SCStreamOutput>
@property (nonatomic, weak) CaptureController *owner;
@end

@implementation StreamOutput

- (void)stream:(SCStream *)stream
didOutputSampleBuffer:(CMSampleBufferRef)sampleBuffer
       ofType:(SCStreamOutputType)type
{
    // Validate
    if (!CMSampleBufferIsValid(sampleBuffer) ||
        CMSampleBufferGetNumSamples(sampleBuffer) == 0) {
        return;
    }

    CVPixelBufferRef pixelBuffer = CMSampleBufferGetImageBuffer(sampleBuffer);
    if (!pixelBuffer) return;

    CVPixelBufferLockBaseAddress(pixelBuffer, kCVPixelBufferLock_ReadOnly);

    void *baseAddress = CVPixelBufferGetBaseAddress(pixelBuffer);
    size_t width      = CVPixelBufferGetWidth(pixelBuffer);
    size_t height     = CVPixelBufferGetHeight(pixelBuffer);
    size_t bpr        = CVPixelBufferGetBytesPerRow(pixelBuffer);
    size_t length     = bpr * height;

    // Copy the frame into a new buffer (BGRA, for example)
    unsigned char *bufferCopy = malloc(length);
    memcpy(bufferCopy, baseAddress, length);

    CVPixelBufferUnlockBaseAddress(pixelBuffer, kCVPixelBufferLock_ReadOnly);

    if (self.owner.frameCallback) {
        self.owner.frameCallback(bufferCopy, length, (int)width, (int)height, (int)bpr);
    } else {
        free(bufferCopy);
    }
}

@end


@interface CaptureController ()
@property (nonatomic, strong) SCStream *stream;
@property (nonatomic, strong) dispatch_queue_t captureQueue;
@property (nonatomic, strong) StreamOutput *streamOutput;
@end

@implementation CaptureController

- (instancetype)init {
    self = [super init];
    if (self) {
        // ...
    }
    return self;
}

- (void)startWindowCapture:(SCWindow *)scWin {
    // 1) Get the window’s global bounding rect.

    // Approach A: Use scWin.frame if supported on your macOS version:
    // (According to some docs, scWin.frame might be available on macOS 13+)
    CGRect windowRect = scWin.frame;

    NSLog(@"Window bounding rect: %@", NSStringFromRect(*(NSRect *)&windowRect));

    // 2) Map the bounding box to a specific display.
    CGDirectDisplayID matchingDisplayID = [self displayForWindowRect:windowRect];
    if (matchingDisplayID == kCGNullDirectDisplay) {
        NSLog(@"Could not find a matching display for window: %@", scWin.title);
        return;
    }

    // 3) Create an SCDisplay for that display. SCDisplay has a +displayWithID: constructor (10.15+) or you can fetch from SCShareableContent.displays
    SCDisplay *matchingSCDisplay = [self scDisplayForCGDirectDisplayID:matchingDisplayID];
    if (!matchingSCDisplay) {
        NSLog(@"Could not convert CGDirectDisplayID %u to an SCDisplay.", matchingDisplayID);
        return;
    }

    SCContentFilter *filter =
      [[SCContentFilter alloc] initWithDisplay:matchingSCDisplay
                             includingWindows:@[ scWin ]];

    SCStreamConfiguration *config = [SCStreamConfiguration new];
    config.queueDepth = 3;
    config.minimumFrameInterval = CMTimeMake(1, 30); // ~30fps
    // config.width = ...   // optionally scale the window
    // config.height = ...

    // Create stream
    self.stream = [[SCStream alloc] initWithFilter:filter
                                     configuration:config
                                          delegate:nil];
    if (!self.stream) {
        NSLog(@"Error: Failed to create SCStream");
        return;
    }

    self.streamOutput = [[StreamOutput alloc] init];
    self.streamOutput.owner = self;  // So the output can call self.frameCallback

    self.captureQueue = dispatch_queue_create("com.example.captureQueue", DISPATCH_QUEUE_SERIAL);

    NSError *addError = nil;
    BOOL didAdd = [self.stream addStreamOutput:self.streamOutput
                                         type:SCStreamOutputTypeScreen
                            sampleHandlerQueue:self.captureQueue
                                         error:&addError];
    if (!didAdd) {
        NSLog(@"Error adding stream output: %@", addError);
        return;
    }

    // 5) Start capture
    [self.stream startCaptureWithCompletionHandler:^(NSError * _Nullable error) {
        if (error) {
            NSLog(@"startCapture error: %@", error);
        } else {
            NSLog(@"Capture started for windowID=%u", scWin.windowID);
        }
    }];
}


- (CGDirectDisplayID)displayForWindowRect:(CGRect)windowRect {
    uint32_t maxDisplays = 16;
    CGDirectDisplayID displayIDs[16];
    uint32_t displayCount = 0;
    CGGetActiveDisplayList(maxDisplays, displayIDs, &displayCount);

    CGDirectDisplayID bestDisplay = kCGNullDirectDisplay;
    CGFloat bestOverlap = 0;

    for (uint32_t i = 0; i < displayCount; i++) {
        CGDirectDisplayID dID = displayIDs[i];
        CGRect displayBounds = CGDisplayBounds(dID);

        CGRect intersection = CGRectIntersection(displayBounds, windowRect);
        CGFloat overlapArea = intersection.size.width * intersection.size.height;
        if (overlapArea > bestOverlap) {
            bestOverlap = overlapArea;
            bestDisplay = dID;
        }
    }

    return bestDisplay;
}

/// Example utility to find/create an SCDisplay for a given CGDirectDisplayID
- (SCDisplay *)scDisplayForCGDirectDisplayID:(CGDirectDisplayID)displayID {
    // Approach A: If you already enumerated SCShareableContent, just loop content.displays
    // to find the SCDisplay whose displayID property == displayID.
    // For instance:
    __block SCDisplay *match = nil;
    dispatch_semaphore_t sem = dispatch_semaphore_create(0);

    [SCShareableContent getShareableContentWithCompletionHandler:^(SCShareableContent * _Nullable content, NSError * _Nullable error) {
        if (content) {
            for (SCDisplay *disp in content.displays) {
                if (disp.displayID == displayID) {
                    match = disp;
                    break;
                }
            }
        }
        dispatch_semaphore_signal(sem);
    }];

    dispatch_semaphore_wait(sem, DISPATCH_TIME_FOREVER);
    return match;
}

- (void)stopCapture {
    [self.stream stopCaptureWithCompletionHandler:^(NSError * _Nullable error) {
        if (error) {
            NSLog(@"stopCapture error: %@", error);
        } else {
            NSLog(@"Capture stopped.");
        }
    }];
}

@end
