// @ts-ignore
import Terminal from 'react-console-emulator'
import ConsoleStyles from "./ConsoleStyles";
import {Paper} from "@mantine/core";

export default function Console() {
    return (
        <Terminal
            style={ConsoleStyles.container}
            contentStyle={ConsoleStyles.content}
            promptLabelStyle={ConsoleStyles.promptLabel}
            inputTextStyle={ConsoleStyles.inputText}
            promptLabel={<b>&gt;</b>}
            styleEchoBack='fullInherit'
            welcomeMessage={[
                'Welcome to the Revolution console.',
                'Type \'help\' for a list of commands. '
            ]}
        />
    )
}