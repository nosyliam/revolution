SetName("vic_path")

-- GO TO CANNON
function goToCannon()
    for i = 1, 5 do
        KeyPress(Key.ZoomOut)
    end


    -- was: walk(3200, "f")
    Walk(Direction.Forward, 83.2)

    -- was: walk(50, "b")
    Walk(Direction.Backward, 1.2)

    -- was: walk(3700, "r")
    Walk(Direction.Right, 96.2)



    KeyDown(Key.Space)
    KeyDown(Key.Right)
    Sleep(50)
    KeyUp(Key.Space)
    Sleep(500)
    KeyDown(Key.Space)
    Sleep(50)
    KeyUp(Key.Space)
    Sleep(1050)
    KeyUp(Key.Right)
end

-- GO TO MT
function goToMT()
    -- Send e (not in KeyPress enumerations)
    -- Send {, 4} => KeyPress(Key.RotLeft) x4
    for i = 1, 4 do
        KeyPress(Key.RotLeft)
    end
    Sleep(1200)

    -- AHK: {d down}, multiple KeyPress(Key.Space)...
    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(2800)
    -- done pressing d

    -- was: walk(5000, "f")
    Walk(Direction.Forward, 5000)
    -- was: walk(7350, "l")
    Walk(Direction.Left, 7350)
    -- was: walk(500, "b")
    Walk(Direction.Backward, 500)

    KeyPress(Key.Space)

    -- Send {. 2} => RotRight x2
    for i = 1, 2 do
        KeyPress(Key.RotRight)
    end

    -- Send {PgDn 9}, {PgUp 2}
    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end
    for i = 1, 2 do
        KeyPress(Key.RotUp)
    end
    Sleep(50)

    -- ;screenshot => PerformDetection()
    PerformDetection()
    Sleep(1000)

    KeyPress(Key.Space)
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

    -- was: walk(7350, "b")
    Walk(Direction.Backward, 7350)
    -- was: walk(3600, "l")
    Walk(Direction.Left, 3600)

    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(350)
    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(3000)

    -- was: walk(3000, "f")
    Walk(Direction.Forward, 3000)
    -- was: walk(1000, "r")
    Walk(Direction.Right, 1000)
    -- was: walk(1000, "f")
    Walk(Direction.Forward, 1000)

    KeyPress(Key.Space)

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end

    for i = 1, 5 do
        KeyPress(Key.RotUp)
        Sleep(50)
    end

    for i = 1, 10 do
        KeyPress(Key.ZoomOut)
    end

    -- ;screenshot => PerformDetection()
    PerformDetection()
    Sleep(1000)

    KeyPress(Key.Space)
end

-- GO TO CAC
function goToCac()
    for i = 1, 2 do
        KeyPress(Key.RotUp)
        Sleep(50)
    end

    -- was: walk(100, "l")
    Walk(Direction.Left, 100)
    -- was: walk(550, "b")
    Walk(Direction.Backward, 550)

    KeyPress(Key.Space)
    Sleep(300)

    -- was: walk(1000, "f")
    Walk(Direction.Forward, 1000)
    -- was: walk(250, "r")
    Walk(Direction.Right, 250)

    -- Press w+d for 5s, etc. => not in enumerations

    -- was: walk(50, "l")
    Walk(Direction.Left, 50)
    -- was: walk(50, "b")
    Walk(Direction.Backward, 50)

    KeyPress(Key.Space)
    Sleep(50)

    -- was: walk(6000, "f")
    Walk(Direction.Forward, 6000)
    -- was: walk(500, "r")
    Walk(Direction.Right, 500)

    KeyPress(Key.Space)

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end

    -- was: walk(600, "b")
    Walk(Direction.Backward, 600)

    -- ;screenshot => PerformDetection()
    PerformDetection()
    Sleep(1000)

    -- was: walk(1950, "b")
    Walk(Direction.Backward, 1950)

    -- ;screenshot => PerformDetection()
    PerformDetection()
    Sleep(1000)

    KeyPress(Key.Space)

    -- was: walk(3000, "f")
    Walk(Direction.Forward, 3000)
end

-- GO TO ROSE
function goToRose()
    for i = 1, 6 do
        KeyPress(Key.RotUp)
        Sleep(50)
    end

    -- was: walk(4000, "l")
    Walk(Direction.Left, 4000)

    KeyPress(Key.Space)
    KeyPress(Key.Space)

    -- was: walk(2700, "l")
    Walk(Direction.Left, 2700)
    -- was: walk(1500, "f")
    Walk(Direction.Forward, 1500)

    KeyPress(Key.Space)

    -- Send {, 2}
    for i = 1, 2 do
        KeyPress(Key.RotLeft)
    end

    -- was: walk(1800, "b")
    Walk(Direction.Backward, 1800)

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end

    for i = 1, 5 do
        KeyPress(Key.RotUp)
        Sleep(50)
    end

    -- was: walk(50, "r")
    Walk(Direction.Right, 50)
    -- was: walk(50, "b")
    Walk(Direction.Backward, 50)

    for i = 1, 10 do
        KeyPress(Key.ZoomOut)
    end

    for i = 1, 2 do
        KeyPress(Key.ZoomIn)
        Sleep(50)
    end

    -- ;screenshot => PerformDetection()
    PerformDetection()
end

-- GO TO PEP
function goToPep()
    goToCannon()

    -- was: walk(2500, "r")
    Walk(Direction.Right, 68.0)
    KeyPress(Key.Space)
    KeyPress(Key.Space)

    -- was: walk(750, "r")
    Walk(Direction.Right, 21)
    -- was: walk(50, "f")
    Walk(Direction.Forward, 2)

    KeyPress(Key.Space)
    KeyPress(Key.Space)

    -- was: walk(2200, "f")
    Walk(Direction.Forward, 57.2)
    KeyDown(Key.Forward)
    Walk(Direction.Right, 15.6)

    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(800)
    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(1800)
    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(2500)

    -- pressed d, etc. => not in enumerations
    KeyDown(Key.Right)
    Sleep(1000)
    KeyUp(Key.Forward)
    -- done pressing w

    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(4000)
    KeyUp(Key.Right)
    -- done pressing d

    -- was: walk(300, "b")
    Walk(Direction.Backward, 7.8)

    KeyPress(Key.Space)

    for i = 1, 9 do
        KeyPress(Key.RotDown)
    end

    for i = 1, 5 do
        KeyPress(Key.RotUp)
        Sleep(50)
    end

    for i = 1, 10 do
        KeyPress(Key.ZoomOut)
    end
    Sleep(50)

    -- ;screenshot => PerformDetection()
    Sleep(100000)
    PerformDetection()

    -- Send, ^p (not in enumerations)
    Sleep(1000)

    KeyPress(Key.Space)

    -- was: walk(50, "l")
    Walk(Direction.Left, 1.3)
    -- was: walk(1000, "b")
    Walk(Direction.Backward, 26)

    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(350)
    KeyPress(Key.Space)
    KeyPress(Key.Space)
    Sleep(1400)

    -- was: walk(5500, "l")
    Walk(Direction.Left, 143)
    -- was: walk(1650, "b")
    Walk(Direction.Backward, 42.9)
    -- was: walk(1200, "l")
    Walk(Direction.Left, 31.2)
end

goToPep()
--goToMT()
--goToSpid()
--goToCac()
--goToRose()