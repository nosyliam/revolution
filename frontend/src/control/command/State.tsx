import {Group, Stack} from "@mantine/core";
import {IconPlayerPauseFilled, IconPlayerPlayFilled, IconPlayerStopFilled} from '@tabler/icons-react';
import ActionButton from "./ActionButton";
import React from "react";
import Switcher from "./Switcher";
import TextContainer from "./TextContainer";


export default function State() {
    const start = (macro: string) => {
    }
    const pause = (macro: string) => {
    }
    const stop = (macro: string) => {
    }

    return (
        <Stack gap={4} p={6}>
            <Group gap={4}>
                <TextContainer label="Status" width={322}>Ready</TextContainer>
                <Switcher type="Preset"/>
            </Group>
            <Group gap={4}>
                <Group gap={4}>
                    <ActionButton action="Start" icon={<IconPlayerPlayFilled size={16}/>} execute={start}/>
                    <ActionButton action="Pause" icon={<IconPlayerPauseFilled size={16}/>} execute={pause}/>
                    <ActionButton action="Stop" icon={<IconPlayerStopFilled size={16}/>} execute={stop}/>
                </Group>
                <Switcher type="Account"/>
            </Group>
        </Stack>
    )
}