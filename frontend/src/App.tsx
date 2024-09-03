import "@mantine/core/styles.css";
import {Box, Button, MantineProvider, Text} from "@mantine/core";
import {theme} from "./theme";
import {Path, RuntimeContext} from "./hooks/useRuntime";
import {useContext, useEffect} from "react";

function App() {
    const runtime = useContext(RuntimeContext)
    const settings = runtime.Object(new Path("settings"))
    const test = settings.Value<number>("test", 0)

    useEffect(runtime.Ready)

    return (
        <MantineProvider>
            <Box id="App">
                <Text size="xl">Value: { test }</Text>
                <Button onClick={() => settings.Set<number>("test", test + 1)}>Click me</Button>
            </Box>
        </MantineProvider>
    )
}

export default App
