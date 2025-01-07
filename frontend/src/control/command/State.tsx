import {Box, Group, Stack} from "@mantine/core";
import {IconPlayerPauseFilled, IconPlayerPlayFilled, IconPlayerStopFilled} from '@tabler/icons-react';
import ActionButton from "./ActionButton";
import React, {useContext} from "react";
import {AccountSwitcher, PresetSwitcher} from "./Switcher";
import TextContainer from "./TextContainer";
import {KeyedObject, RuntimeContext} from "../../hooks/useRuntime";
import {Pause, Start, Stop} from "../../../wailsjs/go/main/Macro";


export default function State() {
    const runtime = useContext(RuntimeContext)
    const state = runtime.Object("state")
    const config = runtime.Object("state.config")
    const activeAccount = config.Value<string>("activeAccount", "default")
    const macros = state.List<KeyedObject>("macros").Values(true)
    const index = Object.fromEntries(macros.map((m) => [m.key, m.object]))

    const activeData = index[activeAccount]
    const loadCheck = (macros.length == 0 ? true : undefined) || (activeData == undefined ? true : undefined)

    return (
        <Box p={6}>
            <Group gap={4}>
                <Stack gap={4}>
                    <TextContainer label="Status">{ loadCheck ? 'Loading' : activeData.Concrete<string>('status') }</TextContainer>
                    <Group gap={4} style={{flexGrow: 1, flexWrap: 'nowrap'}}>
                        <ActionButton
                            action="Start"
                            icon={<IconPlayerPlayFilled size={16}/>}
                            execute={() => Start(activeAccount)}
                            disabled={loadCheck ?? (!activeData.Concrete<boolean>('paused') && activeData.Concrete<boolean>('running'))}
                        />
                        <ActionButton
                            action="Pause"
                            icon={<IconPlayerPauseFilled size={16}/>}
                            execute={() => Pause(activeAccount)}
                            disabled={loadCheck ?? (!activeData.Concrete<boolean>('running') || activeData.Concrete<boolean>('paused'))}
                        />
                        <ActionButton
                            action="Stop"
                            icon={<IconPlayerStopFilled size={16}/>}
                            execute={() => Stop(activeAccount)}
                            disabled={loadCheck ?? !activeData.Concrete<boolean>('running')}
                        />
                    </Group>
                </Stack>
                <Stack gap={4} style={{flexGrow: 1}}>
                    <PresetSwitcher/>
                    <AccountSwitcher/>
                </Stack>
            </Group>
        </Box>
    )
}