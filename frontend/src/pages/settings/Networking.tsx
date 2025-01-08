import {useContext} from "react";
import {RuntimeContext} from "../../hooks/useRuntime";
import {Group, Stack} from "@mantine/core";

export default function Networking() {
    const runtime = useContext(RuntimeContext)
    const state = runtime.State()
    const networking = state.Object("networking")

    const relayStarting = networking.Value("relayStarting", false)
    const relayActive = networking.Value("relayActive", false)
    const relayIdentity = networking.Value("relayIdentity", "")
    const connectedIdentities = networking.List<string>("connectedIdentities")

    const connectedIdentity = networking.Value("connectedIdentity", "")
    const availableIdentities = networking.List<string>("availableIdentities")



    return (
        <Group gap="xxs">

        </Group>
    )
}