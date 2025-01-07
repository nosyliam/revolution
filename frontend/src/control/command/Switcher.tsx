import classes from "./Switcher.module.css"
import {Button, Container, Flex, Group, Text} from "@mantine/core";
import React from "react";
import {IconChevronLeft, IconChevronRight} from "@tabler/icons-react";
import TextContainer from "./TextContainer";

export function PresetSwitcher() {


    return (
        <Group mah={26} style={{flexGrow: 1, gap: 4}}>
            <TextContainer label="Preset">Default</TextContainer>
            <Button
                radius="50%"
                size="xs"
                disabled={true}
                classNames={classes}
            >
                <IconChevronLeft size={16} style={{marginRight: 2}}/>
            </Button>
            <Button
                radius="50%"
                size="xs"
                classNames={classes}
            >
                <IconChevronRight size={16} style={{marginLeft: 2}}/>
            </Button>
        </Group>
    )
}

export function AccountSwitcher() {
    return (
        <Group mah={26} style={{flexGrow: 1, gap: 4}}>
            <TextContainer label="Account">Default</TextContainer>
            <Button
                radius="50%"
                size="xs"
                classNames={classes}
            >
                <IconChevronLeft size={16} style={{marginRight: 2}}/>
            </Button>
            <Button
                radius="50%"
                size="xs"
                classNames={classes}
            >
                <IconChevronRight size={16} style={{marginLeft: 2}}/>
            </Button>
        </Group>
    )
}