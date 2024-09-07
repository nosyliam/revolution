import {
    Accordion,
    BackgroundImage,
    Button,
    Container,
    Grid, Group,
    NumberInput,
    Paper, Select,
    Stack,
    Switch,
    UnstyledButton
} from "@mantine/core";
import ControlBox from "../../components/ControlBox";
import React, {useContext} from "react";
import {Bees} from '../../assets/images/bees/bees'
import {RuntimeContext} from "../../hooks/useRuntime";

interface Mutation {
    value: string,
    label: string,
    minPercent?: number,
    maxPercent?: number,
    minNumeric?: number,
    maxNumeric?: number,
    percent: boolean
    numeric: boolean
}

const Mutations: {[key: string]: Mutation} = {
    "attack": {
        value: "attack",
        label: "Attack",
        minPercent: 2,
        maxPercent: 6,
        minNumeric: 1,
        maxNumeric: 2,
        percent: true,
        numeric: true
    },
    "convert": {
        value: "convert-amount",
        label: "Convert",
        minPercent: 10,
        maxPercent: 30,
        minNumeric: 20,
        maxNumeric: 70,
        percent: true,
        numeric: true
    },
    "gather": {
        value: "gather-amount",
        label: "Gather",
        minPercent: 10,
        maxPercent: 30,
        minNumeric: 2,
        maxNumeric: 78,
        percent: true,
        numeric: true
    },
    "movespeed": {
        value: "movespeed",
        label: "Movespeed",
        minNumeric: 2,
        maxNumeric: 6,
        percent: false,
        numeric: true
    },
    "energy": {
        value: "energy",
        label: "Energy",
        minPercent: 10,
        maxPercent: 40,
        percent: true,
        numeric: false
    },
    "ability-rate": {
        value: "ability-rate",
        label: "Ability Rate",
        minPercent: 1,
        maxPercent: 4,
        percent: true,
        numeric: false
    },
    "critical-chance": {
        value: "critical-chance",
        label: "Critical Chance",
        minPercent: 1,
        maxPercent: 3,
        percent: true,
        numeric: false
    },
    "instant-conversion": {
        value: "instant conversion",
        label: "Instant Conv.",
        minPercent: 8,
        maxPercent: 20,
        percent: true,
        numeric: false
    },
}

export default function JellyTool() {
    const runtime = useContext(RuntimeContext)
    const settings = runtime.Object("settings.tools.jellyTool")
    const enabled = settings.Value<boolean>("enabled", false)
    const allowedTypes = settings.List<string>("beeTypes").Values()
    const requireMutation = settings.Value<boolean>("requireMutation", false)
    const mutationType = settings.Value<string>("mutationType", "movespeed")
    const mutationValue = settings.Value<number>("mutationValue", 0)
    const rollLimit = settings.Value<number>("rollLimit", 0)
    const usePercent = settings.Value<boolean>("usePercent", false)

    return (
        <Grid gutter="xs" style={{flexGrow: 1, height: '100%'}} styles={{inner: {height: '100%'}}}>
            <Grid.Col span={7}>
                <Stack style={{height: '100%', flexGrow: 1, gap: 4}}>
                    <ControlBox title="Enabled">
                        <Switch
                            checked={enabled}
                            onChange={(event) => settings.Set<boolean>("enabled", event.currentTarget.checked)}
                        />
                    </ControlBox>
                    <ControlBox title="Allow All Types">
                        <Switch
                            checked={allowedTypes.length == Object.entries(Bees).length}
                            onChange={(event) => {}}
                            disabled={!enabled}
                        />
                    </ControlBox>
                    <ControlBox title="Require Mutation">
                        <Switch
                            checked={requireMutation}
                            onChange={(event) => settings.Set<boolean>("requireMutation", event.currentTarget.checked)}
                            disabled={!enabled}
                        />
                    </ControlBox>
                    <ControlBox
                        leftAction={
                        <Group style={{gap: 4}}>
                            <NumberInput w={60} size="xs" value={mutationValue} max={100} disabled={!enabled || !requireMutation}
                                         onChange={(value) => settings.Set<number>("mutationValue", value as number)} />
                            <Select
                                size="xs"
                                w={60}
                                withCheckIcon={false}
                                value={usePercent ? '%' : '+'}
                                onChange={(value) => settings.Set<boolean>("usePercent", value == '%')}
                                data={['+', '%']}
                                disabled={!enabled || !requireMutation}
                            />
                        </Group>
                        }
                    >
                        <Select
                            size="xs"
                            w={110}
                            withCheckIcon={false}
                            value={mutationType}
                            onChange={(value) => settings.Set<string>("mutationType", value as string)}
                            data={Object.values(Mutations)}
                            maxDropdownHeight={100}
                            disabled={!enabled || !requireMutation}
                        />
                    </ControlBox>
                    <ControlBox title="Limit Rolls">
                        <Switch
                            checked={rollLimit > 0}
                            onChange={(event) => settings.Set<number>("rollLimit", event.currentTarget.checked ? (rollLimit > 0 ? rollLimit : 100) : 0)}
                            disabled={!enabled}
                        />
                    </ControlBox>
                    <ControlBox title="Roll Limit">
                        <NumberInput w={100} size="xs" value={rollLimit} disabled={!enabled}
                                     onChange={(value) => settings.Set<number>("rollLimit", value as number)} />
                    </ControlBox>
                </Stack>
            </Grid.Col>
            <Grid.Col span={5}>
                <Paper shadow="xs" style={{height: '100%'}} p={4}>
                    <Stack style={{height: '100%', gap: 4}}>
                        <Container p={0} style={(theme) => ({
                            border: `1px solid ${theme.colors.gray[2]}`,
                            width: '163px',
                            flexGrow: 1,
                            borderRadius: 2
                        })}>
                            <Grid styles={{inner: {margin: 0, width: 'unset'}}} columns={5} m={2} pl={2}>
                                {Object.entries(Bees).map((bee, i) => {
                                    return (
                                        <Grid.Col span={1} p={0} key={bee[0]} style={{height: '30px', width: '30px'}}>
                                            <BackgroundImage src={bee[1][1]} style={{height: '29px', width: '29px'}}>
                                                <UnstyledButton size={14} style={{
                                                    width: '100%',
                                                    height: '100%',
                                                }}/>
                                            </BackgroundImage>
                                        </Grid.Col>
                                    )
                                })}
                            </Grid>
                        </Container>
                        <Grid columns={2} gutter="xs">
                            <Grid.Col span={1} pr={2}>
                                <Button size="compact-sm" style={{flexGrow: 1, width: '100%'}} disabled={!enabled}>Edit</Button>
                            </Grid.Col>
                            <Grid.Col span={1} pl={2}>
                                <Button size="compact-sm" style={{flexGrow: 1, width: '100%'}} disabled={!enabled}>Clear</Button>
                            </Grid.Col>
                        </Grid>
                    </Stack>
                </Paper>
            </Grid.Col>
        </Grid>
    )
}