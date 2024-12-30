clang -c CaptureController.m -o CaptureController.o -fobjc-arc \
  -framework Foundation -framework ScreenCaptureKit

clang -c CaptureBridge.m -o CaptureBridge.o -fobjc-arc \
  -framework Foundation -framework ScreenCaptureKit

ar rcs libCapture.a CaptureController.o CaptureBridge.o