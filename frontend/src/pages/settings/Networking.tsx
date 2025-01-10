import React, {useContext, useState} from "react";
import {KeyedObject, Object, RuntimeContext} from "../../hooks/useRuntime";
import {
    Box,
    Button,
    Flex,
    Grid,
    Paper,
    ScrollArea,
    Stack,
    Switch,
    Table,
    Text,
    Tooltip,
    UnstyledButton
} from "@mantine/core";
import ControlBox from "../../components/ControlBox";
import {modals} from '@mantine/modals';
import {
    IconForbid,
    IconForbid2,
    IconForbidFilled,
    IconNetwork,
    IconPlus,
    IconRocket,
    IconSquare,
    IconStar
} from "@tabler/icons-react";
import {useHover} from "@mantine/hooks";

export default function Networking() {
    const [hoveredIndex, setHoveredIndex] = useState("")
    const [actionHovered, setActionHovered] = useState(false)
    const runtime = useContext(RuntimeContext)
    const state = runtime.Object("state.networking")
    const macroState = runtime.State()
    const settings = runtime.Object("settings.networking")
    const networking = macroState.Object("networking")

    const autoConnect = settings.Value("autoConnect", false)

    let relayStarting = networking.Value("relayStarting", false)
    let relayActive = networking.Value("relayActive", false)
    let connectedIdentities = networking.List<string>("connectedIdentities")

    let identity: any = networking.Value("identity", "Unknown/Unknown")
    let connectedIdentity = networking.Value("connectedIdentity", "")
    let availableRelays = networking.List<KeyedObject>("availableRelays").Values(true)
    console.log(state.List<KeyedObject>("savedRelays"))
    let savedRelays = state.List<KeyedObject>("savedRelays").Values(true)



    // @ts-ignore
    relayActive = false
    identity = 'Liam\'s MacBook Air/Liam/Default'

    const identityShorthand = identity.split('/').slice(-2).join('/')

    let mappedSavedRelays = savedRelays.map((v) => ({
        address: v.object.Concrete<string>("address"),
        identity: v.object.Concrete<string>("identity")
    }))

    let mappedRelays = availableRelays.map((v) => ({
        address: v.object.Concrete<string>("address"),
        identity: v.object.Concrete<string>("identity")
    }))

    mappedRelays = [
        {address: '1', identity: 'Liam\'s Mac/Liam/Default'},
        {address: '2', identity: 'Liam\'s PC/Macro 1/Default'},
        {address: '3', identity: 'Liam\'s PC/Macro 2/Default'},
        {address: '4', identity: 'Liam\'s PC/Macro 3/Default'},
        {address: '6', identity: 'Liam\'s PC/Macro 4/Default'},
        {address: '7', identity: 'Liam\'s PC/Macro 5/Default'},
    ]

    const mappedRelayRows = mappedRelays.map((relay) => {
        const saved = !relayActive && mappedSavedRelays.findIndex((r) => r.address == relay.address) != -1
        const actionStyle: React.CSSProperties = {
            display: hoveredIndex == relay.address || saved ? 'inherit' : 'none',
            position: 'absolute',
            bottom: 8,
            right: 8
        }

        const action = () => {
            if (relayActive) {
            } else {
                if (saved) {
                    state.List<KeyedObject>("savedRelays").Delete(relay.address!)
                } else {
                    const item = state.List<KeyedObject>("savedRelays").Append(relay.address!) as KeyedObject
                    item.object.Set("identity", relay.identity!)
                }
            }
        }

        return (
            <Table.Tr key={relay.address} style={{width: '100%'}}
                      onMouseEnter={() => setHoveredIndex(relay.address!)}
                      onMouseLeave={() => setHoveredIndex((i) => i == relay.address ? "" : i)}>
                <Table.Td pb={0} pr={6}><IconNetwork size={16}/></Table.Td>
                <Table.Td style={{display: 'flex', alignItems: 'center', paddingLeft: 0, position: 'relative'}}>
                    <Tooltip label={`${relay.identity} @ ${relay.address}`} withArrow>
                        <Text
                            fz={14}
                            mr={4}
                            c="gray.8"
                            style={{
                                whiteSpace: 'nowrap',
                                overflow: 'hidden',
                                textOverflow: 'ellipsis',
                                width: '130px',
                                direction: 'rtl',
                                textAlign: 'left',
                            }}
                        >
                            {relay.identity!.split('/').slice(-2).join('/')}
                        </Text>
                    </Tooltip>
                    <UnstyledButton onPointerEnter={() => setActionHovered(true)} onPointerLeave={() => setActionHovered(false)} onClick={action}>
                        {relayActive ? (actionHovered ? <IconForbidFilled size={18} style={actionStyle}/> : <IconForbid size={18} style={actionStyle}/>)
                            : <IconStar size={18} style={{...actionStyle, stroke: 'goldenrod'}} fill={actionHovered || saved ? 'goldenrod' : 'none'}/>}
                    </UnstyledButton>
                </Table.Td>
            </Table.Tr>
        )
    })


    const openDisconnectModal = () => modals.openConfirmModal({
        title: 'Please confirm your action',
        children: (
            <Text size="sm">
                Are you sure you want to disconnect from your current relay?
            </Text>
        ),
        labels: {confirm: 'Confirm', cancel: 'Cancel'},
        onCancel: () => console.log('Cancel'),
        onConfirm: () => console.log('Confirmed'),
    });

    const openShutdownModal = () => modals.openConfirmModal({
        title: 'Please confirm your action',
        children: (
            <Text size="sm">
                Are you sure you want to shut down your relay? This will disconnect all connected clients.
            </Text>
        ),
        labels: {confirm: 'Confirm', cancel: 'Cancel'},
        onCancel: () => console.log('Cancel'),
        onConfirm: () => console.log('Confirmed'),
    });


    return (
        <Grid gutter="xxs" style={{flexGrow: 1, height: '100%'}} styles={{inner: {height: '100%'}}}>
            <Grid.Col span={6}>
                <Stack style={{height: '100%', flexGrow: 1, gap: 4}}>
                    <ControlBox height={38} title="Identity">
                        <Tooltip label={identity} withArrow>
                            <Text
                                fz={14}
                                mr={4}
                                c="gray.8"
                                style={{
                                    whiteSpace: 'nowrap',
                                    overflow: 'hidden',
                                    textOverflow: 'ellipsis',
                                    width: '150px',
                                    direction: 'rtl',
                                    textAlign: 'right',
                                }}
                            >
                                {identityShorthand}
                            </Text>
                        </Tooltip>
                    </ControlBox>
                    <ControlBox height={38} title="Relay">
                        <Button
                            onClick={() => {
                            }}
                            rightSection={relayActive ? <IconSquare size={18} fill="white"/> : <IconRocket size={18}/>}
                            color={relayActive ? 'green' : 'blue'}
                            w={100}
                            size="xs"
                        >
                            {relayActive ? 'Active' : 'Start'}
                        </Button>
                    </ControlBox>
                    <ControlBox height={38} title="Auto-Connect">
                        <Switch
                            size="md"
                            checked={autoConnect}
                            onChange={(event) => settings.Set("autoConnect", event.currentTarget.checked)}
                            height={38}
                        />
                    </ControlBox>
                </Stack>
            </Grid.Col>
            <Grid.Col span={6}>
                <Paper style={(theme) => ({
                    border: `1px solid ${theme.colors.gray[4]}`,
                    height: '100%',
                })}>
                    <Flex justify="left" align="center" direction="row" style={(theme) => ({
                        borderBottom: `1px solid ${theme.colors.gray[4]}`,
                        height: '37px',
                    })}>
                        <Text fw={500} fz={14} ml={6}>{relayActive ? 'Connected Identities' : 'Available Relays'}</Text>
                        <Box style={{flexGrow: 1}}/>
                        <Button
                            mr={4}
                            disabled={relayActive}
                            radius="xl"
                            styles={{
                                root: {
                                    width: '22px', // Adjust to your preferred size
                                    height: '22px', // Ensure it's a square
                                    padding: 0,
                                },
                            }}
                        >
                            <IconPlus size={18} style={{marginBottom: '1px'}} color="white" fill="white"/>
                        </Button>
                    </Flex>

                    {/* Scrollable List */}
                    <ScrollArea style={{flexGrow: 1, height: '212px'}} type="scroll" offsetScrollbars>
                        <Table striped style={{width: '212px'}}>
                            <Table.Tbody style={{width: '100%'}}>{mappedRelayRows}</Table.Tbody>
                        </Table>
                    </ScrollArea>
                </Paper>
            </Grid.Col>
        </Grid>
    )
}