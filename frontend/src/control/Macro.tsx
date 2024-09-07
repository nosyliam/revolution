import {Box, Flex, Tabs} from "@mantine/core";
import {
    IconShovel,
    IconTicket,
    IconBucket,
    IconTerminal2,
    IconTool,
    IconSettings,
} from "@tabler/icons-react";
import tabClasses from "./Tab.module.css"
import React from "react";
import Tools from "../pages/tools/Tools";

export default function MacroTabs() {
    return (
        <Flex style={{flexGrow: 1}}>
            <Tabs defaultValue="gallery">
                <Tabs.List bg="gray.0">
                    <Tabs.Tab classNames={tabClasses} value="gather" leftSection={<IconShovel size={16}/>}>
                        Gather
                    </Tabs.Tab>
                    <Tabs.Tab classNames={tabClasses} value="collect" leftSection={<IconTicket size={16}/>}>
                        Collect
                    </Tabs.Tab>
                    <Tabs.Tab classNames={tabClasses} value="planters" leftSection={<IconBucket size={16}/>}>
                        Planters
                    </Tabs.Tab>
                    <Tabs.Tab classNames={tabClasses} value="status" leftSection={<IconTerminal2 size={16}/>}>
                        Status
                    </Tabs.Tab>
                    <Tabs.Tab classNames={tabClasses} value="tools" leftSection={<IconTool size={16}/>}>
                        Tools
                    </Tabs.Tab>
                    <Tabs.Tab classNames={tabClasses} value="settings" leftSection={<IconSettings size={16}/>}>
                        Settings
                    </Tabs.Tab>
                </Tabs.List>

                <Box p={6} style={{height: 'calc(100% - 34px)'}}>
                    <Tabs.Panel value="tools" style={{height: '100%'}}>
                        <Tools/>
                    </Tabs.Panel>
                </Box>
            </Tabs>
        </Flex>
    )
}