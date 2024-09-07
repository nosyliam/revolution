import {Divider, Paper, Text} from "@mantine/core";
import React from "react";

interface ControlBoxProps {
    title?: string
    description?: string
    leftAction?: React.ReactNode
}

export default function(props: React.PropsWithChildren<ControlBoxProps>) {
    return (
        <Paper shadow="xs" style={{
            width: '100%',
            display: 'flex',
            flexDirection: 'row',
            alignItems: 'center',
            padding: '4px',
        }}>
            {props.leftAction ?? <Text fz={14} ml={4}>{ props.title }</Text>}
            <div style={{flexGrow: 1}}/>
            { props.children }
        </Paper>
    )
}