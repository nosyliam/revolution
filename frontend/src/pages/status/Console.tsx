import React, {useRef, useState} from 'react';
import {Terminal} from 'xterm';
import {XTerm} from './XTerm';
import './Console.css';
import {EventsEmit, EventsOn} from "../../../wailsjs/runtime";

type CommandFunction = (...args: string[]) => string | null | undefined | void;
type CommandsMap = Record<string, CommandFunction>;

const commandsMap: CommandsMap = {
    help: () =>
        `Type 'help [command]' for more info\r
Available commands:\r
 listpatterns:  List all patterns\r
 execpattern:   Execute a pattern\r
 moveto:        Move to a field\r
 reset:         Reset to hive\r
 clear:         Clear the terminal\r\n`,

    echo: (...args: string[]) => {
        EventsEmit("command", "echo", args.join(''))
    },

    execpattern: (...args: string[]) => {
        EventsEmit("command", "execpattern", args[0])
    },

    moveto: (...args: string[]) => {
        EventsEmit("command", "moveto", args[0])
    },

    reset: (...args: string[]) => {
        EventsEmit("command", "reset")
    },

    listpatterns: () => {
        EventsEmit("command", "listpatterns")
    },

    clear: () => null,
};

const Console: React.FC = () => {
    const terminalRef = useRef<Terminal | null>(null);
    const [inputBuffer, setInputBuffer] = useState<string>('');
    const [initialized, setInitialized] = useState(false);

    // Initialize the terminal when the component mounts
    const handleInit = (terminal: Terminal) => {
        if (initialized)
            return
        setInitialized(true)
        terminalRef.current = terminal;

        // Set up terminal options and theme
        terminal.options = {
            theme: {
                background: '#ffffff',
                foreground: '#303030',
                cursor: '#333333',
                selection: '#cccccc',
            },
            disableStdin: true,
            cursorBlink: false,
            // @ts-ignore
            cursorStyle: 'none',
            scrollback: 1000,
            fontSize: 12,
            fontFamily: 'monospace',
        };

        // Print welcome message
        terminal.writeln('Welcome to the Revolution Console!');
        terminal.writeln('Type "help" to see all commands.\r\n');

        EventsOn('console', (text: string) => {
            terminal.writeln(text)
        })
    };

    const handleInput = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
        const terminal = terminalRef.current;
        if (!terminal) return;

        if (e.key === 'Enter') {
            e.preventDefault();
            //terminal.write('\r\n');
            const command = inputBuffer.trim();
            setInputBuffer(''); // Clear the input buffer

            // Process command
            if (command === 'clear') {
                terminal.clear();
            } else if (command) {
                const [cmd, ...args] = command.trim().split(/\s+/);
                const commandFunction = commandsMap[cmd];
                if (commandFunction) {
                    const result = commandFunction(...args);
                    if (result) {
                        terminal.writeln(result);
                    }
                } else {
                    terminal.writeln(`Unrecognized command: ${cmd}`);
                }
            }
            //terminal.write('$ ');
        } else if (e.key === 'Backspace') {
            e.preventDefault();
            if (inputBuffer.length > 0) {
                setInputBuffer((prev) => prev.slice(0, -1));
            }
        } else if (e.key.length === 1) {
            setInputBuffer((prev) => prev + e.key);
        }
    };

    return (
        <div id="terminal-container" style={{flexGrow: 1}}>
            <XTerm
                onInit={handleInit}
            />
            <div id="input-container">
                <div id="prompt">&gt;</div>
                <textarea
                    id="input-area"
                    value={inputBuffer}
                    onKeyDown={handleInput}
                    onChange={() => {}}
                    placeholder="Type your command here..."
                />
            </div>
        </div>
    );
};

export default Console;
