import classes from "./Switcher.module.css"
import textClasses from "../../assets/css/NonSelectableText.module.css"
import {Button, Container, Flex, Group, Text} from "@mantine/core";
import React from "react";
import {IconChevronLeft, IconChevronRight} from "@tabler/icons-react";
import TextContainer from "./TextContainer";

interface SwitcherProps {
    type: string
}

export default function Switcher(props: SwitcherProps) {
    return (
        <Group mah={26} style={{flexGrow: 1, gap: 4}}>
            <TextContainer label={props.type}>Default</TextContainer>
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