import "@mantine/core/styles.css";
import {Flex, MantineProvider, Tabs} from "@mantine/core";
import {theme} from "./theme";
import {Path, RuntimeContext} from "./hooks/useRuntime";
import {useContext, useEffect} from "react";
import State from "./control/command/State";
import Macro from "./control/Macro";
import {useViewportSize} from "@mantine/hooks";

function App() {
    const viewport = useViewportSize()
    const runtime = useContext(RuntimeContext)
    useEffect(runtime.Ready)

    return (
        <MantineProvider>
            <Flex
                direction="column"
                style={{height: viewport.height}}
            >
                <Macro/>
                <State/>
            </Flex>
        </MantineProvider>
    )
}

export default App
