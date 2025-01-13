package routines

import (
	"github.com/nosyliam/revolution/opencv"
	. "github.com/nosyliam/revolution/pkg/common"
	. "github.com/nosyliam/revolution/pkg/control/actions"
	revimg "github.com/nosyliam/revolution/pkg/image"
	"image"
)

const DetectNightRoutineKind RoutineKind = "detect-night"

func calcAverage(hist opencv.Mat, totalPixels float32) float32 {
	sum := float32(0)
	for i := 0; i < hist.Rows(); i++ {
		sum += float32(i) * hist.GetFloatAt(i, 0)
	}
	return sum / totalPixels
}

func detectNight(macro *Macro) bool {
	screenshot := macro.Root.Window.Screenshot()
	bounds := screenshot.Bounds()
	offset := macro.MacroState.Object().BaseOriginY
	cropped := revimg.CropRGBA(screenshot, image.Rect(bounds.Dx()-50, offset+50, bounds.Dx(), offset+100))

	// Perform a histogram calculation to determine the average color intensity
	mat, err := opencv.NewMatFromBytes(cropped.Bounds().Dy(), cropped.Bounds().Dx(), opencv.MatTypeCV8UC4, cropped.Pix)
	if err != nil {
		macro.Action(Error("Failed to map night detection image!")(Discord))
		return false
	}
	defer mat.Close()

	channels := opencv.Split(mat)

	histSize := []int{256}
	ranges := []float64{0, 256}

	blueHist := opencv.NewMat()
	greenHist := opencv.NewMat()
	redHist := opencv.NewMat()
	defer blueHist.Close()
	defer greenHist.Close()
	defer redHist.Close()

	opencv.CalcHist([]opencv.Mat{channels[0]}, []int{0}, opencv.NewMat(), &blueHist, histSize, ranges, false)
	opencv.CalcHist([]opencv.Mat{channels[1]}, []int{0}, opencv.NewMat(), &greenHist, histSize, ranges, false)
	opencv.CalcHist([]opencv.Mat{channels[2]}, []int{0}, opencv.NewMat(), &redHist, histSize, ranges, false)

	totalPixels := float32(mat.Rows() * mat.Cols())

	blueAvg := calcAverage(blueHist, totalPixels)
	greenAvg := calcAverage(greenHist, totalPixels)
	redAvg := calcAverage(redHist, totalPixels)

	if blueAvg < 5 && greenAvg < 5 && redAvg < 5 {
		return true
	}
	return false
}

var DetectNightRoutine = Actions{
	Set(NightDetected, false),
	Loop(
		For(9),
		KeyPress(RotUp),
	),
	Loop(
		For(2),
		KeyPress(RotDown),
	),
	Condition(
		If(True(detectNight)),
		Info("Night Detected")(Status, Discord),
		Set(NightDetected, true),
	),
}

func init() {
	DetectNightRoutine.Register(DetectNightRoutineKind)
}
