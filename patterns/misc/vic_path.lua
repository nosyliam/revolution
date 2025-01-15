SetName("vic_path")

-- GO TO CANNON
function goToCannon()
    for i = 1, 5 do
        KeyPress(Key.ZoomOut)
    end

    Walk(Direction.Forward, 83.2)
    Walk(Direction.Backward, 4)
    Walk(Direction.Right, 4)
    Sleep(100)
    Walk(Direction.Right, 96.2)

    KeyDown(Key.Right)
    KeyPress(Key.Space)
    Sleep(1300)
    KeyPress(Key.Space)
    Sleep(600)
    KeyUp(Key.Right)
end

-- GO TO MT
function goToMT()
    Sleep(500)
    KeyPress(Key.E)
    for i = 1, 4 do
        KeyPress(Key.RotLeft)
    end

    Sleep(1700)

    KeyPress(Key.Space)
    KeyPress(Key.Space)
    KeyDown(Key.Right)

    Sleep(1000)

    KeyPress(Key.Space)
    KeyUp(Key.Right)

    -- Alignment
    Walk(Direction.Backward, 12)
    Walk(Direction.Left, 16)
    Walk(Direction.Forward, 16)

    Walk(Direction.Backward, 16)
    Walk(Direction.Right, 48)

    Walk(Direction.Forward, 130)
    Walk(Direction.Left, 191.1)
    Walk(Direction.Backward, 13)

    KeyPress(Key.Shift)

    for i = 1, 2 do
        KeyPress(Key.RotRight)
    end

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end
    for i = 1, 2 do
        KeyPress(Key.RotUp)
    end
    Sleep(50)

    if PerformDetection("mountain") then
        Exit()
    end

    KeyPress(Key.Shift)
end

-- GO TO SPID
function goToSpid()
    for i = 1, 4 do
        KeyPress(Key.RotLeft)
    end
    Sleep(50)

    KeyPress(Key.RotUp)
    Sleep(50)
    KeyPress(Key.RotUp)

    Walk(Direction.Backward, 193.5)
    Walk(Direction.Left, 93.6)

    KeyPress(Key.Space)
    Sleep(350)
    KeyPress(Key.Space)
    Sleep(3200)

    Walk(Direction.Forward, 78)
    Walk(Direction.Right, 26)
    Walk(Direction.Forward, 56)

    KeyPress(Key.Shift)

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end

    for i = 1, 5 do
        KeyPress(Key.RotUp)
    end

    for i = 1, 10 do
        KeyPress(Key.ZoomOut)
    end

    if PerformDetection("spider") then
        Exit()
    end

    KeyPress(Key.Shift)
end

-- GO TO CAC
function goToCac()
    for i = 1, 2 do
        KeyPress(Key.RotUp)
    end

    Walk(Direction.Left, 2.6)
    Walk(Direction.Backward, 14)
    Walk(Direction.Left, 30)

    KeyPress(Key.Space)
    Sleep(200)

    Walk(Direction.Forward, 26)
    Walk(Direction.Right, 39)

    WalkAsync(Direction.Forward, 130)
    Walk(Direction.Right, 130)

    Walk(Direction.Forward, 20)
    KeyPress(Key.Space)
    Walk(Direction.Forward, 46)

    Walk(Direction.Forward, 90)
    Walk(Direction.Right, 30)

    KeyPress(Key.Shift)

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end

    Walk(Direction.Backward, 15.6)
    Walk(Direction.Right, 10)

    Sleep(1000)

    if PerformDetection("cactus") then
        Exit()
    end

    Walk(Direction.Backward, 50.7)
    Walk(Direction.Right, 10)

    if PerformDetection("cactus") then
        Exit()
    end


    KeyPress(Key.Shift)

    Walk(Direction.Forward, 78)
end

-- GO TO ROSE
function goToRose()
    for i = 1, 6 do
        KeyPress(Key.RotUp)
    end

    Walk(Direction.Left, 104)

    KeyPress(Key.Space)

    Walk(Direction.Left, 70.2)
    Walk(Direction.Forward, 39)

    KeyPress(Key.Shift)

    for i = 1, 2 do
        KeyPress(Key.RotLeft)
    end

    Walk(Direction.Backward, 46.8)

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end

    for i = 1, 5 do
        KeyPress(Key.RotUp)
    end

    Walk(Direction.Right, 1.3)
    Walk(Direction.Backward, 1.3)

    for i = 1, 10 do
        KeyPress(Key.ZoomOut)
    end

    if PerformDetection("rose") then
        Exit()
    end

    for i = 1, 2 do
        KeyPress(Key.ZoomIn)
    end

    if PerformDetection("rose") then
        Exit()
    end
end

-- GO TO PEP
function goToPep()
    goToCannon()

    Walk(Direction.Right, 70.0)
    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(50)

    Walk(Direction.Right, 16)
    Walk(Direction.Forward, 12)

    KeyPress(Key.Space)
    KeyPress(Key.Space)

    Sleep(100)

    Walk(Direction.Forward, 58)
    KeyDown(Key.Forward)
    Walk(Direction.Right, 20)

    KeyPress(Key.Space)
    Sleep(800)
    KeyPress(Key.Space)
    Sleep(1800)
    KeyPress(Key.Space)
    Sleep(2500)

    KeyDown(Key.Right)
    Sleep(1000)
    KeyUp(Key.Forward)

    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(4000)
    KeyUp(Key.Right)

    Walk(Direction.Backward, 7.8)

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end

    for i = 1, 5 do
        KeyPress(Key.RotUp)
    end

    for i = 1, 10 do
        KeyPress(Key.ZoomOut)
    end
    Sleep(50)

    KeyPress(Key.Shift)
    Sleep(100)

    if PerformDetection("pepper") then
        Exit()
    end

    KeyPress(Key.Shift)
end

function pepToCannon()
    KeyPress(Key.Space)
    KeyDown(Key.Right)
    Sleep(200)
    KeyUp(Key.Right)
    Walk(Direction.Backward, 40)
    KeyPress(Key.Space)
    KeyDown(Key.Left)
    KeyDown(Key.Backward)
    Sleep(200)
    KeyPress(Key.Space)
    Sleep(2650)
    KeyUp(Key.Backward)
    Sleep(700)
    KeyUp(Key.Left)
    KeyPress(Key.Space)
    Sleep(1000)
end

goToPep()
pepToCannon()
goToMT()
goToSpid()
goToCac()
goToRose()